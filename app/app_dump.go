package app

import (
	"os/exec"
	"io/ioutil"
	"strconv"
	"compress/gzip"
	"os"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

const (
	MYSQLDUMP_BINARY = "mysqldump"
)

func (app *App) dumpDatabaseMysqldump(dbName string) string {

	err := os.MkdirAll(app.dumpDir, 0750)
	if err != nil {
		logrus.Fatalf("failed to create dump dir: %v, error: %+v", app.dumpDir, err)
	}

	outPath := filepath.Join(app.dumpDir, dbName+".sql")

	// append extension if gz compression is configured
	if app.config.Dump.CompressWithGz {
		outPath += ".gz"
	}

	logrus.Infof("dumping database: %v to file: %v \n", dbName, outPath)

	// find binary
	binary, err := exec.LookPath(MYSQLDUMP_BINARY)
	if err != nil {
		logrus.Fatalf("failed to find mysqldump binary: %v, err: %v \n", MYSQLDUMP_BINARY, err)
	}

	// output file
	outFile, err := os.Create(outPath)
	if err != nil {
		logrus.Fatalf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// create command
	cmd := exec.Command(binary, "-P"+strconv.Itoa(app.config.MySQL.Port), "-h"+app.config.MySQL.Host, "-u"+app.config.MySQL.Username, "-p"+app.config.MySQL.Password, dbName)

	// use gz compression if configured
	if app.config.Dump.CompressWithGz {
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
		logrus.Fatalf("failed to create command: %v", err)
	}

	// execute
	if err := cmd.Start(); err != nil {
		logrus.Fatalf("failed to start command: %v", err)
	}

	// read stderr
	stderrorBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		logrus.Fatalf("failed to read from stderr: %v", err)
	}
	println("mysqldump: " + string(stderrorBytes))

	return outPath
}
