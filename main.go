package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"path"
)

// This variable is set at build time using -ldflags parameters. For more info, see:
// http://stackoverflow.com/a/11355611/483528
var VERSION string

type FetchOptions struct {
	RepoUrl       string
	CommitSha     string
	BranchName    string
	TagConstraint string
	GithubToken   string
	SourcePaths   []string
	ReleaseAssets []string
	Unpack        bool
	GpgPublicKey  string
	DownloadDir   string
}

type ReleaseAsset struct {
	Asset     *GitHubReleaseAsset
	Name      string
	LocalPath string
	Tag       string
}

const OPTION_REPO = "repo"
const OPTION_COMMIT = "commit"
const OPTION_BRANCH = "branch"
const OPTION_TAG = "tag"
const OPTION_GITHUB_TOKEN = "github-oauth-token"
const OPTION_SOURCE_PATH = "source-path"
const OPTION_RELEASE_ASSET = "release-asset"
const OPTION_UNPACK = "unpack"
const OPTION_GPG_PUBLIC_KEY = "gpg-public-key"

func main() {
	app := cli.NewApp()
	app.Name = "ghfetch"
	app.Usage = "download a github repo OR selected subfolders/files OR release attachments.\n" +
		"   You can checkout from a specific git commit, branch, or tag.\n" +
		"   Specify a constraint for tags that are semantic version strings!\n" +
		"   Choose to automatically unpack release attachment tars and gzips!!"
	app.UsageText = "ghfetch [global options] /my/downloads/dir\n" +
		"   See https://github.com/opsgang/fetch for examples, argument definitions, and additional docs."
	app.Version = VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  OPTION_REPO,
			Usage: "Required. Fully qualified url of the github repo.\n",
		},
		cli.StringFlag{
			Name:   OPTION_GITHUB_TOKEN,
			Usage:  "A GitHub Personal Access Token, required to download from a private repo.",
			EnvVar: "GITHUB_OAUTH_TOKEN,GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name:  OPTION_COMMIT,
			Usage: "Git commit SHA1 to download. Overrides --branch and --tag.",
		},
		cli.StringFlag{
			Name:  OPTION_BRANCH,
			Usage: "Git branch from which to checkout the latest commit. Overrides --tag.",
		},
		cli.StringFlag{
			Name: OPTION_TAG,
			Usage: "Git tag to download, expressed with Hashicorp's Version Constraint Operators.\n" +
				"\tIf empty, ghfetch will download the latest tag.\n" +
				"\tSee https://github.com/opsgang/fetch#version-constraint-operators for examples.",
		},
		cli.StringSliceFlag{
			Name: OPTION_SOURCE_PATH,
			Usage: "Subfolder (or file) to get from the repo. Subfolder is not created locally.\n" +
				"\tIf this or --release-asset aren't specified, all files are downloaded.\n" +
				"\tSpecify multiple times to download multiple folders or files.\n" +
				"\te.g. # puts libs/* and build.sh in to /opt/libs\n" +
				"\t\t--source-path='/libs/' --source-path='/scripts/build.sh' /opt/libs",
		},
		cli.StringSliceFlag{
			Name: OPTION_RELEASE_ASSET,
			Usage: "Name of github release attachment to download. Requires --tag.\n" +
				"\tSpecify multiple times to grab more than one attachment.\n" +
				"\te.g. # get foo.tgz and bar.txt from latest 1.x release attachments\n" +
				"\t\t--tag='~>1.0' --release-asset='foo.tgz' --release-asset='bar.txt'",
		},
		cli.BoolFlag{
			Name: OPTION_UNPACK,
			Usage: "Whether to unpack a compressed release attachment. Requires --release-asset.\n" +
				"\tOnly unpacks tars, tar-gzip and gzip, otherwise does nothing.\n" +
				"\te.g. # unpacks latest 1.x tag of foo.tgz in to /var/tmp/foo\n" +
				"\t\t--tag='~>1.0' --unpack --release-asset='foo.tgz'",
		},
		cli.StringFlag{
			Name: OPTION_GPG_PUBLIC_KEY,
			Usage: "Path to local armoured GPG public key to verify downloaded release assets.\n" +
				"\tRequires --release-asset.\n" +
				"\tIf set, will look for <asset-name>.asc or <asset-name>.asc.txt attached\n" +
				"\tto the chosen release. That signature and this local key will be used\n" +
				"\tfor gpg verification.\n" +
				"\tIf signature file not found, or verification fails, the release-asset is deleted.",
		},
	}

	app.Action = runFetchWrapper

	// Run the definition of App.Action
	app.Run(os.Args)
}

