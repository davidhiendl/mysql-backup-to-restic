#!/bin/bash

DOCKER_REPO=dhswt/mysql-backup-to-restic
BINARY_PATH=./dist
BINARY_NAME=mysql-backup-to-restic
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

function test-run-docker {
    set-gopath
    image dev
    docker run \
        -ti \
        --net=host \
        -v $(pwd)/conf.test.d/config.yaml:/var/run/secrets/config.yaml \
        ${DOCKER_REPO}:dev
}

function image {
    TAG=${1:-dev}
    echo "Building ${DOCKER_REPO}:${TAG} image, this might take a long time..."; \
	docker build --squash -t $DOCKER_REPO:${TAG} .
}

function push-dev {
    image
    echo "Pushing ${DOCKER_REPO}:dev image, this might take a long time..."; \
	docker push $DOCKER_REPO:dev
}

function push-master {
    image master
    echo "Pushing ${DOCKER_REPO}:master image, this might take a long time..."; \
	docker push $DOCKER_REPO:master
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

    glide)
        exec-glide "${@:2}"
        ;;

    image)
        image "${@:2}"
        ;;

    push-dev)
        push-dev "${@:2}"
        ;;

    push-master)
        push-master "${@:2}"
        ;;

    test-run)
        test-run "${@:2}"
        ;;

    test-run-docker)
        test-run-docker "${@:2}"
        ;;

    *)
        echo $"Usage: $0 {build|build-dev|package-dev|glide|image|test-run}"
        exit 1
esac
