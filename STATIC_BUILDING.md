# Building Statically

        docker pull golang:alpine # latest go env for building golang projects

        docker run -it --name build_static_fetch --rm \
          -e GIT_REF=v0.1.1 \
          -e GIT_URL=https://github.com/opsgang/fetch \
            golang:alpine /bin/sh
        
        # ... inside container
        export GOPATH=/go GOBIN=/usr/local/go/bin
        export LGOBIN=$GOBIN PATH=$PATH:$GOBIN
        apk --no-cache --update add ca-certificates bash gcc musl-dev openssl git curl

        echo "... getting glide (dependency manager)"
        curl https://glide.sh/get | sh || exit 1

        cd /go/src || exit 1
        git clone $GIT_URL --depth 1 --branch $GIT_REF || exit 1

        cd fetch || exit 1
        echo "... installing project deps (using glide)"
        glide install || exit 1

        echo "... building static binary $GOBIN/fetch"
        # -w is to drop debugging related symbols etc for smaller binary
        ts=$(date '+%Y%m%d%H%M%S')
        ldf="-w -extldflags \"-static\" -X main.VERSION=$GIT_REF -X main.TIMESTAMP=build-$ts"
        go build --ldflags "$ldf" -o $GOBIN/fetch . || exit 1

        echo "... basic verification that binary works"
        fetch --help || exit 1
        fetch --version | grep -Po "$GIT_REF$" >/dev/null || exit 1