// We just want to call runFetch(), but app.Action won't permit us to return an error, so call a wrapper function instead.
func runFetchWrapper(c *cli.Context) {
	err := runFetch(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

// Run the ghfetch program
func runFetch(c *cli.Context) error {
	o := parseOptions(c)
	if err := validateOptions(o); err != nil {
		return err
	}

	// Get the tags for the given repo
	tags, err := FetchTags(o.RepoUrl, o.GithubToken)
	if err != nil {
		if err.errorCode == INVALID_GITHUB_TOKEN_OR_ACCESS_DENIED {
			return errors.New(getErrorMessage(INVALID_GITHUB_TOKEN_OR_ACCESS_DENIED, err.details))
		} else if err.errorCode == REPO_DOES_NOT_EXIST_OR_ACCESS_DENIED {
			return errors.New(getErrorMessage(REPO_DOES_NOT_EXIST_OR_ACCESS_DENIED, err.details))
		} else {
			return fmt.Errorf("Error occurred while getting tags from GitHub repo: %s", err)
		}
	}

	specific, desiredTag := isTagConstraintSpecificTag(o.TagConstraint)
	if !specific {
		// Find the specific release that matches the latest version constraint
		latestTag, err := getLatestAcceptableTag(o.TagConstraint, tags)
		if err != nil {
			if err.errorCode == INVALID_TAG_CONSTRAINT_EXPRESSION {
				return errors.New(getErrorMessage(INVALID_TAG_CONSTRAINT_EXPRESSION, err.details))
			} else {
				return fmt.Errorf("Error occurred while computing latest tag that satisfies version contraint expression: %s", err)
			}
		}
		desiredTag = latestTag

		fmt.Printf("Most suitable tag for constraint %s is %s\n", o.TagConstraint, desiredTag)
	}

	// Prepare the vars we'll need to download
	repo, err := ParseUrlIntoGitHubRepo(o.RepoUrl, o.GithubToken)
	if err != nil {
		return fmt.Errorf("Error occurred while parsing GitHub URL: %s", err)
	}

	// If no release assets and no source paths are specified, then by default, download all the source files from
	// the repo
	if len(o.SourcePaths) == 0 && len(o.ReleaseAssets) == 0 {
		o.SourcePaths = []string{"/"}
	}

	// Download any requested source files
	if err := o.downloadSourcePaths(repo, desiredTag); err != nil {
		return err
	}

	// Download any requested release assets
	if err := o.downloadReleaseAssets(repo, desiredTag); err != nil {
		return err
	}

	return nil
}

func parseOptions(c *cli.Context) FetchOptions {
	localDownloadPath := c.Args().First()
	sourcePaths := c.StringSlice(OPTION_SOURCE_PATH)

	return FetchOptions{
		RepoUrl:       c.String(OPTION_REPO),
		CommitSha:     c.String(OPTION_COMMIT),
		BranchName:    c.String(OPTION_BRANCH),
		TagConstraint: c.String(OPTION_TAG),
		GithubToken:   c.String(OPTION_GITHUB_TOKEN),
		SourcePaths:   sourcePaths,
		ReleaseAssets: c.StringSlice(OPTION_RELEASE_ASSET),
		Unpack:        c.Bool(OPTION_UNPACK),
		GpgPublicKey:  c.String(OPTION_GPG_PUBLIC_KEY),
		DownloadDir:   localDownloadPath,
	}
}

func validateOptions(o FetchOptions) error {
	if o.RepoUrl == "" {
		return fmt.Errorf("The --%s flag is required. Run \"fetch --help\" for full usage info.", OPTION_REPO)
	}

	if o.DownloadDir == "" {
		return fmt.Errorf("Missing required arguments specifying the local download dir. Run \"fetch --help\" for full usage info.")
	}

	if o.TagConstraint == "" && o.CommitSha == "" && o.BranchName == "" {
		return fmt.Errorf("You must specify exactly one of --%s, --%s, or --%s. Run \"fetch --help\" for full usage info.", OPTION_TAG, OPTION_COMMIT, OPTION_BRANCH)
	}

	if len(o.ReleaseAssets) > 0 && o.TagConstraint == "" {
		return fmt.Errorf("The --%s flag can only be used with --%s. Run \"fetch --help\" for full usage info.", OPTION_RELEASE_ASSET, OPTION_TAG)
	}

	if len(o.ReleaseAssets) == 0 && o.Unpack {
		return fmt.Errorf("The --%s flag can only be used with --%s. Run \"fetch --help\" for full usage info.", OPTION_UNPACK, OPTION_RELEASE_ASSET)
	}

	if o.GpgPublicKey != "" {
		if len(o.ReleaseAssets) == 0 {
			return fmt.Errorf("The --%s flag can only be used with --%s. Run \"fetch --help\" for full usage info.", OPTION_GPG_PUBLIC_KEY, OPTION_RELEASE_ASSET)
		}

		// check file is readable
		reader, err := os.Open(o.GpgPublicKey)
		if err != nil {
			return fmt.Errorf("GPG public key %s is not a readable file.", o.GpgPublicKey)
		}
		defer reader.Close()
	}

	return nil
}

// Download the specified source files from the given repo
func (o *FetchOptions) downloadSourcePaths(githubRepo GitHubRepo, latestTag string) error {
	if len(o.SourcePaths) == 0 {
		return nil
	}

	// We respect GitHubCommit Hierarchy: "CommitSha > GitTag > BranchName"
	// Note that CommitSha and BranchName are empty unless user passed values.
	// getLatestAcceptableTag() ensures that we have a GitTag value regardless
	// of whether the user passed one or not.
	// So if the user specified nothing, we'd download the latest valid tag.
	gitHubCommit := GitHubCommit{
		Repo:       githubRepo,
		GitTag:     latestTag,
		BranchName: o.BranchName,
		CommitSha:  o.CommitSha,
	}

	// Download that release as a .zip file
	if gitHubCommit.CommitSha != "" {
		fmt.Printf("Downloading git commit \"%s\" of %s ...\n", gitHubCommit.CommitSha, githubRepo.Url)
	} else if gitHubCommit.BranchName != "" {
		fmt.Printf("Downloading latest commit from branch \"%s\" of %s ...\n", gitHubCommit.BranchName, githubRepo.Url)
	} else if gitHubCommit.GitTag != "" {
		fmt.Printf("Downloading tag \"%s\" of %s ...\n", latestTag, githubRepo.Url)
	} else {
		return fmt.Errorf("The commit sha, tag, and branch name are all empty.")
	}

	localZipFilePath, err := downloadGithubZipFile(gitHubCommit, githubRepo.Token)
	if err != nil {
		return fmt.Errorf("Error occurred while downloading zip file from GitHub repo: %s", err)
	}
	defer cleanupZipFile(localZipFilePath)

	// Unzip and move the files we need to our destination
	for _, sourcePath := range o.SourcePaths {
		fmt.Printf("Extracting files from <repo>%s to %s ...\n", sourcePath, o.DownloadDir)
		if err := extractFiles(localZipFilePath, sourcePath, o.DownloadDir); err != nil {
			return fmt.Errorf("Error occurred while extracting files from GitHub zip file: %s", err.Error())
		}
	}

	fmt.Println("Download and file extraction complete.")
	return nil
}

// Download the specified binary files that were uploaded as release assets to the specified GitHub release

func newAsset(name string, path string, asset *GitHubReleaseAsset, tag string) ReleaseAsset {
	return ReleaseAsset{Asset: asset, Name: name, LocalPath: path, Tag: tag}
}

func (o *FetchOptions) downloadReleaseAssets(repo GitHubRepo, tag string) error {
	if len(o.ReleaseAssets) == 0 {
		return nil
	}

	release, err := GetGitHubReleaseInfo(repo, tag)
	if err != nil {
		return err
	}

	// ... create download dir
	os.MkdirAll(o.DownloadDir, 0755)
	for _, assetName := range o.ReleaseAssets {
		asset := findAssetInRelease(assetName, release)
		if asset == nil {
			return fmt.Errorf("Could not find asset %s in release %s", assetName, tag)
		}

		assetPath := path.Join(o.DownloadDir, asset.Name)
		a := newAsset(assetName, assetPath, asset, tag)
		fmt.Printf("Downloading release asset %s to %s\n", asset.Name, assetPath)
		if err := DownloadReleaseAsset(repo, asset.Id, assetPath); err != nil {
			return err
		}

		if o.GpgPublicKey != "" {
			err := a.verifyGpg(o.GpgPublicKey, release, repo)
			if err != nil {
				fmt.Printf("Deleting unverified asset %s\n", assetPath)
				if remErr := os.Remove(assetPath); remErr != nil {
					return fmt.Errorf("%s\nCould not delete it: %s!", err, remErr)
				}

				return err
			}
		}

		if o.Unpack {
			if err := Unpack(assetPath, o.DownloadDir); err != nil {
				return err
			}
		}
	}

	fmt.Println("Download of release assets complete.")
	return nil
}

func (a *ReleaseAsset) verifyGpg(gpgKey string, rel GitHubReleaseApiResponse, githubRepo GitHubRepo) error {
	asc := findAscInRelease(a.Name, rel)
	ascPath := fmt.Sprintf("%s.asc", a.LocalPath)

	if asc == nil {
		return fmt.Errorf("No %s.asc or %s.asc.txt in release %s", a.Name, a.Name, a.Tag)
	}
	fmt.Printf("Downloading gpg sig %s to %s\n", asc.Name, ascPath)
	if err := DownloadReleaseAsset(githubRepo, asc.Id, ascPath); err != nil {
		return err
	}

	err := GpgVerify(gpgKey, ascPath, a.LocalPath)
	if warning := os.Remove(ascPath); warning != nil {
		fmt.Printf("Could not remove sig file %s\n", ascPath)
	}

	return err
}

func findAssetInRelease(assetName string, release GitHubReleaseApiResponse) *GitHubReleaseAsset {
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return &asset
		}
	}

	return nil
}

