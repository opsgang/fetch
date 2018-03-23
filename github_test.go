package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestFetchReleaseTags(t *testing.T) {
	t.Parallel()

	// ... has some published release tags with all requested assets.
	r1 := repo{
		Url:   "https://github.com/foo/bar",
		Owner: "foo",
		Name:  "bar",
		Token: "dummytoken",
		Api:   s.URL,
	}

	// ... has no releases that are both published and have all requested assets.
	r2 := repo{
		Url:   "https://github.com/sna/fu",
		Owner: "sna",
		Name:  "fu",
		Token: "dummytoken",
		Api:   s.URL,
	}

	r3 := repo{
		Url:   "https://github.com/has/no",
		Owner: "has",
		Name:  "no",
		Token: "dummytoken",
		Api:   s.URL,
	}

	var o fetchOpts
	o.relAssets = []string{"foo.tgz", "bar.tgz"}

	// test when there are paginated tags results, and only some have releases
	// ... get results from the stub server for foo/bar
	// ... fiterTags is not responsible for discarding non-semantic version tags,
	// only those that are prereleases, or lack one of more requested attached asset.
	tags, err := o.fetchReleaseTags(r1)
	if err != nil {
		t.Fatalf("error fetching stub releases for foo/bar: err: %s", err)
	} else if !reflect.DeepEqual(tags, relTagsExpected) {
		t.Fatalf("error fetching stub releases for foo/bar: got %#v", tags)
	}

	// test when there are NO valid tags - expect an error, and empty tags
	// ... get results from the stub server for foo/bar
	tags, err = o.fetchReleaseTags(r2)
	if tags != nil {
		t.Fatalf("should have got no go releases back for this test: got %#v", tags)
	}
	if err == nil {
		t.Fatalf("... should have got an error as no releases from sna/fu")
	}

	// test when there are no releases of any kind on repo (api returns empty array)
	// ... get results from the stub server for has/no
	tags, err = o.fetchReleaseTags(r3)
	if tags != nil {
		t.Fatalf("should have got no go releases back for this test: got %#v", tags)
	}
	if err == nil {
		t.Fatalf("... should have got an error as no releases from sna/fu")
	}
}

func TestFetchTagsOnStubApi(t *testing.T) {
	t.Parallel()

	r1 := repo{
		Url:   "https://github.com/foo/bar",
		Owner: "foo",
		Name:  "bar",
		Token: "dummytoken",
		Api:   s.URL,
	}

	r2 := repo{
		Url:   "https://github.com/has/no",
		Owner: "has",
		Name:  "no",
		Token: "dummytoken",
		Api:   s.URL,
	}

	tags, err := fetchTags(r1, false)
	if err != nil {
		t.Fatalf("... should not have erred getting stubbed tags for %s", r1.Url)
	} else if !reflect.DeepEqual(tags, apiTagsExpected) {
		t.Fatalf("... got %#v. Expected %#v", tags, apiTagsExpected)
	}

	// test when no tags returned (empty json array)
	// ... get results from the stub server for has/no
	tags, err = fetchTags(r2, false)
	if tags != nil {
		t.Fatalf("should have got no go tags back for this test: got %#v", tags)
	}
	if err == nil {
		t.Fatalf("... should have got an error as no tags from %s", r2.Url)
	}
}

func TestFetchTagsOnRealRepos(t *testing.T) {
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
		tags, err := fetchTags(r, false)
		if err != nil {
			t.Fatalf("error fetching tags: %s", err)
		}

		if len(tags) != 0 && tc.firstReleaseTag == "" {
			t.Fatalf("expected empty list of tags for repo %s, but got first tag = %s", tc.repoUrl, tags[0])
		}

		if len(tags) == 0 && tc.firstReleaseTag != "" {
			t.Fatalf("expected non-empty list of tags for repo %s, but no tags were found", tc.repoUrl)
		}

		if tags[len(tags)-1] != tc.firstReleaseTag {
			t.Fatalf("error parsing github tags for repo %s. expected first tag = %s, actual = %s", tc.repoUrl, tc.firstReleaseTag, tags[len(tags)-1])
		}

		if tags[0] != tc.lastReleaseTag {
			t.Fatalf("error parsing github tags for repo %s. expected first tag = %s, actual = %s", tc.repoUrl, tc.lastReleaseTag, tags[0])
		}
	}
}

