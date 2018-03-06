package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"gopkg.in/h2non/filetype.v1"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const FAILED_ZIPBALL_DOWNLOAD = `
Failed to download file at the url %s.
Received HTTP Response %d.
`

const CONTENT_TYPE_NOT_ZIP = `
Failed to download file at the url %s.
Expected HTTP Response's "Content-Type" header to be "application/zip", but was "%s"
`

// Download the zip file at the given URL to a temporary local directory.
// Returns the absolute path to the downloaded zip file.
// IMPORTANT: You must call "defer os.RemoveAll(dir)" in the calling function when done with the downloaded zip file!
func getSrcZip(gitHubCommit GitHubCommit, gitHubToken string) (string, int, error) {

	var zipFilePath string
	var rStatus int

	// Create a temp directory
	// Note that ioutil.TempDir has a peculiar interface. We need not specify any meaningful values to achieve our
	// goal of getting a temporary directory.
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return zipFilePath, rStatus, err
	}

	// Download the zip file, possibly using the GitHub oAuth Token
	httpClient := &http.Client{}
	req, err := gitHubZipRequest(gitHubCommit, gitHubToken)
	if err != nil {
		return zipFilePath, rStatus, err
	}

	resp, err := httpClient.Do(req)

	rStatus = resp.StatusCode
	if err != nil {
		return zipFilePath, rStatus, err
	}

	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf(FAILED_ZIPBALL_DOWNLOAD, req.URL.String(), resp.StatusCode)
		return zipFilePath, rStatus, errors.New(errMsg)
	}

	if resp.Header.Get("Content-Type") != "application/zip" {
		errMsg := fmt.Sprintf(CONTENT_TYPE_NOT_ZIP, req.URL.String(), resp.Header.Get("Content-Type"))
		return zipFilePath, rStatus, errors.New(errMsg)
	}

	// Copy the contents of the downloaded file to our empty file
	respBodyBuffer := new(bytes.Buffer)
	respBodyBuffer.ReadFrom(resp.Body)
	err = ioutil.WriteFile(filepath.Join(tempDir, "repo.zip"), respBodyBuffer.Bytes(), 0644)
	if err != nil {
		return zipFilePath, rStatus, err
	}

	zipFilePath = filepath.Join(tempDir, "repo.zip")

	return zipFilePath, rStatus, err
}

// Decompress the file at zipFileAbsPath and move only those files under filesToExtractFromZipPath to localPath
func extractFiles(zipFilePath, filesToExtractFromZipPath, localPath string) error {

	// Open the zip file for reading.
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	// pathPrefix will be stripped from source path before copying the file to localPath
	// E.g. full path = fetch-test-public-0.0.3/folder/file1.txt
	//      path prefix = fetch-test-public-0.0.3
	//      file that will eventually get written = <localPath>/folder/file1.txt

	// By convention, the first file in the zip file is the top-level directory
	pathPrefix := r.File[0].Name

	// Add the path from which we will extract files to the path prefix so we can exclude the appropriate files
	// unless you only want a single file from the zip ...
	pathPrefix = filepath.Join(pathPrefix, filesToExtractFromZipPath)

	os.MkdirAll(localPath, 0755)
	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {

		// If the given file is in the filesToExtractFromZipPath, proceed
		if strings.Index(f.Name, pathPrefix) == 0 {

			// When source-path is a directory, we want to drop
			// the contents in to the local download dir with out
			// the source-path portion of the path ...
			trimmedName := strings.TrimPrefix(f.Name, pathPrefix)
			// ... but if --source-path is a single file the file name and
			// path prefix are the same so we want just the base file name.
			if f.Name == pathPrefix {
				trimmedName = filepath.Base(f.Name)
			}
			// when source-path is a single file, the trimmed name is empty.
			if f.FileInfo().IsDir() {
				// Create a directory
				os.MkdirAll(filepath.Join(localPath, trimmedName), 0777)
			} else {
				// Read the file into a byte array
				readCloser, err := f.Open()
				if err != nil {
					return fmt.Errorf("Failed to open file %s: %s", f.Name, err)
				}

				byteArray, err := ioutil.ReadAll(readCloser)
				if err != nil {
					return fmt.Errorf("Failed to read file %s: %s", f.Name, err)
				}

				// Write the file
				err = ioutil.WriteFile(filepath.Join(localPath, trimmedName), byteArray, 0644)
				if err != nil {
					return fmt.Errorf("Failed to write file: %s", err)
				}
			}
		}
	}

	return nil
}

