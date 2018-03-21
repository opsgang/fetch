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

// FAILED_ZIPBALL_DOWNLOAD : err msg tmpl
const FAILED_ZIPBALL_DOWNLOAD = `
Failed to download file at the url %s.
Received HTTP Response %d.
`

// CONTENT_TYPE_NOT_ZIP : err msg tmpl
const CONTENT_TYPE_NOT_ZIP = `
Failed to download file at the url %s.
Expected HTTP Response's "Content-Type" header to be "application/zip", but was "%s"
`

// getSrcZip ():
// Download the zip file from url to temp dir.
// Returns the absolute path to zip, http resp code, and any err.
func getSrcZip(c commit, gitHubToken string) (string, int, error) {

	var zipFilePath string
	var rStatus int

	// temp dir used for downloading and unpacking zip before copying files.
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return zipFilePath, rStatus, err
	}

	// Download the zip file, possibly using the GitHub oAuth Token
	httpClient := cl
	req, err := gitHubZipRequest(c, gitHubToken)
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

// extractFiles ():
// Decompress zip and filter source for files indicated by --from-path
func extractFiles(zipFilePath, filesToExtractFromZipPath, localPath string) error {

	// Open the zip file for reading.
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	// IMPORTANT - we expect and ignore the parent dir in the zip:
	// GitHub src zip files contain the src of the repo with in a parent folder
	// named the same as the zip file (with out the .zip extension).
	// We ignore this dir, and filter fromPaths past that point.

	// pathPrefix will be stripped from source path before copying the file to localPath
	// E.g. full path = fetch-test-public-0.0.3/folder/file1.txt
	//      path prefix = fetch-test-public-0.0.3
	//      file that will eventually get written = <localPath>/folder/file1.txt

	// By convention, the first file in the zip file is the top-level directory
	pathPrefix := r.File[0].Name

	// Add the path from which we will extract files to, the path prefix so we can exclude
	// the appropriate files unless only a single file was requested.
	pathPrefix = filepath.Join(pathPrefix, filesToExtractFromZipPath)

	os.MkdirAll(localPath, 0755)
	// ... write out zipped files (from fromPath root)
	for _, f := range r.File {

		// If the given file is in the filesToExtractFromZipPath, proceed
		if strings.Index(f.Name, pathPrefix) == 0 {

			// When from-path is a directory, we want to drop the contents in to the local
			// download dir with out the source-path portion of the path ...
			trimmedName := strings.TrimPrefix(f.Name, pathPrefix)

			// ... but if --from-path is a single file the file name and
			// path prefix are the same so we want just the base file name.
			if f.Name == pathPrefix {
				trimmedName = filepath.Base(f.Name)
			}
			fi := f.FileInfo()
			destPath := filepath.Join(localPath, trimmedName)
			if fi.IsDir() {
				os.MkdirAll(destPath, 0777)

			} else if fi.Mode()&os.ModeSymlink != 0 {
				if err := writeSymlinkFromZip(f, destPath); err != nil {
					return err
				}
			} else if fi.Mode().IsRegular() {
				if err := writeRegFileFromZip(f, destPath); err != nil {
					return err
				}
			} else {
				fmt.Printf("... skipping %s in zip as not a dir, symlink, or regular file.", f.Name)
			}
		}
	}

	return nil
}

// writeSymlinkFromZip ():
// ... extracts a symlink from the zip
func writeSymlinkFromZip(f *zip.File, destPath string) (err error) {
	readCloser, err := f.Open()
	if err != nil {
		return fmt.Errorf("Failed to open symlink %s in zip: %s", f.Name, err)
	}

	byteArray, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return fmt.Errorf("Failed to read target of symlink %s in zip: %s", f.Name, err)
	}

	target := string(byteArray[:]) // ... stringify the byteArray

	if err := os.Symlink(target, destPath); err != nil {
		return fmt.Errorf("Could not create symlink %s from zip: %s", f.Name, err)
	}

	return
}

// writeRegFileFromZip ():
// ... extracts a regular file (e.g. not a socket, fifo, dir, symlink etc ...) from a zip
func writeRegFileFromZip(f *zip.File, destPath string) (err error) {
	// Read the file into a byte array
	readCloser, err := f.Open()
	defer readCloser.Close()

	if err != nil {
		return fmt.Errorf("Failed to open file %s: %s", f.Name, err)
	}

	byteArray, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return fmt.Errorf("Failed to read file %s: %s", f.Name, err)
	}

	// Write the file
	err = ioutil.WriteFile(destPath, byteArray, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write file: %s", err)
	}

	return
}

// doUnpack ():
// ... func name prefixed 'do' as field 'unpack' already exists in fetchOpts{}
func (o *fetchOpts) doUnpack(sourceFileName string) (err error) {
	fileExt, err := detectFileType(sourceFileName)
	if err != nil {
		return fmt.Errorf("Error detecting filetype of %s: %s", sourceFileName, err)
	}

	switch fileExt {
	case "gz":
		if err = o.gunzip(sourceFileName); err != nil {
			return err
		}
	case "tar":
		if err = o.untar(sourceFileName); err != nil {
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
func (o *fetchOpts) gunzip(sourceFileName string) error {
	destDir := o.destDir
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

	// need to provide a full path for location of gunzipped file
	gunZipped := filepath.Join(destDir, fmt.Sprintf("%s.gunzipped", filepath.Base(sourceFileName)))

	if o.verbose {
		fmt.Printf("Gunzipping %s\n", sourceFileName)
	}
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
	if err = o.doUnpack(newSource); err != nil {
		return err
	}

	// delete source and gunzipped
	if warning := os.Remove(sourceFileName); warning != nil {
		fmt.Printf("Could not remove intermediary file %s\n", sourceFileName)
	}
	return err
}

// untar : untars arg1 archive in to arg2 dir
func (o *fetchOpts) untar(sourceFileName string) error {
	destDir := o.destDir
	if o.verbose {
		fmt.Printf("Untarring %s\n", sourceFileName)
	}
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
func gitHubZipRequest(c commit, gitHubToken string) (*http.Request, error) {
	var request *http.Request

	// This represents either a commit, branch, or git tag
	var gitRef string
	if c.commitSha != "" {
		gitRef = c.commitSha
	} else if c.branch != "" {
		gitRef = c.branch
	} else if c.GitTag != "" {
		gitRef = c.GitTag
	} else {
		msg := "Neither a commitSha nor a GitTag nor a branch were specified " +
			"so impossible to identify a specific commit to download."
		return request, fmt.Errorf(msg)
	}

	url := fmt.Sprintf(
		"%s/repos/%s/%s/zipball/%s",
		c.Repo.Api,
		c.Repo.Owner,
		c.Repo.Name,
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
