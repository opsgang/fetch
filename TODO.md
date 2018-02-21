# TODO

# destination handling

fubar for multiple assets that include a single file as well.

# --decompress

Handle untgz based on file extension (.tgz or .tar.gz)

Only for release assets.

Should remove archive after extraction.

# --which-tag

FetchTags returns the latest tag or else the one that meets a constraint.

Offer `--which-tag`, which will only display a tag that would have been
downloaded.

* NTH Account for when release asset does not exist - should fail

* Account for when no tag meeting a constraint exists - should fail, not default
  to latest.

# --tag-prefix (filter)

Account for only fetching tags that have a prefix with an optional delimiter
before x.y.z (can be a \., \- or \_)

