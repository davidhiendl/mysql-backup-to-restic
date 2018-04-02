# build binary first
FROM    golang:1.9.4-alpine3.7

# install tools
RUN     apk add --no-cache \
            glide \
            git \
            upx

# add glide config and install dependencies with glide in a separate step to speed up subsequent builds
WORKDIR /go/src/github.com/davidhiendl/mysql-backup-to-s3
ADD     glide.lock glide.yaml /go/src/github.com/davidhiendl/mysql-backup-to-s3/
RUN     glide install

# add source and build package
ADD     . /go/src/github.com/davidhiendl/mysql-backup-to-s3/
RUN     glide install \
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
            ca-certificates \
&&      rm /usr/bin/myisam_ftdump \
&&      rm /usr/bin/mysql_waitpid \
&&      rm /usr/bin/mysqladmin \
&&      rm /usr/bin/mysqlcheck \
&&      rm /usr/bin/mysqlimport \
&&      rm /usr/bin/mysqlshow

# add binary from previous stage
COPY    --from=0 /mysql-backup-to-s3 /
ENTRYPOINT ["/mysql-backup-to-s3"]
CMD        [""]
