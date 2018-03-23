package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

var s *httptest.Server

func TestMain(m *testing.M) {
	s = apiStub()
	defer s.Close()
	code := m.Run()
	os.Exit(code)
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

func apiStub() *httptest.Server {
	var resp string
	var err bool
	var counter = 0 // used for keeping track of retries

	return httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				if v, ok := r.Header["Authorization"]; ok {
					w.Header().Set("X-Authorization", v[0])
				}
				switch {
				case strings.HasPrefix(r.RequestURI, "/repos/foo/bar"):
					resp, err = stubFooBar(w, r)
					if resp == "" || err == false {
						return
					}

				case strings.HasPrefix(r.RequestURI, "/repos/has/no"):
					resp, err = stubHasNo(w, r)
					if resp == "" || err == false {
						return
					}

				case strings.HasPrefix(r.RequestURI, "/repos/sna/fu"):
					resp, err = stubSnaFu(w, r)
					if resp == "" || err == false {
						return
					}

				// keep this one in the closure, as it needs counter
				// and we don't need to be passing pointers to ints ...
				case strings.HasPrefix(r.RequestURI, "/fail/twice"):
					if counter == 2 {
						resp = fmt.Sprintf("counter: %d", counter)
					} else {
						counter++
						http.Error(w, "Remote failure", http.StatusBadGateway)
						return
					}

				case strings.HasPrefix(r.RequestURI, "/fail/always"):
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

func stubSnaFu(w http.ResponseWriter, r *http.Request) (resp string, ok bool) {
	switch r.RequestURI {

	case "/repos/sna/fu/releases?per_page=100&page=1":
		w.Header().Set("Link", noValidRelsPage1Link)
		return noValidRelsPage1, true

	case "/repos/sna/fu/releases?per_page=100&page=2":
		w.Header().Set("Link", noValidRelsPage2Link)
		return noValidRelsPage2, true

	default:
		http.Error(w, "not found", http.StatusNotFound)
	}

	return "", false
}

func stubHasNo(w http.ResponseWriter, r *http.Request) (resp string, ok bool) {
	switch r.RequestURI {
	case "/repos/has/no/tags?per_page=100&page=1":
		return apiNoTags, true

	case "/repos/has/no/releases?per_page=100&page=1":
		return apiNoReleases, true

	default:
		http.Error(w, "not found", http.StatusNotFound)
	}

	return "", false
}

func stubFooBar(w http.ResponseWriter, r *http.Request) (resp string, ok bool) {

	switch r.RequestURI {
	case "/repos/foo/bar/tags?per_page=100&page=1":
		w.Header().Set("Link", apiTagsPage1Link)
		return apiTagsPage1, true

	case "/repos/foo/bar/tags?per_page=100&page=2":
		w.Header().Set("Link", apiTagsPage2Link)
		return apiTagsPage2, true

	case "/repos/foo/bar/zipball/46.0.0":
		if err := serveFile(w, "test-fixtures/46.0.0.zip"); err != nil {
			http.Error(w, "Could not serve file", http.StatusInternalServerError)
			return "", false
		}
		return "", true

	case "/repos/foo/bar/releases?per_page=100&page=1":
		w.Header().Set("Link", relsPage1Link)
		return relsPage1, true

	case "/repos/foo/bar/releases?per_page=100&page=2":
		w.Header().Set("Link", relsPage2Link)
		return relsPage2, true

	case "/repos/foo/bar/releases/tags/7.6.5?per_page=100&page=":
		return relWithAsset, true

	case "/repos/foo/bar/releases/assets/7654783?per_page=100&page=":
		if err := serveFile(w, "test-fixtures/packed.tgz"); err != nil {
			http.Error(w, "Could not serve file", http.StatusInternalServerError)
			return "", false
		}
		return "", true

	default:
		http.Error(w, "not found", http.StatusNotFound)
	}

	return "", false
}