func findAscInRelease(assetName string, release GitHubReleaseApiResponse) *GitHubReleaseAsset {
	for _, asset := range release.Assets {
		asc := fmt.Sprintf("%s.asc", assetName)
		ascTxt := fmt.Sprintf("%s.asc.txt", assetName)
		if asset.Name == asc || asset.Name == ascTxt {
			return &asset
		}
	}

	return nil
}

// Delete the given zip file.
func cleanupZipFile(localZipFilePath string) error {
	err := os.Remove(localZipFilePath)
	if err != nil {
		return fmt.Errorf("Failed to delete local zip file at %s", localZipFilePath)
	}

	return nil
}

func getErrorMessage(errorCode int, errorDetails string) string {
	switch errorCode {
	case INVALID_TAG_CONSTRAINT_EXPRESSION:
		return fmt.Sprintf(`
The --tag value you entered is not a valid constraint expression.
See https://github.com/opsgang/fetch#version-constraint-operators for examples.

Underlying error message:
%s
`, errorDetails)
	case INVALID_GITHUB_TOKEN_OR_ACCESS_DENIED:
		return fmt.Sprintf(`
Received an HTTP 401 Response when attempting to query the repo for its tags.

Either your GitHub OAuth Token is invalid, or that you don't have access to
the repo with that token. Is the repo private?

Underlying error message:
%s
`, errorDetails)
	case REPO_DOES_NOT_EXIST_OR_ACCESS_DENIED:
		return fmt.Sprintf(`
Received an HTTP 404 Response when attempting to query the repo for its tags.

Either the URL does not exist, or you don't have permission to access it.
If the repo is private, you will need to set GITHUB_TOKEN (or GITHUB_OAUTH_TOKEN)
in the env before invoking fetch.

Underlying error message:
%s
`, errorDetails)
	}

	return ""
}
