package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"path"
)

// VERSION: set at build time with -ldflags
var VERSION string

// TIMESTAMP: set at build time with -ldflags
var TIMESTAMP string

// fetchOpts: user defined opts
type fetchOpts struct {
	repoUrl       string
	commitSha     string
	branch        string
	tagConstraint string
	githubToken   string
	fromPaths      []string
	ReleaseAssets []string
	Unpack        bool
	GpgPublicKey  string
	DownloadDir   string
}

// releaseDl: data to complete download of a release asset
type releaseDl struct {
	Asset     *GitHubReleaseAsset
	Name      string
	LocalPath string
	Tag       string
}

const opt_repo = "repo"
const opt_commit = "commit"
const opt_branch = "branch"
const opt_tag = "tag"
const opt_github_token = "oauth-token"
const opt_from_path = "from-path"
const opt_release_asset = "release-asset"
const opt_unpack = "unpack"
const opt_gpg_public_key = "gpg-public-key"

func main() {
	app := cli.NewApp()

	app.Name = "ghfetch"

	app.Usage = txt_usage + "   " + TIMESTAMP

	app.UsageText = usage_lead

	app.Version = VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  opt_repo,
			Usage: txt_repo,
		},
		cli.StringFlag{
			Name:  opt_commit,
			Usage: txt_commit,
		},
		cli.StringFlag{
			Name:  opt_branch,
			Usage: txt_branch,
		},
		cli.StringFlag{
			Name: opt_tag,
			Usage: txt_tag,
		},
		cli.StringSliceFlag{
			Name: opt_from_path,
			Usage: txt_from_path,
		},
		cli.StringSliceFlag{
			Name: opt_release_asset,
			Usage: txt_release_asset,
		},
		cli.BoolFlag{
			Name: opt_unpack,
			Usage: txt_unpack,
		},
		cli.StringFlag{
			Name: opt_gpg_public_key,
			Usage: txt_gpg_public_key,
		},
		cli.StringFlag{
			Name:   opt_github_token,
			Usage:  txt_token,
			EnvVar: "GITHUB_OAUTH_TOKEN,GITHUB_TOKEN",
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
	tags, err := FetchTags(o.repoUrl, o.githubToken)
	if err != nil {
		if err.errorCode == INVALID_GITHUB_TOKEN_OR_ACCESS_DENIED {
			return errors.New(getErrorMessage(INVALID_GITHUB_TOKEN_OR_ACCESS_DENIED, err.details))
		} else if err.errorCode == REPO_DOES_NOT_EXIST_OR_ACCESS_DENIED {
			return errors.New(getErrorMessage(REPO_DOES_NOT_EXIST_OR_ACCESS_DENIED, err.details))
		} else {
			return fmt.Errorf("Error occurred while getting tags from GitHub repo: %s", err)
		}
	}

	specific, desiredTag := isTagConstraintSpecificTag(o.tagConstraint)
	if !specific {
		// Find the specific release that matches the latest version constraint
		latestTag, err := getLatestAcceptableTag(o.tagConstraint, tags)
		if err != nil {
			if err.errorCode == INVALID_TAG_CONSTRAINT_EXPRESSION {
				return errors.New(getErrorMessage(INVALID_TAG_CONSTRAINT_EXPRESSION, err.details))
			} else {
				return fmt.Errorf("Error occurred while computing latest tag that satisfies version contraint expression: %s", err)
			}
		}
		desiredTag = latestTag

		fmt.Printf("Most suitable tag for constraint %s is %s\n", o.tagConstraint, desiredTag)
	}

	// Prepare the vars we'll need to download
	repo, err := ParseUrlIntoGitHubRepo(o.repoUrl, o.githubToken)
	if err != nil {
		return fmt.Errorf("Error occurred while parsing GitHub URL: %s", err)
	}

	// If no release assets and no source paths are specified, then by default, download all the source files from
	// the repo
	if len(o.fromPaths) == 0 && len(o.ReleaseAssets) == 0 {
		o.fromPaths = []string{"/"}
	}

	// Download any requested source files
	if err := o.downloadFromPaths(repo, desiredTag); err != nil {
		return err
	}

	// Download any requested release assets
	if err := o.downloadReleaseAssets(repo, desiredTag); err != nil {
		return err
	}

	return nil
}

