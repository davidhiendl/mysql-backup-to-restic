package app

import (
	"database/sql"
	"github.com/davidhiendl/mysql-backup-to-restic/app/config"
	"github.com/davidhiendl/mysql-backup-to-restic/app/restic"
	"os"
	"path/filepath"
)

type App struct {
	config  *config.Config
	db      *sql.DB
	restic  *restic.Restic
	dumpDir string
}

// Create new config and populate it from environment
func NewApp(config *config.Config) (*App) {
	app := App{
		config:  config,
		dumpDir: filepath.Join(config.Common.ScratchDir, "sqldumps"),
	}

	app.db = app.connectToDb("")
	app.restic = restic.New(config)

	return &app
}

// run templates against containers and generate config
func (app *App) Run() {
	databases := app.getDatabases()

	dumpedFiles := make([]string, 0)

	// export databases
	for db, ok := range databases {
		if !ok {
			continue
		}

		// execute dump
		dumpFile := app.dumpDatabaseMysqldump(db)
		dumpedFiles = append(dumpedFiles, dumpFile)
	}

	// run restic
	app.restic.InitRepositoryIfAbsent()
	app.restic.Backup(app.dumpDir, nil);
	if (app.config.RetentionPolicy.HasKeepDirective()) {
		app.restic.Forget(&app.config.RetentionPolicy)

		if (app.config.RetentionPolicy.Check) {
			app.restic.Check()
		}
	}

	// clean up
	for _, file := range dumpedFiles {
		os.Remove(file)
	}
}

func (app *App) ShouldIncludeDatabase(name string) (bool, string) {
	// exclude system
	if app.config.Databases.ExcludeSystem {
		if name == "performance_schema" || name == "information_schema" || name == "mysql" {
			return false, "system database excluded"
		}
	}

	// exclude specific list
	for _, exclude := range app.config.Databases.Exclude {
		if name == exclude {
			return false, "on exclude list"
		}
	}

	// if include list enabled
	if len(app.config.Databases.Include) > 0 {
		for _, include := range app.config.Databases.Include {
			if name == include {
				return true, "on include list"
			}
		}

		// only allow databases that are included
		return false, "not on include list"
	}

	return true, "not on any list"
}

func (app *App) Close() {
	app.db.Close()
}
