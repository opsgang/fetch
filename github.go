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
	"strings"
	"strconv"
	"time"
)

type GitHubRepo struct {
	Url   string // The URL of the GitHub repo
	Owner string // The GitHub account name under which the repo exists
	Name  string // The GitHub repo name
	Token string // The personal access token to access this repo (if it's a private repo)
}

// TODO: Client Should have keep-alives and timeouts set.
// cl : our single instance of http.Client to be reused throughout.
var cl http.Client

type headers map[string]string

/*
Hierarchy:
* commitSha > branch > GitTag
* Example: GitTag and branch are both specified; use the GitTag
* Example: GitTag and commitSha are both specified; use the commitSha
* Example: branch alone is specified; use branch
*/

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
}

/*
fetchReleaseTags ():
returns a list of tags to releases that are:
i) published
ii) contain the desired --release-assets
iii) starts with --tag-prefix if specified
*/
func (o *fetchOpts) fetchReleaseTags() ([]string, error) {
	var tagsString []string

	repo, err := ParseUrlIntoGitHubRepo(o.repoUrl, o.githubToken)
	if err != nil {
		return tagsString, err
	}

	// TODO - abstract and iterate over pages using header data
	url := createGitHubRepoUrlForPath(repo, "releases")
	resps, err := repo.callGitHubApi(url, map[string]string{})
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
			var missing_asset bool
			for _, a := range o.ReleaseAssets {
				if stringInSlice(a, relAssetsList) {
					continue
				} else {
					fmt.Printf("... ignoring rel tag %s: %s not attached.\n", rel.Tag_name, a)
					missing_asset = true
					break
				}
			}

			if missing_asset {
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
func FetchTags(githubRepoUrl string, githubToken string) ([]string, error) {
	var tagsString []string

	repo, err := ParseUrlIntoGitHubRepo(githubRepoUrl, githubToken)
	if err != nil {
		return tagsString, err
	}

	url := createGitHubRepoUrlForPath(repo, "tags")
	resps, err := repo.callGitHubApi(url, map[string]string{})
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
func ParseUrlIntoGitHubRepo(url string, token string) (GitHubRepo, error) {
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
	}

	return gitHubRepo, nil
}

// Download the release asset with the given id and return its body
func FetchReleaseAsset(repo GitHubRepo, assetId int, destPath string) error {

	url := createGitHubRepoUrlForPath(repo, fmt.Sprintf("releases/assets/%d", assetId))

	// ... don't need to use callGitHubApi as that only wraps for pagination.
	resp, _, err := repo.apiResp(url, "", map[string]string{"Accept": "application/octet-stream"})
	if err != nil {
		return err
	}

	return writeResponseToDisk(resp, destPath)
}

// Get information about the GitHub release with the given tag
func GetGitHubReleaseInfo(repo GitHubRepo, tag string) (GitHubReleaseApiResponse, error) {
	release := GitHubReleaseApiResponse{}

	url := createGitHubRepoUrlForPath(repo, fmt.Sprintf("releases/tags/%s", tag))
	resp, _, err := repo.apiResp(url, "", map[string]string{})
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
func createGitHubRepoUrlForPath(repo GitHubRepo, path string) string {
	return fmt.Sprintf("repos/%s/%s/%s", repo.Owner, repo.Name, path)
}

/*
callGitHubApi ():
Call the GitHub API, return HTTP response, and next page number if any
*/
func (r GitHubRepo) apiResp(path string, page string, h headers) (*http.Response, string, error) {

	var resp *http.Response

	url := fmt.Sprintf("https://api.github.com/%s?per_page=100&page=%s", path, page)

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

	attempt := 1
	sleep := time.Second
	for attempt <= 3 {
		resp, err = cl.Do(request)
		if resp != nil && resp.StatusCode < 500 {
			break // success or an http client err (retry pointless)
		}

		attempt++
		jitter := time.Duration(rand.Int63n(int64(sleep)))
		time.Sleep(sleep + jitter/2)
		sleep = sleep * 2

		if attempt <3 {
			fmt.Printf("Remote failure. Will retry call to %s\n", url)
		}
	}
	if err != nil {
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

func (r GitHubRepo) callGitHubApi(path string, h headers) ([]*http.Response, error) {
	var resps []*http.Response
	var err error
	page := "1"

	for page != "" {
		resp, n, err := r.apiResp(path, page, h)
		if err !=nil {
			return resps, err
		}
		page = n
		resps = append(resps, resp)
	}

	return resps, err
}

// Write the body of the given HTTP response to disk at the given path
func writeResponseToDisk(resp *http.Response, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}

	defer out.Close()
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
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
	case status >= 500 && status  < 600:
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

