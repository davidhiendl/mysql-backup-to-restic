package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	S3AccessKey      string `envconfig:"S3_ACCESS_KEY"`
	S3SecretKey      string `envconfig:"S3_SECRET_KEY"`
	S3Endpoint       string `envconfig:"S3_ENDPOINT" default:"s3.us-west-1.amazonaws.com"`
	S3Region         string `envconfig:"S3_REGION" default:"us-west-1"`
	S3Bucket         string `envconfig:"S3_BUCKET"`
	S3PathPrefix     string `envconfig:"S3_PATH_PREFIX" default:""`
	S3ForcePathStyle bool   `envconfig:"S3_FORCE_PATH_STYLE" default:"true"`

	MySQLDumpBinary string `envconfig:"MYSQLDUMP_BINARY" default:"mysqldump"`

	MySQLUser string `envconfig:"MYSQL_USER" default:"root"`
	MySQLPass string `envconfig:"MYSQL_PASS"`
	MySQLHost string `envconfig:"MYSQL_HOST" default:"127.0.0.1"`
	MySQLPort int    `envconfig:"MYSQL_PORT" default:"3306"`

	ExportTempDir       string `envconfig:"EXPORT_TEMP_DIR" default:"/tmp"`
	SkipSystemDatabases bool   `envconfig:"SKIP_SYSTEM_DATABASES" default:"true"`
	CompressWithGz      bool   `envconfig:"COMPRESS_WITH_GZ" default:"true"`

	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

func Load() *Config {
	c := Config{}

	err := envconfig.Process("", &c)
	if err != nil {
		logrus.Panicf("failed to parse log level, %+v", err)
	}

	// set log level
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		logrus.Panicf("failed to parse log level, %+v", err)
	}
	logrus.SetLevel(level)

	return &c
}
