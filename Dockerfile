# build binary first
FROM    golang:1.9.4-alpine3.7

WORKDIR /go/src/github.com/davidhiendl/mysql-backup-to-s3

# install upx to compress binary
RUN     apk add --no-cache \
            git \
            mercurial \
            upx \
            curl \

# install go dep
&&      curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# add sources
ADD     . /go/src/github.com/davidhiendl/mysql-backup-to-s3/

# fetch remaining dependencies and build package
RUN     dep ensure \
&&      go build -i \
            -o /mysql-backup-to-s3 \
            -ldflags="-s -w" \
            ./main.go \

# compress binary
&&      upx /mysql-backup-to-s3

# build container next
FROM    alpine
LABEL   maintainer="David Hiendl <david.hiendl@dhswt.de>"

RUN     apk add --no-cache \
            mariadb-client \
            ca-certificates

# add binary from previous stage
COPY    --from=0 /mysql-backup-to-s3 /
ENTRYPOINT ["/mysql-backup-to-s3"]
CMD        [""]
