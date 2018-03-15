package app

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/davidhiendl/mysql-backup-to-s3/app/logger"
	"fmt"
)

func (app *App) connectToDb(name string) *sql.DB {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
			app.config.MySQLUser,
			app.config.MySQLPass,
			app.config.MySQLHost,
			app.config.MySQLPort,
			name))

	if err != nil {
		logger.Fatalf("failed to configure database: %v \n", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatalf("failed to connect to database: %v \n", err)
	}

	return db
}

func (app *App) getDatabases() []string {
	var databaseList []string

	var (
		name string
	)

	rows, err := app.db.Query("show databases")
	if err != nil {
		logger.Fatalf("failed to query databases: %v \n", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			logger.Fatalf("failed to fetch row: %v \n", err)
		}
		logger.Infof("listed db: %v \n", name)
		databaseList = append(databaseList, name)
	}

	err = rows.Err()
	if err != nil {
		logger.Fatalf("failed to fetch rows: %v \n", err)
	}

	return databaseList
}


