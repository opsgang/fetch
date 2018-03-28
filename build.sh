# vim: et sr sw=4 ts=4 smartindent syntax=sh:
setup_fetch() {
  export GOPATH=${GOPATH:-/go} GOBIN=/usr/local/go/bin
  export LGOBIN=$GOBIN PATH=$PATH:$GOBIN
  export __WD=$GOPATH/src/github.com/opsgang/fetch
  export PATH=$GOBIN:$PATH
  export __V=${FETCH_VERSION:-v1.0.0}

  get_glide || return 1
}

wd() {
    [[ -z "$__WD" ]] && echo "Run setup_*() first" && return 1
    [[ "$PWD" != "$__WD" ]] && cd $__WD
    return 0
}

build_linux() {
    (
        _build
    )
}

build_macos() {
    (
        export GOOS=darwin GOARCH=amd64
        _build
    )
}

get_glide() {
    (
        command -v glide >/dev/null && return 0
        wd || return 1
        if on_alpine && not_root
        then
            curl https://glide.sh/get | su-exec root sh || exit 1
            su-exec root glide install
        else
            curl https://glide.sh/get | sh || exit 1
            glide install
        fi
    )
}

# ... account for building on alpine
_build() {
    local fl="-w -extldflags \"-static\" -X main.VERSION=$__V -X main.TIMESTAMP=build-$(date '+%Y%m%d%H%M%S')"
    (
        wd || return 1
        export CGO_ENABLED=0
        if on_alpine && not_root
        then
            su-exec root go build --ldflags "$fl" -o $GOBIN/ghfetch .
        else
            go build --ldflags "$fl" -o $GOBIN/ghfetch .
        fi
    )
}

# if on alpine we need to use su-exec root to build successfully.
on_alpine() {
    (
        . /etc/os-release 2>/dev/null || exit 1
        [[ "$NAME" =~ Alpine ]] && exit 0
        exit 1
    )
}

not_root() {
    [[ "$(id -u)" -ne 0 ]]
}
