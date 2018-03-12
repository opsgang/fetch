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

type GitHubRepo struct {
	Url   string // The URL of the GitHub repo
	Owner string // The GitHub account name under which the repo exists
	Name  string // The GitHub repo name
	Token string // The personal access token to access this repo (if it's a private repo)
	Api   string // https://api.github.com (or stub server - see github_test.go)
}

// cl : our single instance of http.Client to be reused throughout.
var cl *http.Client

type headers map[string]string

// Hierarchy:
// * commitSha > branch > GitTag
// * Example: GitTag and branch are both specified; use the GitTag
// * Example: GitTag and commitSha are both specified; use the commitSha
// * Example: branch alone is specified; use branch

// GitHubCommit {}:
// A specific git commit.
type GitHubCommit struct {
	Repo      GitHubRepo // The GitHub repo where this release lives
	GitTag    string     // The specific git tag for this release
	branch    string     // If specified, will find HEAD commit
	commitSha string     // Specific sha
}

// GitHubTagsApiResponse {}:
type GitHubTagsApiResponse struct {
	Name       string // The tag name
	ZipBallUrl string // The URL where a ZIP of the release can be downloaded
	TarballUrl string // The URL where a Tarball of the release can be downloaded
	Commit     GitHubTagsCommitApiResponse
}

// GitHubTagsCommitApiResponse {}:
type GitHubTagsCommitApiResponse struct {
	Sha string // The SHA of the commit associated with a given tag
	Url string // The URL to get more commit info
}

// GitHubReleaseApiResponse {}:
// https://developer.github.com/v3/repos/releases/#get-a-release-by-tag-name
type GitHubReleaseApiResponse struct {
	Id         int    // release id
	Url        string // release url
	Name       string // release name (not tag)
	Prerelease bool   // not published?
	Tag_name   string // the associated git tag
	Assets     []GitHubReleaseAsset
}

// GitHubReleaseAsset {}: (release attachment)
// https://developer.github.com/v3/repos/releases/#get-a-release-by-tag-name
type GitHubReleaseAsset struct {
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
func (o *fetchOpts) fetchReleaseTags(repo GitHubRepo) (tagsString []string, err error) {

	url := createGitHubRepoUrlForPath(repo, "releases")
	resps, err := repo.callGitHubApi(url, headers{})
	if err != nil {
		return tagsString, err
	}

	// Convert the response body to a byte array

	for _, resp := range resps {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		jsonResp := buf.Bytes()

		// Extract the JSON into our array of gitHubTagsCommitApiResponse's
		var rels []GitHubReleaseApiResponse
		if err := json.Unmarshal(jsonResp, &rels); err != nil {
			return tagsString, err
		}

		for _, rel := range rels {
			// ... skip if prerelease
			if rel.Prerelease {
				fmt.Printf("... ignoring rel tag %s: prerelease.\n", rel.Tag_name)
				continue
			}
			// ... skip if release contains fewer assets than number requested
			if len(rel.Assets) < len(o.ReleaseAssets) {
				fmt.Printf("... ignoring rel tag %s: not all requested assets.\n", rel.Tag_name)
				continue
			}

			var relAssetsList []string
			for _, a := range rel.Assets {
				relAssetsList = append(relAssetsList, a.Name)
			}
			// ... skip if desired asset not in list of attached assets
			var missingAsset bool
			for _, a := range o.ReleaseAssets {
				if stringInSlice(a, relAssetsList) {
					continue
				} else {
					fmt.Printf("... ignoring rel tag %s: %s not attached.\n", rel.Tag_name, a)
					missingAsset = true
					break
				}
			}

			if missingAsset {
				continue
			}

			tagsString = append(tagsString, rel.Tag_name)
		}
	}

	if len(tagsString) == 0 {
		return tagsString, fmt.Errorf("No single release found with all requested assets")
	}
	return tagsString, nil
}

// Fetch all tags from the given GitHub repo
func FetchTags(r GitHubRepo) ([]string, error) {
	var tagsString []string

	url := createGitHubRepoUrlForPath(r, "tags")
	resps, err := r.callGitHubApi(url, headers{})
	if err != nil {
		return tagsString, err
	}

	for _, resp := range resps {
		// Convert the response body to a byte array
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		jsonResp := buf.Bytes()

		// Extract the JSON into our array of gitHubTagsCommitApiResponse's
		var tags []GitHubTagsApiResponse
		if err := json.Unmarshal(jsonResp, &tags); err != nil {
			return tagsString, err
		}

		for _, tag := range tags {
			tagsString = append(tagsString, tag.Name)
		}
	}

	return tagsString, nil
}

// Convert a URL into a GitHubRepo struct
func urlToGitHubRepo(url string, token string) (GitHubRepo, error) {
	var gitHubRepo GitHubRepo

	regex, regexErr := regexp.Compile("https?://(?:www\\.)?github.com/(.+?)/(.+?)(?:$|\\?|#|/)")
	if regexErr != nil {
		return gitHubRepo, fmt.Errorf("GitHub Repo URL %s is malformed.", url)
	}

	matches := regex.FindStringSubmatch(url)
	if len(matches) != 3 {
		return gitHubRepo, fmt.Errorf("GitHub Repo URL %s could not be parsed correctly", url)
	}

	gitHubRepo = GitHubRepo{
		Url:   url,
		Owner: matches[1],
		Name:  matches[2],
		Token: token,
		Api:   "https://api.github.com",
	}

	return gitHubRepo, nil
}

// Download the release asset with the given id and return its body
func FetchReleaseAsset(r GitHubRepo, assetId int, destPath string) error {

	url := createGitHubRepoUrlForPath(r, fmt.Sprintf("releases/assets/%d", assetId))

	// ... don't need to use callGitHubApi as that only wraps for pagination.
	resp, _, err := r.apiResp(url, "", headers{"Accept": "application/octet-stream"})
	if err != nil {
		return err
	}

	return writeResponseToDisk(resp, destPath)
}

// Get information about the GitHub release with the given tag
func GetGitHubReleaseInfo(r GitHubRepo, tag string) (GitHubReleaseApiResponse, error) {
	release := GitHubReleaseApiResponse{}

	url := createGitHubRepoUrlForPath(r, fmt.Sprintf("releases/tags/%s", tag))
	resp, _, err := r.apiResp(url, "", headers{})
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
func createGitHubRepoUrlForPath(r GitHubRepo, path string) string {
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

func (r GitHubRepo) apiResp(path string, page string, h headers) (*http.Response, string, error) {

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
func (r GitHubRepo) callGitHubApi(path string, h headers) (resps []*http.Response, err error) {
	page := "1"

	for page != "" {
		resp, n, err := r.apiResp(path, page, h)
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