func parseOptions(c *cli.Context) fetchOpts {
	localDownloadPath := c.Args().First()

	return fetchOpts{
		repoUrl:       c.String(opt_repo),
		commitSha:     c.String(opt_commit),
		branch:        c.String(opt_branch),
		tagConstraint: c.String(opt_tag),
		githubToken:   c.String(opt_github_token),
		fromPaths:     c.StringSlice(opt_from_path),
		ReleaseAssets: c.StringSlice(opt_release_asset),
		Unpack:        c.Bool(opt_unpack),
		GpgPublicKey:  c.String(opt_gpg_public_key),
		DownloadDir:   localDownloadPath,
	}
}

func validateOptions(o fetchOpts) error {
	if o.repoUrl == "" {
		return fmt.Errorf("The --%s flag is required. Run \"fetch --help\" for full usage info.", opt_repo)
	}

	if o.DownloadDir == "" {
		return fmt.Errorf("Missing required arguments specifying the local download dir. Run \"fetch --help\" for full usage info.")
	}

	if o.tagConstraint == "" && o.commitSha == "" && o.branch == "" {
		return fmt.Errorf("You must specify exactly one of --%s, --%s, or --%s. Run \"fetch --help\" for full usage info.", opt_tag, opt_commit, opt_branch)
	}

	if len(o.ReleaseAssets) > 0 && o.tagConstraint == "" {
		return fmt.Errorf("The --%s flag can only be used with --%s. Run \"fetch --help\" for full usage info.", opt_release_asset, opt_tag)
	}

	if len(o.ReleaseAssets) == 0 && o.Unpack {
		return fmt.Errorf("The --%s flag can only be used with --%s. Run \"fetch --help\" for full usage info.", opt_unpack, opt_release_asset)
	}

	if o.GpgPublicKey != "" {
		if len(o.ReleaseAssets) == 0 {
			return fmt.Errorf("The --%s flag can only be used with --%s. Run \"fetch --help\" for full usage info.", opt_gpg_public_key, opt_release_asset)
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
func (o *fetchOpts) downloadFromPaths(githubRepo GitHubRepo, latestTag string) error {
	if len(o.fromPaths) == 0 {
		return nil
	}

	// We respect GitHubCommit Hierarchy: "commitSha > GitTag > branch"
	// Note that commitSha and branch are empty unless user passed values.
	// getLatestAcceptableTag() ensures that we have a GitTag value regardless
	// of whether the user passed one or not.
	// So if the user specified nothing, we'd download the latest valid tag.
	gitHubCommit := GitHubCommit{
		Repo:       githubRepo,
		GitTag:     latestTag,
		branch: o.branch,
		commitSha:  o.commitSha,
	}

	// Download that release as a .zip file
	if gitHubCommit.commitSha != "" {
		fmt.Printf("Downloading git commit \"%s\" of %s ...\n", gitHubCommit.commitSha, githubRepo.Url)
	} else if gitHubCommit.branch != "" {
		fmt.Printf("Downloading latest commit from branch \"%s\" of %s ...\n", gitHubCommit.branch, githubRepo.Url)
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
	for _, fromPath := range o.fromPaths {
		fmt.Printf("Extracting files from <repo>%s to %s ...\n", fromPath, o.DownloadDir)
		if err := extractFiles(localZipFilePath, fromPath, o.DownloadDir); err != nil {
			return fmt.Errorf("Error occurred while extracting files from GitHub zip file: %s", err.Error())
		}
	}

	fmt.Println("Download and file extraction complete.")
	return nil
}

// Download the specified binary files that were uploaded as release assets to the specified GitHub release

func newAsset(name string, path string, asset *GitHubReleaseAsset, tag string) releaseDl {
	return releaseDl{Asset: asset, Name: name, LocalPath: path, Tag: tag}
}

func (o *fetchOpts) downloadReleaseAssets(repo GitHubRepo, tag string) error {
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
		if err := FetchReleaseAsset(repo, asset.Id, assetPath); err != nil {
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

func (a *releaseDl) verifyGpg(gpgKey string, rel GitHubReleaseApiResponse, githubRepo GitHubRepo) error {
	asc := findAscInRelease(a.Name, rel)
	ascPath := fmt.Sprintf("%s.asc", a.LocalPath)

	if asc == nil {
		return fmt.Errorf("No %s.asc or %s.asc.txt in release %s", a.Name, a.Name, a.Tag)
	}
	fmt.Printf("Downloading gpg sig %s to %s\n", asc.Name, ascPath)
	if err := FetchReleaseAsset(githubRepo, asc.Id, ascPath); err != nil {
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
