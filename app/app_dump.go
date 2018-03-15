package app

import (
	"os/exec"
	"github.com/davidhiendl/mysql-backup-to-s3/app/logger"
	"io/ioutil"
	"strconv"
	"github.com/JamesStewy/go-mysqldump"
	"time"
	"compress/gzip"
	"os"
)

func (app *App) getDumpTime() string {
	t := time.Now()
	return t.Format("2006-01-02T15_04-05-9999") + t.Format("Z07:00")
}

func (app *App) dumpDatabaseMysqldump(dbName string) string {
	outPath := app.config.ExportTempDir + "/" + dbName + ".sql"

	// append extension if gz compression is configured
	if app.config.CompressWithGz {
		outPath += ".gz"
	}

	logger.Infof("dumping database: %v to file: %v \n", dbName, outPath)

	// find binary
	binary, err := exec.LookPath(app.config.MySQLDumpBinary)
	if err != nil {
		logger.Fatalf("failed to find binary: %v, err: %v \n", app.config.MySQLDumpBinary, err)
	}

	// output file
	outFile, err := os.Create(outPath)
	if err != nil {
		logger.Fatalf("failed to create command: %v", err)
	}
	defer outFile.Close()

	// create command
	cmd := exec.Command(binary, "-P"+strconv.Itoa(app.config.MySQLPort), "-h"+app.config.MySQLHost, "-u"+app.config.MySQLUser, "-p"+app.config.MySQLPass, dbName)

	// use gz compression if configured
	if app.config.CompressWithGz {
		// create gzWriter
		gzWriter := gzip.NewWriter(outFile)
		defer gzWriter.Flush()
		defer gzWriter.Close()

		// attach compressed writer to gz command
		cmd.Stdout = gzWriter
	} else {
		cmd.Stdout = outFile
	}

	// attach stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Fatalf("failed to create command: %v", err)
	}

	// execute
	if err := cmd.Start(); err != nil {
		logger.Fatalf("failed to start command: %v", err)
	}

	// read stderr
	stderrorBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		logger.Fatalf("failed to read from stderr: %v", err)
	}
	println("mysqldump: " + string(stderrorBytes))

	return outPath
}

func (app *App) dumpDatabaseInGo(dbName string) string {
	logger.Infof("dumping db %v", dbName)

	db := app.connectToDb(dbName)

	// Register database with mysqldump
	dumper, err := mysqldump.Register(db, app.config.ExportTempDir, time.ANSIC)
	if err != nil {
		logger.Fatalf("error registering db for dumb: %v \n", err)
	}

	outFile, err := dumper.Dump()
	if err != nil {
		logger.Fatalf("error dumping db for dumb: %v \n", err)
	}

	logger.Infof("dumped db %v to %v", dbName, outFile)

	return outFile
}
