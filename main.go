// Package main is the Highlander of namespaces.
// *There can be only one.*
package main

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/urfave/cli"
)

// VERSION : set at build time with -ldflags
var VERSION string

// TIMESTAMP : set at build time with -ldflags
var TIMESTAMP string

// fetchOpts : user defined opts
type fetchOpts struct {
	repoUrl       string
	commitSha     string
	branch        string
	tagConstraint string
	apiToken      string
	fromPaths     []string
	relAssets     []string
	tagRegex      string
	semverMatcher *regexp.Regexp
	unpack        bool
	verbose       bool
	whichTag      bool
	gpgPubKey     string
	destDir       string
	timeout       int
}

// releaseDl : data to complete download of a release asset
type releaseDl struct {
	name      string
	localPath string
	tag       string
	verbose   bool
}

const optRepo = "repo"
const optCommit = "commit"
const optBranch = "branch"
const optTag = "tag"
const optApiToken = "api-token"
const optFromPath = "from-path"
const optReleaseAsset = "release-asset"
const optUnpack = "unpack"
const optGpgPubKey = "gpg-public-key"
const optVerbose = "verbose"
const optWhichTag = "which-tag"
const optTimeout = "timeout"
const optTagRegex = "tag-regex"

const semverRegx = `[vV]?(?P<semver>(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?)`

