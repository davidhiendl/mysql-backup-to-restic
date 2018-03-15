.DEFAULT_GOAL := list

SHELL:=/bin/bash

list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

image-repo  = dhswt/mysql-backup-to-s3
source-path = /go/src/github.com/davidhiendl/mysql-backup-to-s3
binary-path = ./dist
binary-name = mysql-backup-to-s3

run:
	source .test-env-export ; \
	export GOPATH=$$GOPATH:$$(dirname $$(dirname $$(dirname $$(dirname $$PWD)))); \
	echo "GOPATH=$${GOPATH}" ; \
	go run main.go

run-container:
	docker run -ti --net=host --env-file ./.test-env $(image-repo):master

# build compressed binary using local go
binary:
	GOPATH=$$GOPATH:$$PWD/../../../../ \
	&& echo $$GOPATH \
	&& go get . \
	&& go build -i -ldflags="-s -w" -o $(binary-path)/$(binary-name) ./main.go \
	&& upx $(binary-path)/$(binary-name)

# build using local go
binary-dev:
	GOPATH=$$GOPATH:$$PWD/../../../../ \
	&& echo $$GOPATH \
	&& go get . \
	&& go build -i -o $(binary-path)/$(binary-name) ./main.go

# build image
image:
	echo "building image, this might take a long time..." && \
	docker build -t $(image-repo):master .

tag-push-testing:
	docker tag $(image-repo):master $(image-repo):testing && \
	docker push $(image-repo):testing

show-images:
	docker images | grep "$(image-repo)"

# Remove dangling images
clean-images:
	docker images -a -q \
		--filter "reference=$(image-repo)" \
		--filter "dangling=true" \
	| xargs docker rmi

# Remove all images
clear-images:
	docker images -a -q \
		--filter "reference=$(image-repo)" \
	| xargs docker rmi
