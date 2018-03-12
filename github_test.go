package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestGetListOfReleasesFromGitHubRepo(t *testing.T) {
	t.Parallel()

	cases := []struct {
		repoUrl          string
		firstReleaseTag  string
		lastReleaseTag   string
		gitHubOAuthToken string
	}{
		// Test on a public repo whose sole purpose is to be a test fixture for this tool
		{"https://github.com/opsgang/fetch", "v0.0.1", "v0.1.1", ""},

		// Private repo equivalent
		{"https://github.com/opsgang/docker_aws_env", "1.0.0", "1.2.1", os.Getenv("GITHUB_OAUTH_TOKEN")},
	}

	for _, tc := range cases {
		r, err := urlToGitHubRepo(tc.repoUrl, tc.gitHubOAuthToken)
		releases, err := FetchTags(r)
		if err != nil {
			t.Fatalf("error fetching releases: %s", err)
		}

		if len(releases) != 0 && tc.firstReleaseTag == "" {
			t.Fatalf("expected empty list of releases for repo %s, but got first release = %s", tc.repoUrl, releases[0])
		}

		if len(releases) == 0 && tc.firstReleaseTag != "" {
			t.Fatalf("expected non-empty list of releases for repo %s, but no releases were found", tc.repoUrl)
		}

		if releases[len(releases)-1] != tc.firstReleaseTag {
			t.Fatalf("error parsing github releases for repo %s. expected first release = %s, actual = %s", tc.repoUrl, tc.firstReleaseTag, releases[len(releases)-1])
		}

		if releases[0] != tc.lastReleaseTag {
			t.Fatalf("error parsing github releases for repo %s. expected first release = %s, actual = %s", tc.repoUrl, tc.lastReleaseTag, releases[0])
		}
	}
}

func TestUrlToGitHubRepo(t *testing.T) {
	t.Parallel()

	cases := []struct {
		repoUrl string
		owner   string
		name    string
		token   string
	}{
		{"https://github.com/brikis98/ping-play", "brikis98", "ping-play", ""},
		{"http://github.com/brikis98/ping-play", "brikis98", "ping-play", ""},
		{"https://github.com/gruntwork-io/script-modules", "gruntwork-io", "script-modules", ""},
		{"http://github.com/gruntwork-io/script-modules", "gruntwork-io", "script-modules", ""},
		{"http://www.github.com/gruntwork-io/script-modules", "gruntwork-io", "script-modules", ""},
		{"http://www.github.com/gruntwork-io/script-modules/", "gruntwork-io", "script-modules", ""},
		{"http://www.github.com/gruntwork-io/script-modules?foo=bar", "gruntwork-io", "script-modules", "token"},
		{"http://www.github.com/gruntwork-io/script-modules?foo=bar&foo=baz", "gruntwork-io", "script-modules", "token"},
	}

	for _, tc := range cases {
		repo, err := urlToGitHubRepo(tc.repoUrl, tc.token)
		if err != nil {
			t.Fatalf("error extracting url %s into a GitHubRepo struct: %s", tc.repoUrl, err)
		}

		if repo.Owner != tc.owner {
			t.Fatalf("while extracting %s, expected owner %s, received %s", tc.repoUrl, tc.owner, repo.Owner)
		}

		if repo.Name != tc.name {
			t.Fatalf("while extracting %s, expected name %s, received %s", tc.repoUrl, tc.name, repo.Name)
		}

		if repo.Url != tc.repoUrl {
			t.Fatalf("while extracting %s, expected url %s, received %s", tc.repoUrl, tc.repoUrl, repo.Url)
		}

		if repo.Token != tc.token {
			t.Fatalf("while extracting %s, expected token %s, received %s", tc.repoUrl, tc.token, repo.Token)
		}
	}
}

func TestParseUrlThrowsErrorOnMalformedUrl(t *testing.T) {
	t.Parallel()

	cases := []struct {
		repoUrl string
	}{
		{"https://githubb.com/brikis98/ping-play"},
		{"github.com/brikis98/ping-play"},
		{"curl://github.com/brikis98/ping-play"},
	}

	for _, tc := range cases {
		_, err := urlToGitHubRepo(tc.repoUrl, "")
		if err == nil {
			t.Fatalf("Expected error on malformed url %s, but no error was received.", tc.repoUrl)
		}
	}
}

func TestGetGitHubReleaseInfo(t *testing.T) {
	t.Parallel()

	token := os.Getenv("GITHUB_OAUTH_TOKEN")

	expectedReleaseAsset := GitHubReleaseAsset{
		Id:   5354782,
		Url:  "https://api.github.com/repos/opsgang/fetch/releases/assets/5354782",
		Name: "fetch.tgz",
	}

	expectedFetchTestPublicRelease := GitHubReleaseApiResponse{
		Id:         8471364,
		Url:        "https://api.github.com/repos/opsgang/fetch/releases/8471364",
		Name:       "static binary for amd64 linux",
		Prerelease: false,
		Tag_name:   "v0.1.1",
		Assets:     append([]GitHubReleaseAsset{}, expectedReleaseAsset),
	}

	cases := []struct {
		repoUrl   string
		repoToken string
		tag       string
		expected  GitHubReleaseApiResponse
	}{
		{"https://github.com/opsgang/fetch", token, "v0.1.1", expectedFetchTestPublicRelease},
	}

	for _, tc := range cases {
		repo, err := urlToGitHubRepo(tc.repoUrl, tc.repoToken)
		if err != nil {
			t.Fatalf("Failed to parse %s into GitHub URL due to error: %s", tc.repoUrl, err.Error())
		}

		resp, err := GetGitHubReleaseInfo(repo, tc.tag)
		if err != nil {
			t.Fatalf("Failed to fetch GitHub release info for repo %s due to error: %s", tc.repoToken, err.Error())
		}

		if !reflect.DeepEqual(tc.expected, resp) {
			t.Fatalf("Expected GitHub release %v but got GitHub release %v", tc.expected, resp)
		}
	}
}

