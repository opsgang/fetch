package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Although other tests besides those in this file require this env var, this init() func will cover all tests.
func init() {
	if os.Getenv("GITHUB_OAUTH_TOKEN") == "" && os.Getenv("GITHUB_TOKEN") == "" {
		fmt.Println("ERROR: Before running these tests, set GITHUB_OAUTH_TOKEN or GITHUB_TOKEN.")
		fmt.Println("See the tests cases to see which GitHub repos the oAuth token needs access to.")
		os.Exit(1)
	}
}

func TestDownloadGitTagZipFile(t *testing.T) {
	t.Parallel()

	cases := []struct {
		repoOwner   string
		repoName    string
		gitTag      string
		githubToken string
	}{
		{"opsgang", "fetch", "v0.1.1", ""},
		{"opsgang", "fetch", "v0.0.2", os.Getenv("GITHUB_OAUTH_TOKEN")},
	}

	for _, tc := range cases {
		gitHubCommit := GitHubCommit{
			Repo: GitHubRepo{
				Owner: tc.repoOwner,
				Name:  tc.repoName,
			},
			GitTag: tc.gitTag,
		}

		zipFilePath, err := downloadGithubZipFile(gitHubCommit, tc.githubToken)
		defer os.RemoveAll(zipFilePath)
		if err != nil {
			t.Fatalf("Failed to download file: %s", err)
		}

		if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
			t.Fatalf("Downloaded file doesn't exist at the expected path of %s", zipFilePath)
		}
	}
}

func TestDownloadGitBranchZipFile(t *testing.T) {
	t.Parallel()

	cases := []struct {
		repoOwner   string
		repoName    string
		branch  string
		githubToken string
	}{
		{"opsgang", "fetch", "enable-fetch-to-pull-from-branch", ""},
		{"opsgang", "fetch", "enable-fetch-to-pull-from-branch", os.Getenv("GITHUB_OAUTH_TOKEN")},
	}

	for _, tc := range cases {
		gitHubCommit := GitHubCommit{
			Repo: GitHubRepo{
				Owner: tc.repoOwner,
				Name:  tc.repoName,
			},
			branch: tc.branch,
		}

		zipFilePath, err := downloadGithubZipFile(gitHubCommit, tc.githubToken)
		defer os.RemoveAll(zipFilePath)
		if err != nil {
			t.Fatalf("Failed to download file: %s", err)
		}

		if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
			t.Fatalf("Downloaded file doesn't exist at the expected path of %s", zipFilePath)
		}
	}
}

func TestDownloadBadGitBranchZipFile(t *testing.T) {
	t.Parallel()

	cases := []struct {
		repoOwner   string
		repoName    string
		branch  string
		githubToken string
	}{
		{"opsgang", "fetch", "branch-that-doesnt-exist", ""},
	}

	for _, tc := range cases {
		gitHubCommit := GitHubCommit{
			Repo: GitHubRepo{
				Owner: tc.repoOwner,
				Name:  tc.repoName,
			},
			branch: tc.branch,
		}

		zipFilePath, err := downloadGithubZipFile(gitHubCommit, tc.githubToken)
		defer os.RemoveAll(zipFilePath)
		if err == nil {
			t.Fatalf("Expected that attempt to download repo %s/%s for branch \"%s\" would fail, but received no error.", tc.repoOwner, tc.repoName, tc.branch)
		}
	}
}

func TestDownloadGitCommitFile(t *testing.T) {
	t.Parallel()

	cases := []struct {
		repoOwner   string
		repoName    string
		commitSha   string
		githubToken string
	}{
		{"opsgang", "fetch", "f5790b465750498bf781169bae74747a6a7b536e", ""},
		{"opsgang", "fetch", "9815bb39119e66c89d5f1c3abeb9d980993ef0a4", ""},
		{"opsgang", "fetch", "f5790b465750498bf781169bae74747a6a7b536e", os.Getenv("GITHUB_OAUTH_TOKEN")},
	}

	for _, tc := range cases {
		gitHubCommit := GitHubCommit{
			Repo: GitHubRepo{
				Owner: tc.repoOwner,
				Name:  tc.repoName,
			},
			commitSha: tc.commitSha,
		}

		zipFilePath, err := downloadGithubZipFile(gitHubCommit, tc.githubToken)
		defer os.RemoveAll(zipFilePath)
		if err != nil {
			t.Fatalf("Failed to download file: %s", err)
		}

		if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
			t.Fatalf("Downloaded file doesn't exist at the expected path of %s", zipFilePath)
		}
	}
}

