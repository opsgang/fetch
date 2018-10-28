# TODO

## --tag-regex-filter (filter)
Account for only fetching tags that have a prefix leading the semver string
and / or a suffix after the semver string.

To avoid ambiguous results, the user is urged to use fixed strings whereever
possible.

EXAMPLE

repo contains src for multiple built artefacts. Tags exist for each type of
artefact:

foo-app-1.0.8-20180711120000
foo-app-1.0.10-20180711120000 # (ver 1.0.10)
bar-app-1.0.3-20180711120000
bar-app-1.0.10-20180711120000 # (ver. 1.0.10 as above, but different app)
bar-app-2.1.10-20180711120000

If I want the latest bar-app, I should use a filter pattern like:

`^bar\-app\-(\d+\.\d+\.\d+)` # if we know the semver ALWAYS follows a set prefix

OR

`^bar\-app\-(.*)\-\d+$` # match of semver based on known suffix and prefix format.

OR

`^bar\-app\-(\d+\.\d+\.\d+)\-\d+$` # strict format of desired tag string.

The following patterns are bad:

`(\d+\.\d+\.\d+)`        # will match all foo-app tags as well.
`.*\-(\d\.\d+\.\d+)\-.*` # will match all foo-bar tags as well.
`bar-app\-(.*)`          # will assume the timestamp suffix is part of the semver string.

* pattern MUST contain EXACTLY ONE capture group to match the semver.

* all other parts of the pattern are to identify which tags are relevant.

* Account for when no tag meeting a constraint exists - should fail, not default
  to latest.

* Account for when many tags meet a pattern - should fail and tell user to
  tighten up the pattern.

  e.g. if tags are foo01-1.0.11-20180711120000 and bar01-1.0.11-some-string,
  the pattern `.*-SEMVER-.*`, both tags get matched with SEMVER of 1.0.11.

## end-to-end tests

* remove tests against real repos from \_test.go files

* add a script to run after a successful build that does a few
    `fetch` runs e.g. against a private repo, using a constraint,
    for a release, using from-path, using commit, using branch,
    using gpg. Compare to expected output.

## refactoring

* rewrite it separating methods for fetchOptions, rel and tag structs.

* correct docco.

* stop fetching tags if only need a commit or branch!

* abstract API calls so can work for different git providers e.g. gitlab

## rename to ghget | glget | bbget (it's shorter)

_fetch_ is far too overreaching a name - it only works with the github api, much to
the exasperation of my bitbucket-using-colleagues.

Renaming this to ghfetch is not only clearer, it also frees up the fetch namespace
for new variants such as bbfetch. At which point maybe my colleagues will speak
to me again.

Still need to rename all vars, methods that use Fetch to use something else.

# IN PROGRESS

# DONE

## --which-tag
FetchTags returns the latest tag or else the one that meets a constraint.

Offer `--which-tag`, which will only display a tag that would have been
downloaded.


* Account for when release asset does not exist - should fail

## --timeout i

Let user specify net/http timeout (for those larger repos and assets)

## --verbose
Suppress output unless this is specified.

## mocking of api calls

## better tests
* - [x] for large tag sets
* - [x] for repo with no tags
* - [x] for repo with no releases
* - [x] for releases with out all assets requested

## http transport settings

* keep-alive
* connection timeouts
* idle timeouts

## http retry on 5xx
* on all http calls with back-off

## pagination
* on all api calls

## error handling

I want to reduce cyclomatic complexitiy of functions. This can be achieved
by returning only the one type of error from a function. e.g. not 'error' and FetchError,
but only FetchError.

Separate out http response codes - see newError e.g. in github.go:callGitHubApi()
Stop with error codes for anything else.

That means FetchError only needs a message constructed.

**In fact we probably don't need FetchError at all, but just error, and some functions
to construct the error messages.**

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

## NOTES ON SEMVER

valid semver pattern

v?\d+\.\d+\.\d+(\-[\w]+[-\w\.]*)?(+[\w]+[-\w\.]*)?
