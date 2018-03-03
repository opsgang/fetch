package main

// the const below contain zero-width spaces to force urfave/cli-generated
// help text to separate sections with a blank line.

const usage_lead = `
	GITHUB_TOKEN=$access_token ghfetch [global options] /my/downloads/dir
	See https://github.com/opsgang/fetch for examples, argument definitions, and additional docs.
`

const txt_usage = `
	Download selected subfolders/files OR release attachments of a GitHub repo.
	Select a specific git commit, branch, or tag.
	Specify a semantic version constraint for tags e.g. '>=1.0,<2.0' !
	Choose to automatically unpack release attachment tars and gzips!!
	Verify downloaded release attachments against gpg asc signature files!!!
	​
`

const txt_repo = `
Required. Fully qualified url of the github repo.
	​
`

const txt_commit = `
Git commit SHA1 to checkout. Overrides --branch and --tag.
	​
`

const txt_branch = `
Git branch - will checkout HEAD commit. Overrides --tag.
	​
`

const txt_tag = `
Git tag to download, expressed with Hashicorp's version constraint operators.
	If empty, ghfetch will download the latest tag.
	Examples: https://github.com/opsgang/fetch#version-constraint-operators.
	​
`

const txt_from_path = `
Get contents of `+"`/PATH/IN/REPO`"+`. The folder itself is not created locally,
	only its contents. /PATH/IN/REPO can also be the path to a single file.
	If this or --release-asset aren't specified, all files are downloaded.
	Specify multiple times to download multiple folders or files.
	​
	e.g. # puts libs/* and build.sh in to /opt/libs
		--from-path='/libs/' --from-path='/scripts/build.sh' /opt/libs
	​
`

const txt_release_asset = `
Name of github release attachment to download. Requires --tag.
	Specify multiple times to grab more than one attachment.
	​
	e.g. # get foo.tgz and bar.txt from latest 1.x release attachments
		--tag='~>1.0' --release-asset='foo.tgz' --release-asset='bar.txt',
	​
`

const txt_unpack = `
Whether to unpack a compressed release attachment. Requires --release-asset.
	Only unpacks tars, tar-gzip and gzip, otherwise does nothing.
	​
	e.g. # unpack latest 1.x tag of foo.tgz in to /var/tmp/foo
		--tag='~>1.0' --unpack --release-asset='foo.tgz',
	​
`

const txt_gpg_public_key = "`/PATH/TO/KEY` " + `to verify downloaded release assets.
	Requires --release-asset.
	If set, will look for <asset-name>.asc or <asset-name>.asc.txt attached
	to the chosen release. That signature and /PATH/TO/KEY will be used
	for gpg verification.
	If signature file not found, or verification fails, the release-asset is deleted.
	​
`

const txt_token = `
GitHub Personal Access Token, required to download from a private repo.
	Also enables GitHub api's higher rate-limit. It is recommended to set
	this in the env var GITHUB_TOKEN before invoking ghfetch,
	instead of via this commandline option.
	​
`

