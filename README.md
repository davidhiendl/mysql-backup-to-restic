# Telegraf Docker Service Discovery

## Description
Create MySQL backups with `mysqldump`

## Docker Image
**Using pre-built image:**
```bash
docker run -ti \
    -e ...
    dhswt/mysql-backup-to-s3:stable
```

**Building the image yourself:** \
The entire build (including building the binary) is included in the [Dockerfile](Dockerfile).
```bash
docker build -t yourprefix/mysql-backup-to-s3:<tag>
```

## Configuration Variables
| Variable             | Default                  | Description                                                                                 |
| ---                  | ---                      | ---                                                                                         |
| TSD_TEMPLATE_DIR     | /etc/telegraf/sd-tpl.d   | Where configurations templates are taken from                                               |
| TSD_CONFIG_DIR       | /etc/telegraf/telegraf.d | Where configurations are written to, the telegraf config directory                          |
| TSD_TAG_SWARM_LABELS | true                     | If docker swarm labels should be imported as tags. See `Container Detection > Swarm Labels` |
| TSD_TAG_LABELS       | none                     | A list of comma separated labels that should be added as tags                               |
| TSD_QUERY_INTERVAL   | 15                       | Interval in seconds between querying of the docker api for changes                          |



## Dependencies
- GO >= 1.9

