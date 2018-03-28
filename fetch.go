package main

import (
	"fmt"
	"time"
)

func (o *fetchOpts) do(r repo) (err error) {

	// Get the tags for the given repo
	// or tags for actual releases if getting release assets.
	tags, err := o.tagsList(r)
	if err != nil {
		return err
	}

	tag, err := o.tagToGet(tags)
	if err != nil {
		return err
	}

	if o.whichTag {
		return
	}

	// If no release assets or from-paths specified, assume
	// user wants all files from zipball.
	if len(o.fromPaths) == 0 && len(o.relAssets) == 0 {
		o.fromPaths = []string{"/"}
	}

	// Download any requested source files
	if err := o.downloadFromPaths(r, tag); err != nil {
		return err
	}

	// Download any requested release assets
	if err := o.downloadReleaseAssets(r, tag); err != nil {
		return err
	}

	return
}

// tagsList ():
// returns str slice of tags - release tags, if --release-asset is specified.
func (o *fetchOpts) tagsList(r repo) (tags []string, err error) {
	if len(o.relAssets) == 0 {
		tags, err = fetchTags(r, o.verbose)
	} else {
		tags, err = o.fetchReleaseTags(r)
	}
	if err != nil {
		return tags, fmt.Errorf("Error occurred while getting tags from GitHub repo: %s", err)
	}

	return
}

// setTimeout ():
// sets net/http transport's Timeout attr
func (o *fetchOpts) setTimeout() {
	if o.timeout != 0 {
		if o.verbose {
			fmt.Printf("... setting download timeout to %d seconds", o.timeout)
		}
		cl.Timeout = time.Second * time.Duration(o.timeout)
	}
}
