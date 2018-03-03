# TODO

## --tag-prefix (filter)

Account for only fetching tags that have a prefix with an optional delimiter
before x.y.z (can be a \., \- or \_)

## --which-tag

FetchTags returns the latest tag or else the one that meets a constraint.

Offer `--which-tag`, which will only display a tag that would have been
downloaded.

* NTH Account for when release asset does not exist - should fail

* Account for when no tag meeting a constraint exists - should fail, not default
  to latest.

# IN PROGRESS

## rename to ghfetch

_fetch_ is far too overreaching a name - it only works with the github api, much to
the exasperation of my bitbucket-using-colleagues.

Renaming this to ghfetch is not only clearer, it also frees up the fetch namespace
for new variants such as bbfetch. At which point maybe my colleagues will speak
to me again.

Still need to rename all vars, methods that use Fetch to use something else.

# DONE

## destination handling

FUBAR when a source-path is a single file.

The destination parameter should **always** be treated as a directory.

This dir should always be created.

## --unpack

Handle untgz based on file extension (.tgz or .tar.gz)

Only for release assets.

Should remove archive after extraction.

Should handle symlinks in tars.

## --signer-public-key=/path/to/file

See https://gist.github.com/lsowen/d420a64821414cd2adfb for using openpgp

Only if --release-asset specified.

Will not keep any downloaded {{release-asset}} unless a {{release-asset}}.asc
or {{release-asset}}.asc.txt is attached to the release as well.

The file and asset are downloaded, then the pgp sig is checked and
unverified download are deleted.

On success the downloaded .asc is deleted.

