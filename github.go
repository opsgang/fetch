package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type repo struct {
	Url   string // The URL of the GitHub repo
	Owner string // The GitHub account name under which the repo exists
	Name  string // The GitHub repo name
	Token string // The personal access token to access this repo (if it's a private repo)
	Api   string // https://api.github.com etc (or stub server - see github_test.go)
}

// cl : our single instance of http.Client to be reused throughout.
var cl *http.Client

type headers map[string]string

// Hierarchy:
// * commitSha > GitTag > branch
// * Example: GitTag and branch are both specified; use the GitTag
// * Example: GitTag and commitSha are both specified; use the commitSha
// * Example: branch alone is specified; use branch

// commit {}:
// A specific git commit.
type commit struct {
	Repo      repo   // The GitHub repo where this release lives
	GitTag    string // The specific git tag for this release
	branch    string // If specified, will find HEAD commit
	commitSha string // Specific sha
}

// tag {}:
type tag struct {
	Name       string // The tag name
	ZipBallUrl string // The URL where a ZIP of the release can be downloaded
	TarballUrl string // The URL where a Tarball of the release can be downloaded
	Commit     taggedCommit
}

// taggedCommit {}:
type taggedCommit struct {
	Sha string // The SHA of the commit associated with a given tag
	Url string // The URL to get more commit info
}

// release {}:
// https://developer.github.com/v3/repos/releases/#get-a-release-by-tag-name
type release struct {
	Id         int    // release id
	Url        string // release url
	Name       string // release name (not tag)
	Prerelease bool   // not published?
	TagName    string `json:"tag_name"` // the associated git tag
	Assets     []relAsset
}

// relAsset {}: (release attachment)
// https://developer.github.com/v3/repos/releases/#get-a-release-by-tag-name
type relAsset struct {
	Id   int    // asset id (not release id)
	Url  string // url to retrieve asset
	Name string // asset name
}

func init() {
	rand.Seed(time.Now().UnixNano())

	cl = &http.Client{
		Timeout:   time.Second * 120,
		Transport: http.DefaultTransport,
	}
}

// fetchReleaseTags ():
// returns a list of tags to releases that are:
// i) published
// ii) contain the desired --release-assets
// iii) starts with --tag-prefix if specified
func (o *fetchOpts) fetchReleaseTags(r repo) (tagsForValidRels []string, err error) {

	url := createGitHubRepoUrlForPath(r, "releases")
	resps, err := r.callGitHubApi(url, headers{}, o.verbose)
	if err != nil {
		return tagsForValidRels, err
	}

	// Convert the response body to a byte array

	for _, resp := range resps {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		jsonResp := buf.Bytes()

		// Extract the JSON into our array of gitHubTagsCommitApiResponse's
		var rels []release
		if err := json.Unmarshal(jsonResp, &rels); err != nil {
			return tagsForValidRels, err
		}

		tags := o.filterReleaseTags(rels)

		if len(tags) != 0 {
			tagsForValidRels = append(tagsForValidRels, tags...)
		}

	}

	if len(tagsForValidRels) == 0 {
		return tagsForValidRels, fmt.Errorf("No single release found with all requested assets")
	}

	return
}

// filterReleaseTags ():
func (o *fetchOpts) filterReleaseTags(rels []release) (tags []string) {

	for _, rel := range rels {
		// ... skip if prerelease
		if rel.Prerelease {
			if o.verbose {
				fmt.Printf("... ignoring rel tag %s: prerelease.\n", rel.TagName)
			}
			continue
		}
		// ... skip if release contains fewer assets than number requested
		if len(rel.Assets) < len(o.relAssets) {
			if o.verbose {
				fmt.Printf("... ignoring rel tag %s: not all requested assets.\n", rel.TagName)
			}
			continue
		}

		var relAssetsList []string
		for _, a := range rel.Assets {
			relAssetsList = append(relAssetsList, a.Name)
		}
		// ... skip if desired asset not in list of attached assets
		var missingAsset bool
		for _, a := range o.relAssets {
			if stringInSlice(a, relAssetsList) {
				continue
			} else {
				if o.verbose {
					fmt.Printf("... ignoring rel tag %s: %s not attached.\n", rel.TagName, a)
				}
				missingAsset = true
				break
			}
		}

		if missingAsset {
			continue
		}

		tags = append(tags, rel.TagName)
	}

	return
}

// Fetch all tags from the given GitHub repo
func fetchTags(r repo, verbose bool) ([]string, error) {
	var tagsList []string

	url := createGitHubRepoUrlForPath(r, "tags")
	resps, err := r.callGitHubApi(url, headers{}, verbose)
	if err != nil {
		return tagsList, err
	}

	for _, resp := range resps {
		// ... response to bytes for unmarshalling
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		jsonResp := buf.Bytes()

		// ... unmarshall resp to structs
		var tags []tag
		if err := json.Unmarshal(jsonResp, &tags); err != nil {
			return tagsList, err
		}

		for _, tag := range tags {
			tagsList = append(tagsList, tag.Name)
		}
	}

	if len(tagsList) == 0 {
		return tagsList, fmt.Errorf("No tags found from %s", url)
	}

	return tagsList, err
}

