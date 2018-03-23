package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestTagsList(t *testing.T) {
	t.Parallel()

	// fetch tags or release tags based on cmd-line flags
	r := repo{
		Url:   "https://github.com/foo/bar",
		Owner: "foo",
		Name:  "bar",
		Token: "dummytoken",
		Api:   s.URL,
	}

	// happy path ... release assets
	oRel := fetchOpts{}
	oRel.relAssets = []string{"foo.tgz", "bar.tgz"}
	oRel.verbose = true

	if tags, err := oRel.tagsList(r); err != nil {
		t.Fatalf("... expected a list of published release tags from foo/bar, got\n%s", err)
	} else if !reflect.DeepEqual(tags, relTagsExpected) {
		t.Fatalf("got: %#v\nexpected: %#v", tags, relTagsExpected)
	}

	// happy path ... tags
	oTag := fetchOpts{}
	oRel.verbose = true

	if tags, err := oTag.tagsList(r); err != nil {
		t.Fatalf("... expected a tag list from foo/bar, got\n%s", err)
	} else if !reflect.DeepEqual(tags, apiTagsExpected) {
		t.Fatalf("got: %#v\nexpected: %#v", tags, apiTagsExpected)
	}

}

func TestDo(t *testing.T) {
	t.Parallel()

	r := repo{
		Url:   "https://github.com/foo/bar",
		Owner: "foo",
		Name:  "bar",
		Token: "dummytoken",
		Api:   s.URL,
	}

	testDoFromPaths(t, r)
	testDoReleaseAssets(t, r)
}

func testDoFromPaths(t *testing.T, r repo) {

	// default behaviour:
	// will download whole repo from latest tag.
	//
	// * from-path of / if user specified neither from-path nor release-asset
	// * latest tag used.
	//
	// NOTE cli prevents defaulting to latest tag as user is required to
	// choose a tag or tag constraint or branch or commit.

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	o := fetchOpts{}
	o.destDir = tempDir

	err = o.do(r)
	if err != nil {
		t.Fatalf("Test failed to download latest foo/bar tag's src.%s", err)
	}

	// ... check for file in decompressed /root folder
	f := fmt.Sprintf("%s/symlink", tempDir)
	if yes := fileExists(f); !yes {
		t.Fatalf("expected file %s does not exist", f)
	}
}

func testDoReleaseAssets(t *testing.T, r repo) {

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(tempDir)

	o := fetchOpts{}
	o.destDir = tempDir
	o.relAssets = []string{"packed.tgz"}

	err = o.do(r)
	if err != nil {
		t.Fatalf("Test failed to download release (7.6.5).\n%s", err)
	}
}
