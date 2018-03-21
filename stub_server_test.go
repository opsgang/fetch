package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var s *httptest.Server

func TestMain(m *testing.M) {
	s = apiStub()
	defer s.Close()
	code := m.Run()
	os.Exit(code)
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

				case "/repos/has/no/tags?per_page=100&page=1":
					resp = apiNoTags

				case "/repos/has/no/releases?per_page=100&page=1":
					resp = apiNoReleases

				case "/repos/foo/bar/tags?per_page=100&page=1":
					w.Header().Set("Link", apiTagsPage1Link)
					resp = apiTagsPage1

				case "/repos/foo/bar/tags?per_page=100&page=2":
					w.Header().Set("Link", apiTagsPage2Link)
					resp = apiTagsPage2

				case "/repos/foo/bar/zipball/46.0.0":
					if err := serveFile(w, "test-fixtures/46.0.0.zip"); err != nil {
						http.Error(w, "Could not serve file", http.StatusInternalServerError)
						return
					}

				case "/repos/foo/bar/releases?per_page=100&page=1":
					w.Header().Set("Link", relsPage1Link)
					resp = relsPage1

				case "/repos/foo/bar/releases?per_page=100&page=2":
					w.Header().Set("Link", relsPage2Link)
					resp = relsPage2

				case "/repos/foo/bar/releases/tags/7.6.5?per_page=100&page=":
					resp = relWithAsset

				case "/repos/foo/bar/releases/assets/7654783?per_page=100&page=":
					if err := serveFile(w, "test-fixtures/packed.tgz"); err != nil {
						http.Error(w, "Could not serve file", http.StatusInternalServerError)
						return
					}

				case "/repos/sna/fu/releases?per_page=100&page=1":
					w.Header().Set("Link", noValidRelsPage1Link)
					resp = noValidRelsPage1

				case "/repos/sna/fu/releases?per_page=100&page=2":
					w.Header().Set("Link", noValidRelsPage2Link)
					resp = noValidRelsPage2

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

// streams a local file
func serveFile(w http.ResponseWriter, fn string) (err error) {
	f, err := os.Open(fn)
	defer f.Close() //Close after function return

	if err != nil {
		//File not found, send 404
		http.Error(w, "Test File not found.", 404)
		return
	}

	fileHeader := make([]byte, 512)
	//Copy the headers into the fileHeader buffer
	f.Read(fileHeader)
	//Get content type of file
	fileContentType := http.DetectContentType(fileHeader)

	//Get the file size
	fileStat, _ := f.Stat()                            //Get info from file
	fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fn)
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)

	//Send the file
	//We read 512 bytes from the file already so we reset the offset back to 0
	f.Seek(0, 0)
	io.Copy(w, f) //'Copy' the file to the client
	return
}