func main() {
	app := cli.NewApp()

	defaultFlagStringer := cli.FlagStringer

	// prefer help text to separate long flag descriptions with newline.
	cli.FlagStringer = func(f cli.Flag) string {
		return fmt.Sprintf("%s\n\t", defaultFlagStringer(f))
	}

	app.Name = "ghfetch"

	app.Usage = txtUsage + "   " + TIMESTAMP

	app.UsageText = usageLead

	app.Version = VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  optRepo,
			Usage: txtRepo,
		},
		cli.StringFlag{
			Name:  optCommit,
			Usage: txtCommit,
		},
		cli.StringFlag{
			Name:  optBranch,
			Usage: txtBranch,
		},
		cli.StringFlag{
			Name:  optTag,
			Usage: txtTag,
		},
		cli.StringSliceFlag{
			Name:  optFromPath,
			Usage: txtFromPath,
		},
		cli.StringSliceFlag{
			Name:  optReleaseAsset,
			Usage: txtReleaseAsset,
		},
		cli.BoolFlag{
			Name:  optUnpack,
			Usage: txtUnpack,
		},
		cli.BoolFlag{
			Name:  optVerbose,
			Usage: txtVerbose,
		},
		cli.BoolFlag{
			Name:  optWhichTag,
			Usage: txtWhichTag,
		},
		cli.StringFlag{
			Name:  optTagRegex,
			Usage: txtTagRegex,
		},
		cli.StringFlag{
			Name:  optGpgPubKey,
			Usage: txtGpgPubKey,
		},
		cli.StringFlag{
			Name:   optApiToken,
			Usage:  txtToken,
			EnvVar: "API_TOKEN,API_OAUTH_TOKEN,GITHUB_TOKEN,GITHUB_OAUTH_TOKEN",
		},
		cli.IntFlag{
			Name:  optTimeout,
			Value: 120,
			Usage: txtTimeout,
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
func runFetch(c *cli.Context) (err error) {
	o := parseOptions(c)
	if err := validateOptions(o); err != nil {
		return err
	}

	if o.tagRegex != "" {
		tokenToReplace := regexp.MustCompile(`SEMVER`)
		// ... insert our semver matching pattern in to the user-defined pattern.
		pattern := tokenToReplace.ReplaceAllLiteralString(o.tagRegex, semverRegx)
		if pattern == o.tagRegex {
			return fmt.Errorf("--tag-regex pattern must include string SEMVER. See --help")
		}

		// wrap in defer and use recover() to avoid ugly panic msg?
		o.semverMatcher = regexp.MustCompile(pattern)
	}

	if o.apiToken == "" {
		fmt.Println("WARNING: no api token provided - rate-limited and can't access private repos")
	}

	o.setTimeout()

	// Prepare the vars we'll need to download
	r, err := urlToGitHubRepo(o.repoUrl, o.apiToken)
	if err != nil {
		return fmt.Errorf("Error occurred while parsing GitHub URL: %s", err)
	}

	if err := o.do(r); err != nil {
		return err
	}

	return
}

func parseOptions(c *cli.Context) fetchOpts {
	localDownloadPath := c.Args().First()

	return fetchOpts{
		repoUrl:       c.String(optRepo),
		commitSha:     c.String(optCommit),
		branch:        c.String(optBranch),
		tagConstraint: c.String(optTag),
		apiToken:      c.String(optApiToken),
		fromPaths:     c.StringSlice(optFromPath),
		relAssets:     c.StringSlice(optReleaseAsset),
		tagRegex:      c.String(optTagRegex),
		timeout:       c.Int(optTimeout),
		unpack:        c.Bool(optUnpack),
		verbose:       c.Bool(optVerbose),
		whichTag:      c.Bool(optWhichTag),
		gpgPubKey:     c.String(optGpgPubKey),
		destDir:       localDownloadPath,
	}
}

func validateOptions(o fetchOpts) error {
	if o.repoUrl == "" {
		return fmt.Errorf("The --%s flag is required.", optRepo)
	}

	if o.destDir == "" && !o.whichTag {
		return fmt.Errorf("Final argument must be the destination dir (unless calling --which-tag)")
	}

	// ... must specify only one of sha, tag or branch
	c := 0
	for _, val := range []string{o.tagConstraint, o.commitSha, o.branch} {
		if val != "" {
			c = c + 1
		}
	}
	if c != 1 {
		return fmt.Errorf(
			"You must specify one (and only one) of --%s, --%s, or --%s.",
			optTag, optCommit, optBranch,
		)
	}

	if len(o.relAssets) > 0 && o.tagConstraint == "" {
		return fmt.Errorf("The --%s flag can only be used with --%s.", optReleaseAsset, optTag)
	}

	if o.tagConstraint == "" && o.whichTag {
		return fmt.Errorf("The --%s flag makes no sense without --%s.", optWhichTag, optTag)
	}

	if len(o.relAssets) > 0 && len(o.fromPaths) > 0 {
		return fmt.Errorf("Specify only --%s or --%s, not both.", optReleaseAsset, optFromPath)
	}

	if len(o.relAssets) == 0 && o.unpack {
		return fmt.Errorf("The --%s flag can only be used with --%s.", optUnpack, optReleaseAsset)
	}

	if o.tagRegex != "" && o.tagConstraint == "" {
		return fmt.Errorf("The --%s flag can only be used with --%s.", optTagRegex, optTag)
	}

	if o.gpgPubKey != "" {
		if len(o.relAssets) == 0 {
			return fmt.Errorf("The --%s flag can only be used with --%s.", optGpgPubKey, optReleaseAsset)
		}

		// check file is readable
		reader, err := os.Open(o.gpgPubKey)
		if err != nil {
			return fmt.Errorf("GPG public key %s is not a readable file.", o.gpgPubKey)
		}
		defer reader.Close()
	}

	if o.timeout <= 0 {
		return fmt.Errorf("--timeout expects a POSITIVE, non-zero number of seconds!, not %d", o.timeout)
	}

	return nil
}

// Download the specified source files from the given repo
func (o *fetchOpts) downloadFromPaths(r repo, latestTag string) error {
	if len(o.fromPaths) == 0 {
		return nil
	}

	// We respect commit Hierarchy: "commitSha > GitTag > branch"
	// Note that commitSha and branch are empty unless user passed values.
	// bestFitTag() ensures that we have a GitTag value regardless
	// of whether the user passed one or not.
	// So if the user specified nothing, we'd download the latest valid tag.
	c := commit{
		Repo:      r,
		GitTag:    latestTag,
		branch:    o.branch,
		commitSha: o.commitSha,
	}

	// Download that release as a .zip file
	if c.commitSha != "" {
		fmt.Printf("Downloading git commit \"%s\" of %s ...\n", c.commitSha, r.Url)
	} else if c.branch != "" {
		fmt.Printf("Downloading latest commit from branch \"%s\" of %s ...\n", c.branch, r.Url)
	} else if c.GitTag != "" {
		fmt.Printf("Downloading tag \"%s\" of %s ...\n", latestTag, r.Url)
	} else {
		return fmt.Errorf("The commit sha, tag, and branch name are all empty.")
	}

	localZipFilePath, _, err := getSrcZip(c, r.Token)
	if err != nil {
		return fmt.Errorf("Error occurred while downloading zip file from GitHub repo: %s", err)
	}
	defer cleanupZipFile(localZipFilePath)

	// Unzip and move the files we need to our destination
	for _, fromPath := range o.fromPaths {
		fmt.Printf("Extracting files from <repo>%s to %s ...\n", fromPath, o.destDir)
		if err := extractFiles(localZipFilePath, fromPath, o.destDir); err != nil {
			return fmt.Errorf("Error occurred while extracting files from GitHub zip file: %s", err)
		}
	}

	fmt.Println("Download and file extraction complete.")
	return nil
}

// newAsset ():
//
func newAsset(name string, path string, tag string, verbose bool) releaseDl {
	return releaseDl{name, path, tag, verbose}
}

// downloadReleaseAssetts ():
// Download the user-defined release attachments.
// Also performs GPG check if needed.
func (o *fetchOpts) downloadReleaseAssets(r repo, tag string) error {
	if len(o.relAssets) == 0 {
		return nil
	}

	release, err := GetGitHubReleaseInfo(r, tag)
	if err != nil {
		fmt.Println("getting release info")
		return err
	}

	// ... create download dir
	os.MkdirAll(o.destDir, 0755)
	for _, assetName := range o.relAssets {
		asset := findAssetInRelease(assetName, release)
		if asset == nil {
			return fmt.Errorf("Could not find asset %s in release %s", assetName, tag)
		}

		assetPath := path.Join(o.destDir, asset.Name)
		a := newAsset(assetName, assetPath, tag, o.verbose)
		fmt.Printf("Downloading release asset %s to %s\n", asset.Name, assetPath)
		if err := FetchReleaseAsset(r, asset.Id, assetPath); err != nil {
			return err
		}

		if o.gpgPubKey != "" {
			err := a.verifyGpg(o.gpgPubKey, release, r)
			if err != nil {
				fmt.Printf("Deleting unverified asset %s\n", assetPath)
				if remErr := os.Remove(assetPath); remErr != nil {
					return fmt.Errorf("%s\nCould not delete it: %s!", err, remErr)
				}

				return err
			}
		}

		if o.unpack {
			if err := o.doUnpack(assetPath); err != nil {
				return err
			}
		}
	}

	fmt.Println("Download of release assets complete.")
	return nil
}

func (a *releaseDl) verifyGpg(gpgKey string, rel release, gr repo) error {
	asc := findAscInRelease(a.name, rel)
	ascPath := fmt.Sprintf("%s.asc", a.localPath)

	if asc == nil {
		return fmt.Errorf("No %s.asc or %s.asc.txt in release %s", a.name, a.name, a.tag)
	}

	if a.verbose {
		fmt.Printf("Downloading gpg sig %s to %s\n", asc.Name, ascPath)
	}

	if err := FetchReleaseAsset(gr, asc.Id, ascPath); err != nil {
		return err
	}

	err := gpgVerify(gpgKey, ascPath, a.localPath)
	if warning := os.Remove(ascPath); warning != nil {
		fmt.Printf("Could not remove sig file %s\n", ascPath)
	}

	return err
}

func findAssetInRelease(assetName string, release release) *relAsset {
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return &asset
		}
	}

	return nil
}

func findAscInRelease(assetName string, release release) *relAsset {
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

// Return true if the given slice contains the given string
func stringInSlice(s string, slice []string) bool {
	for _, val := range slice {
		if val == s {
			return true
		}
	}
	return false
}
