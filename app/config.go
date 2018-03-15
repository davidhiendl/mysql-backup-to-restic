package app

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/davidhiendl/mysql-backup-to-s3/app/logger"
)

type Config struct {
	S3AccessKey  string `envconfig:"S3_ACCESS_KEY"`
	S3SecretKey  string `envconfig:"S3_SECRET_KEY"`
	S3Endpoint   string `envconfig:"S3_ENDPOINT"`
	S3Region     string `envconfig:"S3_REGION"`
	S3Bucket     string `envconfig:"S3_BUCKET"`
	S3PathPrefix string `envconfig:"S3_PATH_PREFIX"`

	MySQLDumpBinary string `envconfig:"MYSQLDUMP_BINARY"`

	MySQLUser string `envconfig:"MYSQL_USER"`
	MySQLPass string `envconfig:"MYSQL_PASS"`
	MySQLHost string `envconfig:"MYSQL_HOST"`
	MySQLPort int    `envconfig:"MYSQL_PORT"`

	ExportTempDir       string `envconfig:"EXPORT_TEMP_DIR"`
	SkipSystemDatabases bool   `envconfig:"SKIP_SYSTEM_DATABASES"`
	CompressWithGz      bool   `envconfig:"COMPRESS_WITH_GZ"`

	LogLevel int `envconfig:"LOG_LEVEL"`
}

var DefaultConfig = Config{
	LogLevel: logger.LOG_INFO,

	S3Region:   "us-west-1",
	S3Endpoint: "s3.us-west-1.amazonaws.com",

	MySQLDumpBinary: "mysqldump",

	MySQLHost: "127.0.0.1",
	MySQLPort: 3306,

	ExportTempDir:       "/tmp",
	SkipSystemDatabases: true,
	CompressWithGz:      true,
}

// Create new config and populate it from environment
func NewConfig() (*Config) {
	c := DefaultConfig
	return &c
}

func (c *Config) LoadFromEnv() error {
	err := envconfig.Process("", c)
	return err
}