func TestDownloadBadGitCommitFile(t *testing.T) {
	t.Parallel()

	cases := []struct {
		repoOwner   string
		repoName    string
		commitSha   string
		githubToken string
	}{
		{"opsgang", "fetch", "hello-world", ""},
		{"opsgang", "fetch", "i-am-a-non-existent-commit", ""},
		// remove a single letter from the beginning of an otherwise legit commit sha
		// interestingly, through testing I found that GitHub will attempt to find the right commit sha if you
		// truncate the end of it.
		{"opsgang", "fetch", "7752e7f1df0acbd3c1e61545d5c4d0e87699d84", ""},
	}

	for _, tc := range cases {
		gitHubCommit := GitHubCommit{
			Repo: GitHubRepo{
				Owner: tc.repoOwner,
				Name:  tc.repoName,
			},
			commitSha: tc.commitSha,
		}

		zipFilePath, err := downloadGithubZipFile(gitHubCommit, tc.githubToken)
		defer os.RemoveAll(zipFilePath)
		if err == nil {
			t.Fatalf("Expected that attempt to download repo %s/%s at commmit sha \"%s\" would fail, but received no error.", tc.repoOwner, tc.repoName, tc.commitSha)
		}
	}
}

func TestDownloadZipFileWithBadRepoValues(t *testing.T) {
	t.Parallel()

	cases := []struct {
		repoOwner   string
		repoName    string
		gitTag      string
		githubToken string
	}{
		{"https://github.com/opsgang/fetch/archive/does-not-exist.zip", "MyNameIsWhat", "x.y.z", ""},
	}

	for _, tc := range cases {
		gitHubCommit := GitHubCommit{
			Repo: GitHubRepo{
				Owner: tc.repoOwner,
				Name:  tc.repoName,
			},
			GitTag: tc.gitTag,
		}

		_, err := downloadGithubZipFile(gitHubCommit, tc.githubToken)
		if err == nil && err.errorCode != 500 {
			t.Fatalf("Expected error for bad repo values: %s/%s:%s", tc.repoOwner, tc.repoName, tc.gitTag)
		}
	}
}

func TestExtractFiles(t *testing.T) {
	t.Parallel()

	cases := []struct {
		localFilePath     string
		filePathToExtract string
		expectedNumFiles  int
		nonemptyFiles     []string
	}{
		{"test-fixtures/fetch-test-public-0.0.1.zip", "/", 1, nil},
		{"test-fixtures/fetch-test-public-0.0.1.zip", "/README.md", 1, nil}, // single file as --source-path
		{"test-fixtures/fetch-test-public-0.0.2.zip", "/", 2, nil},
		{"test-fixtures/fetch-test-public-0.0.3.zip", "/", 4, []string{"/README.md"}},
		{"test-fixtures/fetch-test-public-0.0.3.zip", "/folder", 2, nil},
	}

	for _, tc := range cases {
		// Create a temp directory
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %s", err)
		}
		defer os.RemoveAll(tempDir)

		err = extractFiles(tc.localFilePath, tc.filePathToExtract, tempDir)
		if err != nil {
			t.Fatalf("Failed to extract files: %s", err)
		}

		// Count the number of files in the directory
		var numFiles int
		filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				numFiles++
			}
			return nil
		})

		if numFiles != tc.expectedNumFiles {
			t.Fatalf("While extracting %s, expected to find %d file(s), but found %d. Local path = %s", tc.localFilePath, tc.expectedNumFiles, numFiles, tempDir)
		}

		// Ensure that files declared to be non-empty are in fact non-empty
		filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			relativeFilename := strings.TrimPrefix(path, tempDir)

			if !info.IsDir() && stringInSlice(relativeFilename, tc.nonemptyFiles) {
				if info.Size() == 0 {
					t.Fatalf("Expected %s in %s to have non-zero file size, but found file size = %d.\n", relativeFilename, tc.localFilePath, info.Size())
				}
			}
			return nil
		})

	}
}

