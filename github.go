package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

type GitHubRepo struct {
	Url   string // The URL of the GitHub repo
	Owner string // The GitHub account name under which the repo exists
	Name  string // The GitHub repo name
	Token string // The personal access token to access this repo (if it's a private repo)
}

// Represents a specific git commit.
// Note that code using GitHub Commit should respect the following hierarchy:
// - commitSha > branch > GitTag
// - Example: GitTag and branch are both specified; use the GitTag
// - Example: GitTag and commitSha are both specified; use the commitSha
// - Example: branch alone is specified; use branch
type GitHubCommit struct {
	Repo      GitHubRepo // The GitHub repo where this release lives
	GitTag    string     // The specific git tag for this release
	branch    string     // If specified, will find HEAD commit
	commitSha string     // Specific sha
}

// Modeled directly after the api.github.com response
type GitHubTagsApiResponse struct {
	Name       string // The tag name
	ZipBallUrl string // The URL where a ZIP of the release can be downloaded
	TarballUrl string // The URL where a Tarball of the release can be downloaded
	Commit     GitHubTagsCommitApiResponse
}

// Modeled directly after the api.github.com response
type GitHubTagsCommitApiResponse struct {
	Sha string // The SHA of the commit associated with a given tag
	Url string // The URL at which additional API information can be found for the given commit
}

// Modeled directly after the api.github.com response (but only includes the fields we care about). For more info, see:
// https://developer.github.com/v3/repos/releases/#get-a-release-by-tag-name
type GitHubReleaseApiResponse struct {
	Id         int
	Url        string
	Name       string
	Prerelease bool
	Tag_name   string
	Assets     []GitHubReleaseAsset
}

// The "assets" portion of the GitHubReleaseApiResponse. Modeled directly after the api.github.com response (but only
// includes the fields we care about). For more info, see:
// https://developer.github.com/v3/repos/releases/#get-a-release-by-tag-name
type GitHubReleaseAsset struct {
	Id   int
	Url  string
	Name string
}

/*
fetchReleaseTags ()
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
	url := createGitHubRepoUrlForPath(repo, "releases?per_page=100")
	resp, err := repo.callGitHubApi(url, map[string]string{})
	if err != nil {
		return tagsString, err
	}

	// Convert the response body to a byte array
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

		// ... compare rel.Name
		fmt.Printf("... adding %s for consideration\n",rel.Tag_name)
		tagsString = append(tagsString, rel.Tag_name)
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

	url := createGitHubRepoUrlForPath(repo, "tags?per_page=100")
	resp, err := repo.callGitHubApi(url, map[string]string{})
	if err != nil {
		return tagsString, err
	}

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
	resp, err := repo.callGitHubApi(url, map[string]string{"Accept": "application/octet-stream"})
	if err != nil {
		return err
	}

	return writeResponseToDisk(resp, destPath)
}

// Get information about the GitHub release with the given tag
func GetGitHubReleaseInfo(repo GitHubRepo, tag string) (GitHubReleaseApiResponse, error) {
	release := GitHubReleaseApiResponse{}

	url := createGitHubRepoUrlForPath(repo, fmt.Sprintf("releases/tags/%s", tag))
	resp, err := repo.callGitHubApi(url, map[string]string{})
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

// Call the GitHub API at the given path and return the HTTP response
func (r GitHubRepo) callGitHubApi(path string, headers map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("https://api.github.com/%s", path)
	httpClient := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed creating request object for %s: %s", url, err)
	}

	if r.Token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("token %s", r.Token))
	}

	for headerName, headerValue := range headers {
		request.Header.Set(headerName, headerValue)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		// Convert the resp.Body to a string
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respBody := buf.String()

		// Return err on non-200
		return nil, ghApiErr(resp.StatusCode, url, respBody)
	}

	return resp, nil
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
