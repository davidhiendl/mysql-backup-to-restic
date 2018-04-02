#!/bin/bash

DOCKER_REPO=dhswt/mysql-backup-to-s3
BINARY_PATH=./dist
BINARY_NAME=mysql-backup-to-s3
BINARY_TARGET="${BINARY_PATH}/${BINARY_NAME}"

function build {
   set-gopath
   echo "building binary..."
   echo "> target: ${BINARY_TARGET}"
   go build -i -ldflags="-s -w" -o "${BINARY_TARGET}" ./main.go
   upx "${BINARY_TARGET}"
}

function build-dev {
   set-gopath
   echo "building binary..."
   echo "> target: ${BINARY_TARGET}"
   go build -i -ldflags="-s -w" -o "${BINARY_TARGET}" ./main.go
}

function set-gopath {
    LOCAL_GO_PATH=$(realpath $PWD/../../../..)
    GOPATH=$GOPATH:$LOCAL_GO_PATH
}

function exec-glide {
    set-gopath
    glide "${@:1}"
}

function test-run {
    set-gopath
   	eval $(egrep -v '^#' .test-env | xargs) go run main.go
}

function image {
    echo "Building ${DOCKER_REPO} image, this might take a long time..."; \
	docker build --squash -t $DOCKER_REPO:dev .
}

function push-dev {
    echo "Pushing ${DOCKER_REPO} image, this might take a long time..."; \
	docker push $DOCKER_REPO:dev
}

case "$1" in
    build)
        build "${@:2}"
        ;;

    build-dev)
        build-dev "${@:2}"
        ;;

    package-deb)
        package-deb "${@:2}"
        ;;

    exec-glide)
        exec-glide "${@:2}"
        ;;

    image)
        image "${@:2}"
        ;;

    push-dev)
        push-dev "${@:2}"
        ;;

    test-run)
        test-run "${@:2}"
        ;;

    *)
        echo $"Usage: $0 {build|build-dev|package-dev|exec-glide|image|test-run}"
        exit 1
esac