func TestFetchReleaseAsset(t *testing.T) {
	t.Parallel()

	token := os.Getenv("GITHUB_OAUTH_TOKEN")

	cases := []struct {
		repoUrl   string
		repoToken string
		tag       string
		assetId   int
	}{
		{"https://github.com/opsgang/fetch", token, "v0.1.1", 5354782},
	}

	for _, tc := range cases {
		repo, err := urlToGitHubRepo(tc.repoUrl, tc.repoToken)
		if err != nil {
			t.Fatalf("Failed to parse %s into GitHub URL due to error: %s", tc.repoUrl, err.Error())
		}

		tmpFile, tmpErr := ioutil.TempFile("", "test-download-release-asset")
		if tmpErr != nil {
			t.Fatalf("Failed to create temp file due to error: %s", tmpErr.Error())
		}

		if err := FetchReleaseAsset(repo, tc.assetId, tmpFile.Name()); err != nil {
			t.Fatalf("Failed to download asset %d to %s from GitHub URL %s due to error: %s", tc.assetId, tmpFile.Name(), tc.repoUrl, err.Error())
		}

		defer os.Remove(tmpFile.Name())

		if !fileExists(tmpFile.Name()) {
			t.Fatalf("Got no errors downloading asset %d to %s from GitHub URL %s, but %s does not exist!", tc.assetId, tmpFile.Name(), tc.repoUrl, tmpFile.Name())
		}
	}
}

func TestApiResp(t *testing.T) {
	t.Parallel()

	s := apiStub()
	defer s.Close()

	r := GitHubRepo{
		Url:   "https://github.com/foo/bar",
		Owner: "foo",
		Name:  "bar",
		Token: "dummytoken",
		Api:   s.URL,
	}

	// ... check we get next page back
	resp, next, err := r.apiResp("foo/bar/tags", "1", headers{})

	if err != nil {
		t.Fatalf("did not expect error when calling /foo/bar/tags, page 1")
	}
	if next != "2" {
		t.Fatalf("expected 2nd page of results from Link response header")
	}

	// ... check auth header correct format
	if v, ok := resp.Header["X-Authorization"]; !ok {
		t.Fatalf("auth header should have been copied back in response for test")
	} else {
		if v[0] != "token dummytoken" {
			t.Fatalf("auth header should have value 'token dummytoken', not '%s'", v[0])
		}
	}

	// ... check no further pages after last
	_, next, err = r.apiResp("foo/bar/tags", "2", headers{})
	if err != nil {
		t.Fatalf("did not expect error when calling /foo/bar/tags, page 1")
	}
	if next != "" {
		t.Fatalf("expected no more results from Link response header")
	}

}

func TestRetryReq(t *testing.T) {
	t.Parallel()

	s := apiStub()
	defer s.Close()

	// ... test we are retrying (check the test counter)
	failTwice := fmt.Sprintf("%s/fail/twice", s.URL)
	request, err := http.NewRequest("GET", failTwice, nil)
	resp, err := retryReq(request, failTwice)
	if err != nil {
		t.Fatalf("/fail/twice should not err , not %s", err)
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respBody := buf.String()
		if respBody != "counter: 2" {
			t.Fatalf("/fail/twice should succeed when counter is 2, got %s", respBody)
		}
	}

	// ... after retries exceeded, should fail.
	failAlways := fmt.Sprintf("%s/fail/always", s.URL)
	request, err = http.NewRequest("GET", failAlways, nil)
	resp, err = retryReq(request, failAlways)
	if err != nil || resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("/fail/always should have failed with resp code %d", http.StatusInternalServerError)
	}
}

func apiStub() *httptest.Server {
	var resp string
	var counter = 0 // used for keeping track of retries

	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if v, ok := r.Header["Authorization"]; ok {
					w.Header().Set("X-Authorization", v[0])
				}
				switch r.RequestURI {

				case "/foo/bar/tags?per_page=100&page=1":
					w.Header().Set("Link", apiTagsPage1Link)
					resp = apiTagsPage1

				case "/foo/bar/tags?per_page=100&page=2":
					w.Header().Set("Link", apiTagsPage2Link)
					resp = apiTagsPage2

				case "/fail/twice":
					if counter == 2 {
						resp = fmt.Sprintf("counter: %d", counter)
					} else {
						counter++
						http.Error(w, "Remote failure", http.StatusBadGateway)
						return
					}

				case "/fail/always":
					http.Error(w, "Remote failure", http.StatusInternalServerError)
					return

				default:
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
				w.Write([]byte(resp))
			},
		),
	)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
