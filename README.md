# Moved to GitLab: https://gitlab.com/davidhiendl/mysql-backup-to-restic

# MySQL Backup to Restic

## WARNING: ALPHA QUALITY; MINIMAL TESTING; USE AT YOUR OWN RISK

## Description
Create MySQL backups with `mysqldump` and backup them with `Restic` to any supported backend (note: initial implementation only supports S3).

## Docker Image
**Using pre-built image:**
```bash
docker run -ti \
    -v ./config-dir/config.yaml:/var/run/secrets/config.yaml \
    dhswt/mysql-backup-to-restic:alpha
```

**Building the image yourself:** \
The entire build (including building the binary) is included in the [Dockerfile](Dockerfile).
```bash
docker build -t your-username/mysql-backup-to-restic:<tag>
```

## Binary env variables
| Variable    | Default     | Description                                             |
| ---         | ---         | ---                                                     |
| LOG_LEVEL   | info        | Log level, any of: debug,info,warning,error,fatal,panic |
| CONFIG_FILE | config.yaml | the absolute or relative path to the configuration file |

## Docker Env Variables
| Variable    | Default                      | Description                                             |
| ---         | ---                          | ---                                                     |
| LOG_LEVEL   | info                         | Log level, any of: debug,info,warning,error,fatal,panic |
| CONFIG_FILE | /var/run/secrets/config.yaml | the absolute or relative path to the configuration file |

## Dependencies
- GO >= 1.9

