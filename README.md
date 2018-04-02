# MySQL Backup to S3

## WARNING: ALPHA QUALITY; MINIMAL TESTING; USE AT YOUR OWN RISK

## Description
Create MySQL backups with `mysqldump` and store them to a S3-compatible
storage. This allows for long-term, off-site and backup storage while
eliminating the need for any self-hosted backup storage.

Why `mysqldump`?: Because it is reliable and performance is acceptable. No point in re-inventing a backup tool.

## Docker Image
**Using pre-built image:**
```bash
docker run -ti \
    -e ...
    dhswt/mysql-backup-to-s3:alpha
```

**Building the image yourself:** \
The entire build (including building the binary) is included in the [Dockerfile](Dockerfile).
```bash
docker build -t yourprefix/mysql-backup-to-s3:<tag>
```

## Configuration Variables
| Variable              | Default                    | Description                                                                                                       |
| ---                   | ---                        | ---                                                                                                               |
| S3_ACCESS_KEY         |                            | S3 Access Key                                                                                                     |
| S3_SECRET_KEY         |                            | S3 Secret Key                                                                                                     |
| S3_ENDPOINT           | s3.us-west-1.amazonaws.com | The endpoint to use, also works with other s3 compatibile services like DigitalOcean spaces                       |
| S3_REGION             | us-west-1                  | Region to use for S3                                                                                              |
| S3_BUCKET             |                            | Bucket to use for backups                                                                                         |
| S3_PATH_PREFIX        | mysql-backup               | A path prefix inside the bucket                                                                                   |
| S3_FORCE_PATH_STYLE   | true                       | Force S3 URLs to use bucket as folder inside of subdomain (more compatible with minio and similar solutions       |
| MYSQLDUMP_BINARY      | mysqldump                  | The location of the binary, if it is in path it should just work. Otherwise a full path may be required           |
| MYSQL_USER            | root                       | MySQL User                                                                                                        |
| MYSQL_PASS            |                            | MySQL Password                                                                                                    |
| MYSQL_HOST            | localhost                  | MySQL Host (Warning: localhost will use unix:// socket, use 127.0.0.1 instead if localhost via tcp/ip is required |
| MYSQL_PORT            | 3306                       | MySQL Port                                                                                                        |
| EXPORT_TEMP_DIR       | /tmp                       | A temporary directory where exports are stored before they are uploaded to S3                                     |
| SKIP_SYSTEM_DATABASES | true                       | If system databases (information_schema, performance_schema, mysql) should be skipped during backup               |
| COMPRESS_WITH_GZ      | true                       | If the backup files should be compress with gz (will also add .gz extenion)                                       |

## Dependencies
- GO >= 1.9