func TestUnpack(t *testing.T) {
	t.Parallel()

	fixDir := "test-fixtures"
	tmpDirBase := "/tmp/ghfetch-TestUnpack"

	allFiles := []string{
		"dirA",
		"dirA/file",
		"dirB",
		"dirB/dirC",
		"dirB/file.sh",
		"symlink",
	}

	cases := []struct {
		sourceFileName string
		expectedFiles  []string
	}{
		{"packed.tgz", allFiles},
		{"packed.tar.gz", allFiles},
		{"packed.tar", allFiles},
		{"file.gz", []string{"file"}},
		{"dodgygz", []string{"dodgygz"}},
	}

	os.MkdirAll(tmpDirBase, 0755)
	for _, tc := range cases {
		tempDir, err := ioutil.TempDir(tmpDirBase, "")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %s", err)
		}
		defer os.RemoveAll(tempDir)

		sourceFileOriginal := fmt.Sprintf("%s/%s", fixDir, tc.sourceFileName)
		sourceFile := fmt.Sprintf("%s/%s", tmpDirBase, tc.sourceFileName)

		if err := copyFile(sourceFileOriginal, sourceFile); err != nil {
			t.Fatalf("Failed to copy file %s: %s", sourceFileOriginal, err)
		}

		// suppress output from Unpack
		err = Unpack(sourceFile, tempDir)
		if err != nil {
			t.Fatalf("Failed to Unpack files: %s", err)
		}
		// Ensure that files declared to be non-empty are in fact non-empty
		filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			if path != tempDir {
				relativeFilename := strings.TrimPrefix(path, fmt.Sprintf("%s/", tempDir))

				if !stringInSlice(relativeFilename, tc.expectedFiles) {
					t.Fatalf("Unexpected file %s in pack %s.\n", relativeFilename, tc.sourceFileName)
				}
			}
			return nil
		})

	}

	return

}

// TestUntar - check that untarred contents have expected objects and permissions
func TestUntar(t *testing.T) {
	t.Parallel()

	fixDir := "test-fixtures"
	tmpDirBase := "/tmp/ghfetch-TestUntar"

	type fileInfo struct {
		isDir     bool
		isSymLink bool
		isRegular bool
		filePerms string
	}

	dirPerms := "drwxr-xr-x"
	regPerms := "-rw-r--r--"
	exePerms := "-rwxr-xr-x"
	symlinkPerms := "Lrwxrwxrwx"
	r := make(map[string]fileInfo)
	r["dirA"] = fileInfo{isDir: true, filePerms: dirPerms}
	r["dirA/file"] = fileInfo{isRegular: true, filePerms: regPerms}
	r["dirB"] = fileInfo{isDir: true, filePerms: dirPerms}
	r["dirB/dirC"] = fileInfo{isDir: true, filePerms: dirPerms}
	r["dirB/file.sh"] = fileInfo{isRegular: true, filePerms: exePerms}
	r["symlink"] = fileInfo{isSymLink: true, filePerms: symlinkPerms}

	os.MkdirAll(tmpDirBase, 0755)
	tempDir, err := ioutil.TempDir(tmpDirBase, "")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)
	sourceFileName := "packed.tar"
	sourceFileOriginal := fmt.Sprintf("%s/%s", fixDir, sourceFileName)
	sourceFile := fmt.Sprintf("%s/%s", tmpDirBase, sourceFileName)

	if err := copyFile(sourceFileOriginal, sourceFile); err != nil {
		t.Fatalf("Failed to copy file %s: %s", sourceFileOriginal, err)
	}

	err = Untar(sourceFile, tempDir)
	if err != nil {
		t.Fatalf("Failed to Untar files: %s", err)
	}

	// Ensure that files declared to be non-empty are in fact non-empty
	filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if path != tempDir {
			relativePath := strings.TrimPrefix(path, fmt.Sprintf("%s/", tempDir))
			if i, ok := r[relativePath]; !ok {
				t.Fatalf("File %s in tar %s not expected in results!", relativePath, sourceFileName)
			} else {
				if info.IsDir() && !i.isDir {
					t.Fatalf("path %s in tar expected to be a dir", relativePath)
				}
				if info.Mode().IsRegular() && !i.isRegular {
					t.Fatalf("path %s in tar expected to be a regular file", relativePath)
				}
				if info.Mode()&os.ModeSymlink != 0 && !i.isSymLink {
					t.Fatalf("path %s in tar expected to be a symlink", relativePath)
				}
				if fmt.Sprintf("%s", info.Mode()) != i.filePerms {
					t.Fatalf("path %s expected perms %s, not %s", relativePath, i.filePerms, info.Mode())
				}
			}
		}
		return nil
	})
}

// Return ture if the given slice contains the given string
func stringInSlice(s string, slice []string) bool {
	for _, val := range slice {
		if val == s {
			return true
		}
	}
	return false
}

func copyFile(src string, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(dst, data, 0644); err != nil {
		return err
	}
	return nil
}
