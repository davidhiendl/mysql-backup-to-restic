package app

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/davidhiendl/mysql-backup-to-s3/app/logger"
	"os"
	"bytes"
	"net/http"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

func (app *App) connectToS3() {
	var token = ""

	creds := credentials.NewStaticCredentials(app.config.S3AccessKey, app.config.S3SecretKey, token)
	_, err := creds.Get()
	if err != nil {
		logger.Fatalf("failed to configure s3 credentials: %v \n", err)
	}

	conf := aws.NewConfig().WithRegion(app.config.S3Region).WithEndpoint(app.config.S3Endpoint).WithCredentials(creds)
	svc := s3.New(session.New(), conf)
	app.s3svc = svc
}

func (app *App) storeFile(localPath string, s3Path string) {
	file, err := os.Open(localPath)
	if err != nil {
		logger.Fatalf("err opening file: %s", err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	params := &s3.PutObjectInput{
		Bucket:        aws.String(app.config.S3Bucket),
		Key:           aws.String(s3Path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	resp, err := app.s3svc.PutObject(params)
	if err != nil {
		logger.Fatalf("bad response: %s", err)
	}
	logger.Infof("stored file %v to %v, response: ", localPath, s3Path, awsutil.StringValue(resp))
}
