package app

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/sirupsen/logrus"
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
		logrus.Fatalf("failed to configure database:%v ", err)
	}

	err = db.Ping()
	if err != nil {
		logrus.Fatalf("failed to connect to database:%v ", err)
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
		logrus.Fatalf("failed to query databases: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			logrus.Fatalf("failed to fetch row: %v", err)
		}
		logrus.Infof("found database: %v", name)
		databaseList = append(databaseList, name)
	}

	err = rows.Err()
	if err != nil {
		logrus.Fatalf("failed to fetch rows: %v", err)
	}

	return databaseList
}