func unpack(sourceFileName, destDir string) (err error) {
	if err != nil {
		return err
	}
	fileExt, err := detectFileType(sourceFileName)
	if err != nil {
		return fmt.Errorf("Error detecting filetype of %s: %s", sourceFileName, err)
	}

	switch fileExt {
	case "gz":
		if err = gunzip(sourceFileName, destDir); err != nil {
			return err
		}
	case "tar":
		if err = untar(sourceFileName, destDir); err != nil {
			return err
		}
	}
	return nil
}

// detectFileType : we only care if the asset if a tar or gzip
// ... otherwise we deliver as is.
func detectFileType(source string) (fileExt string, err error) {
	buf, err := ioutil.ReadFile(source)
	if err != nil {
		err = fmt.Errorf("Failed to read file %s to get mimetype: %s", source, err)
	} else {
		kind, unknown := filetype.Match(buf)
		if unknown == nil {
			fileExt = kind.Extension
		}
	}
	return
}

// gunzip : Remember, a gzip will only contain a single file
func gunzip(sourceFileName, destDir string) error {
	reader, err := os.Open(sourceFileName)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	// need to provide a full path for location of ungzipped file
	gunZipped := filepath.Join(destDir, fmt.Sprintf("%s.gunzipped", filepath.Base(sourceFileName)))
	fmt.Printf("Gunzipping %s\n", sourceFileName)
	writer, err := os.Create(gunZipped)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)

	newSource := strings.TrimSuffix(gunZipped, ".gunzipped")
	newSource = strings.TrimSuffix(newSource, ".gz")
	if strings.HasSuffix(newSource, ".tgz") {
		newSource = strings.Replace(newSource, ".tgz", ".tar", 1)
	}

	if err = os.Rename(gunZipped, newSource); err != nil {
		return err
	}

	// now untar if needed
	if err = unpack(newSource, destDir); err != nil {
		return err
	}

	// delete source and gunzipped
	if warning := os.Remove(sourceFileName); warning != nil {
		fmt.Printf("Could not remove intermediary file %s\n", sourceFileName)
	}
	return err
}

// untar : untars arg1 archive in to arg2 dir
func untar(sourceFileName, destDir string) error {
	fmt.Printf("Untarring %s\n", sourceFileName)
	reader, err := os.Open(sourceFileName)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(destDir, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		if info.Mode()&os.ModeSymlink != 0 {
			if err := os.Symlink(header.Linkname, path); err != nil {
				return err
			} else {
				continue
			}
		}

		if info.Mode().IsRegular() {
			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}
		}
	}
	if warning := os.Remove(sourceFileName); warning != nil {
		fmt.Printf("Could not remove intermediary file %s\n", sourceFileName)
	}
	return nil
}

// gitHubZipRequest : returns HTTP request for zipball
// Sha trumps branch which trumps tag.
func gitHubZipRequest(gitHubCommit GitHubCommit, gitHubToken string) (*http.Request, error) {
	var request *http.Request

	// This represents either a commit, branch, or git tag
	var gitRef string
	if gitHubCommit.commitSha != "" {
		gitRef = gitHubCommit.commitSha
	} else if gitHubCommit.branch != "" {
		gitRef = gitHubCommit.branch
	} else if gitHubCommit.GitTag != "" {
		gitRef = gitHubCommit.GitTag
	} else {
		msg := "Neither a commitSha nor a GitTag nor a branch were specified " +
			"so impossible to identify a specific commit to download."
		return request, fmt.Errorf(msg)
	}

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/zipball/%s",
		gitHubCommit.Repo.Owner,
		gitHubCommit.Repo.Name,
		gitRef,
	)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return request, err
	}

	if gitHubToken != "" {
		request.Header.Set("Authorization", fmt.Sprintf("token %s", gitHubToken))
	}

	return request, nil
}