// Convert a URL into a repo struct
func urlToGitHubRepo(url string, token string) (repo, error) {
	var r repo

	regex, regexErr := regexp.Compile("https?://(?:www\\.)?github.com/(.+?)/(.+?)(?:$|\\?|#|/)")
	if regexErr != nil {
		return r, fmt.Errorf("GitHub Repo URL %s is malformed.", url)
	}

	matches := regex.FindStringSubmatch(url)
	if len(matches) != 3 {
		return r, fmt.Errorf("GitHub Repo URL %s could not be parsed correctly", url)
	}

	r = repo{
		Url:   url,
		Owner: matches[1],
		Name:  matches[2],
		Token: token,
		Api:   "https://api.github.com",
	}

	return r, nil
}

// FetchReleaseAsset ():
// Download the release asset with the given id and return its body
func FetchReleaseAsset(r repo, assetId int, destPath string) error {

	url := createGitHubRepoUrlForPath(r, fmt.Sprintf("releases/assets/%d", assetId))

	// ... don't need to use callGitHubApi as that only wraps for pagination.
	resp, _, err := r.apiResp(url, "", headers{"Accept": "application/octet-stream"}, false)
	if err != nil {
		return err
	}

	return writeResponseToDisk(resp, destPath)
}

// GetGitHubReleaseInfo ():
// Get information about the GitHub release with the given tag
func GetGitHubReleaseInfo(r repo, tag string) (release, error) {
	release := release{}

	url := createGitHubRepoUrlForPath(r, fmt.Sprintf("releases/tags/%s", tag))
	resp, _, err := r.apiResp(url, "", headers{}, false)
	if err != nil {
		return release, err
	}

	// Convert the response body to a byte array
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	jsonResp := buf.Bytes()

	err = json.Unmarshal(jsonResp, &release)

	return release, err
}

// Craft a URL for the GitHub repos API of the form repos/:owner/:repo/:path
func createGitHubRepoUrlForPath(r repo, path string) string {
	return fmt.Sprintf("repos/%s/%s/%s", r.Owner, r.Name, path)
}

func retryReq(request *http.Request, url string) (resp *http.Response, err error) {

	attempt := 1
	sleep := time.Second
	for attempt <= 3 {
		resp, err = cl.Do(request)
		if resp != nil && resp.StatusCode < 500 {
			break // success or client err so retry unnecessary
		}

		if attempt < 3 {
			fmt.Printf("Remote failure. Will retry call to %s\n", url)
		}

		attempt++
		jitter := time.Duration(rand.Int63n(int64(sleep)))
		time.Sleep(sleep + jitter/2)
		sleep = sleep * 2

	}

	return
}

func (r repo) apiResp(path string, page string, h headers, v bool) (*http.Response, string, error) {

	var resp *http.Response

	url := fmt.Sprintf("%s/%s?per_page=100&page=%s", r.Api, path, page)

	next := "" // next page of results if any, assume none

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("Failed creating request object for %s: %s", url, err)
	}

	if r.Token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("token %s", r.Token))
	}

	for headerName, headerValue := range h {
		request.Header.Set(headerName, headerValue)
	}

	if v {
		fmt.Printf("... fetching page %s of results from api\n", page)
	}
	resp, err = retryReq(request, url)
	if err != nil { // not checking resp code, only whether http transport succeeded
		return nil, "", err
	}

	if resp.StatusCode != 200 {
		// Convert the resp.Body to a string
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respBody := buf.String()

		// Return err on non-200
		return nil, "", ghApiErr(resp.StatusCode, url, respBody)
	}

	if link, ok := resp.Header["Link"]; ok {
		for _, l := range link {
			if strings.Contains(l, "rel=\"next\"") {
				p, _ := strconv.Atoi(page)
				n := p + 1
				next = strconv.Itoa(n)
			}
		}
	}

	return resp, next, err
}

// callGitHubApi ():
// Call the GitHub API, return HTTP response, and next page number if any
func (r repo) callGitHubApi(path string, h headers, v bool) (resps []*http.Response, err error) {
	page := "1"

	for page != "" {
		resp, n, err := r.apiResp(path, page, h, v)
		if err != nil {
			return resps, err
		}
		page = n
		resps = append(resps, resp)
	}

	return
}

// Write the body of the given HTTP response to disk at the given path
func writeResponseToDisk(resp *http.Response, destPath string) (err error) {
	out, err := os.Create(destPath)
	if err != nil {
		return
	}

	defer out.Close()
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return
}

// ghApiErr : returns a valid `error` obj
func ghApiErr(status int, url string, body string) error {
	var tmpl string
	switch {

	case status == 401:
		tmpl = `
Received an HTTP %d Response when attempting to query the repo.
url: %s

Either your GitHub OAuth Token is invalid, or you don't have access to
the repo with that token. Is the repo private?

http response:
%s
`
	case status == 404:
		tmpl = `
Received an HTTP %d Response when attempting to query the repo.
url: %s

Either the URL does not exist, or you don't have permission to access it.
If the repo is private, you need to set GITHUB_TOKEN (or GITHUB_OAUTH_TOKEN)
in the env before invoking fetch.

http response:
%s
`
	case status >= 500 && status < 600:
		tmpl = `
Received HTTP response %d from GitHub. Is it down?
url: %s

http response:
%s
`
	default:
		tmpl = `
Received non-200 HTTP response %d from GitHub. Could not fulfill request."
url: %s

http response:
%s
`
	} // end switch

	return fmt.Errorf(tmpl, status, url, body)
}
