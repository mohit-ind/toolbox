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
	"github.com/pkg/errors"
)

// s3Uploader is every object which is capable of uploading files to S3.
type s3Uploader interface {
	Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error)
}

// Uploader is responsible for uploading files to S3.
type Uploader struct {
	ctx context.Context
	api s3Uploader
}

// NewUploader creates a new Uploader.
// It may return an error if the AWS configuration cannot be created, however it will create
// the Uploader with incorrect configs (it does not test access key and secret),
// so it will only fail when Upload() is called.
func NewUploader(ctx context.Context) (*Uploader, error) {
	awsConf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create AWS Configs")
	}

	return &Uploader{
		ctx: ctx,
		api: manager.NewUploader(s3.NewFromConfig(awsConf)),
	}, nil
}

// Upload a file from filePath to target S3 bucket/key.
func (ul *Uploader) Upload(bucket string, key string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	// Get file size and read the file content into a buffer.
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	if _, err = file.Read(buffer); err != nil {
		return err
	}

	// upload the content of the buffer to the target S3 bucket/key.
	if _, err := ul.api.Upload(ul.ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          bytes.NewBuffer(buffer),
		ContentLength: size,
		ContentType:   aws.String(http.DetectContentType(buffer)),
	}); err != nil {
		return errors.Wrap(err, "Failed to upload file")
	}

	return file.Close()
}

// Upload create a file in S3 bucket/key with data from memory.
func (ul *Uploader) UploadFromMemory(bucket string, key string, data []byte) error {
	// upload the data to the target S3 bucket/key.
	if _, err := ul.api.Upload(ul.ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          bytes.NewBuffer(data),
		ContentLength: int64(len(data)),
		ContentType:   aws.String(http.DetectContentType(data)),
	}); err != nil {
		return errors.Wrap(err, "Failed to upload data from memory")
	}
	return nil
}
