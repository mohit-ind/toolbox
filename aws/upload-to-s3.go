package aws

import (
	"bytes"
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Essentials struct {
	BucketName string
	FileDir    string
}

// AddFileToS3
func (e S3Essentials) AddFileToS3() error {
	file, err := os.Open(e.FileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	if _, err = file.Read(buffer); err != nil {
		return err
	}
	uploader, err := getUploader()
	if err != nil {
		return err
	}
	return uploaderNeeds{
		buffer:   buffer,
		size:     size,
		bucket:   e.BucketName,
		filePath: e.FileDir,
	}.upload(uploader)
}

type uploaderNeeds struct {
	buffer   []byte
	size     int64
	bucket   string
	filePath string
}

func (u uploaderNeeds) upload(uploader AWS) error {
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:               aws.String(u.bucket),
		Key:                  aws.String(u.filePath),
		Body:                 bytes.NewReader(u.buffer),
		ContentLength:        u.size,
		ContentType:          aws.String(http.DetectContentType(u.buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: "AES256",
	})
	return err
}

func getUploader() (*manager.Uploader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return manager.NewUploader(s3.NewFromConfig(cfg)), nil
}
