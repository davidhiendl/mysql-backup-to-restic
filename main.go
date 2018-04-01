package main

import (
	"github.com/davidhiendl/mysql-backup-to-s3/app"
	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
	"os"
	"github.com/davidhiendl/mysql-backup-to-s3/app/config"
)

func main() {


	// configure logger
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// retrieve config
	cfg := config.Load()
	// print config
	m := structs.Map(cfg)
	for key, value := range m {
		if key == "EnvMap" {
			continue
		}
		logrus.WithFields(logrus.Fields{"key": key, "value": value}).Infof("configuration loaded")
	}


	instance := app.NewApp(cfg)
	instance.Run()
}
