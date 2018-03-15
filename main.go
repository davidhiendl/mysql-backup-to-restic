package main

import (
	"github.com/davidhiendl/mysql-backup-to-s3/app"
	"github.com/davidhiendl/mysql-backup-to-s3/app/logger"
	"github.com/fatih/structs"
)

func main() {

	// retrieve config
	config := app.NewConfig()
	err := config.LoadFromEnv()
	if err != nil {
		logger.Fatalf("failed to parse configuration from environment: %v \n", err)
	}

	logger.SetLevel(config.LogLevel)

	m := structs.Map(config)
	for key, value := range m {
		logger.Infof("Config.%v = %v", key, value)
	}

	instance := app.NewApp(config)
	instance.Run()
}
