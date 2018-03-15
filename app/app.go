package app

import (
	"database/sql"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/davidhiendl/mysql-backup-to-s3/app/logger"
	"os"
)

type App struct {
	config *Config
	db     *sql.DB
	s3svc  *s3.S3
}

// Create new config and populate it from environment
func NewApp(config *Config) (*App) {
	app := App{
		config: config,
	}

	app.db = app.connectToDb("")
	app.connectToS3()

	return &app
}

// run templates against containers and generate config
func (app *App) Run() {
	databases := app.getDatabases()

	for _, db := range databases {
		if app.shouldSkipDb(db) {
			logger.Infof("skipping db: %v", db)
			continue
		}

		// TODO improve workflow to skip temporary disk storage of backup files

		// determine paths
		dumpFile := app.dumpDatabaseMysqldump(db)
		s3Path := app.config.S3PathPrefix + "/" + db + "/" + app.getDumpTime() + ".sql"

		// append extension if gz compression is configured
		if app.config.CompressWithGz {
			s3Path += ".gz"
		}

		// store file to s3
		app.storeFile(dumpFile, s3Path)

		// cleanup file from local storage
		os.Remove(dumpFile)
	}

}

func (app *App) shouldSkipDb(name string) bool {
	if app.config.SkipSystemDatabases {
		if name == "performance_schema" || name == "information_schema" || name == "mysql" {
			return true
		}
	}

	// TODO skip by regex

	return false
}

func (app *App) Close() {
	app.db.Close()
}