func TestFilterTags(t *testing.T) {
	t.Parallel()

	var o fetchOpts
	o.relAssets = []string{"magic.rb", "wizardry.py", "sourcery.go"}

	allAssets := []relAsset{
		{1, "irrelevant-for-this", "magic.rb"},
		{2, "irrelevant-for-this", "wizardry.py"},
		{3, "irrelevant-for-this", "sourcery.go"},
	}

	missingAsset := []relAsset{
		{1, "irrelevant-for-this", "magic.rb"},
		{2, "irrelevant-for-this", "sourcery.go"},
	}

	wrongAsset := []relAsset{
		{1, "irrelevant-for-this", "magic.rb"},
		{2, "irrelevant-for-this", "magic.c"},
		{2, "irrelevant-for-this", "magic.h"},
	}

	respAllValid := []release{
		{1, "irrelevant-for-this", "MagicFoo1", false, "v1.0.0", allAssets},
		{2, "irrelevant-for-this", "MagicFoo2", false, "v2.0.0", allAssets},
		{3, "irrelevant-for-this", "MagicFoo3", false, "v3.0.0", allAssets},
	}

	respAllPrerelease := []release{
		{1, "irrelevant-for-this", "MagicFoo1", true, "v1.0.0", allAssets},
		{2, "irrelevant-for-this", "MagicFoo2", true, "v2.0.0", allAssets},
		{3, "irrelevant-for-this", "MagicFoo3", true, "v3.0.0", allAssets},
	}

	respTooFewAssets := []release{
		{1, "irrelevant-for-this", "MagicFoo1", false, "v1.0.0", allAssets},
		{2, "irrelevant-for-this", "MagicFoo2", false, "v2.0.0", missingAsset},
		{3, "irrelevant-for-this", "MagicFoo3", false, "v3.0.0", allAssets},
	}

	respWrongAsset := []release{
		{1, "irrelevant-for-this", "MagicFoo1", false, "v1.0.0", allAssets},
		{2, "irrelevant-for-this", "MagicFoo2", false, "v2.0.0", allAssets},
		{3, "irrelevant-for-this", "MagicFoo3", false, "v3.0.0", wrongAsset},
	}

	cases := []struct {
		resps  []release
		result []string
	}{
		{respAllValid, []string{"v1.0.0", "v2.0.0", "v3.0.0"}},
		{respAllPrerelease, nil},
		{respTooFewAssets, []string{"v1.0.0", "v3.0.0"}},
		{respWrongAsset, []string{"v1.0.0", "v2.0.0"}},
	}

	for _, tc := range cases {
		tags := o.filterTags(tc.resps)
		if !reflect.DeepEqual(tags, tc.result) {
			t.Fatalf("tags string did not match for %#v\n\tGot %#v", tc, tags)
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
			t.Fatalf("error extracting url %s into a repo struct: %s", tc.repoUrl, err)
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

	expectedReleaseAsset := relAsset{
		Id:   5354782,
		Url:  "https://api.github.com/repos/opsgang/fetch/releases/assets/5354782",
		Name: "fetch.tgz",
	}

	expectedFetchTestPublicRelease := release{
		Id:         8471364,
		Url:        "https://api.github.com/repos/opsgang/fetch/releases/8471364",
		Name:       "static binary for amd64 linux",
		Prerelease: false,
		TagName:    "v0.1.1",
		Assets:     append([]relAsset{}, expectedReleaseAsset),
	}

	cases := []struct {
		repoUrl   string
		repoToken string
		tag       string
		expected  release
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

	r := repo{
		Url:   "https://github.com/foo/bar",
		Owner: "foo",
		Name:  "bar",
		Token: "dummytoken",
		Api:   s.URL,
	}

	// ... check we get next page back
	resp, next, err := r.apiResp("repos/foo/bar/tags", "1", headers{}, false)

	if err != nil {
		t.Fatalf("did not expect error when calling repos/foo/bar/tags, page 1: %s", err)
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
	_, next, err = r.apiResp("repos/foo/bar/tags", "2", headers{}, false)
	if err != nil {
		t.Fatalf("did not expect error when calling repos/foo/bar/tags, page 1")
	}
	if next != "" {
		t.Fatalf("expected no more results from Link response header")
	}

}

func TestRetryReq(t *testing.T) {
	t.Parallel()

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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
