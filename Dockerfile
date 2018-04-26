# build binary first
FROM    golang:1.9.4-alpine3.7

# install tools
RUN     apk add --no-cache \
            glide \
            git \
            upx

# add glide config and install dependencies with glide in a separate step to speed up subsequent builds
WORKDIR /go/src/github.com/davidhiendl/mysql-backup-to-restic
ADD     glide.lock glide.yaml /go/src/github.com/davidhiendl/mysql-backup-to-restic/
RUN     glide install

# add source and build package
ADD     . /go/src/github.com/davidhiendl/mysql-backup-to-restic/
RUN     glide install \
&&      go build -i \
            -o /mysql-backup-to-restic \
            -ldflags="-s -w" \
            ./main.go \

# compress binary
&&      upx /mysql-backup-to-restic

# build container next
FROM    alpine
LABEL   maintainer="David Hiendl <david.hiendl@dhswt.de>"

ENV     LOG_LEVEL=info \
        CONFIG_FILE=/var/run/secrets/config.yaml

RUN     apk add --no-cache \
            mariadb-client \
            ca-certificates \
&&      rm /usr/bin/myisam_ftdump \
&&      rm /usr/bin/mysql_waitpid \
&&      rm /usr/bin/mysqladmin \
&&      rm /usr/bin/mysqlcheck \
&&      rm /usr/bin/mysqlimport \
&&      rm /usr/bin/mysqlshow

RUN     apk add --no-cache --virtual .build-deps \
            bzip2 \
            curl \
&&      curl -L https://github.com/restic/restic/releases/download/v0.8.3/restic_0.8.3_linux_amd64.bz2 \
            -o restic.bz2 \
&&      bzip2 -d restic.bz2 \
&&      chmod 755 restic \
&&      mv restic /usr/bin/restic \
&&      apk del .build-deps

# add binary from previous stage
COPY    --from=0 /mysql-backup-to-restic /
ENTRYPOINT ["/mysql-backup-to-restic"]
CMD        [""]
