package app

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"bytes"
	"net/http"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/sirupsen/logrus"
)

func (app *App) connectToS3() {
	var token = ""

	creds := credentials.NewStaticCredentials(app.config.S3AccessKey, app.config.S3SecretKey, token)
	_, err := creds.Get()
	if err != nil {
		logrus.Fatalf("failed to configure s3 credentials: %v \n", err)
	}

	conf := aws.NewConfig()
	conf.WithRegion(app.config.S3Region)
	conf.WithEndpoint(app.config.S3Endpoint)
	conf.WithCredentials(creds)
	conf.WithS3ForcePathStyle(app.config.S3ForcePathStyle)

	svc := s3.New(session.New(), conf)
	app.s3svc = svc
}

func (app *App) storeFile(localPath string, s3Path string) {
	file, err := os.Open(localPath)
	if err != nil {
		logrus.Fatalf("err opening file: %s", err)
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
		logrus.Fatalf("bad response: %s", err)
	}
	logrus.Infof("stored file %v to %v, response: ", localPath, s3Path, awsutil.StringValue(resp))
}
