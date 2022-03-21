// Package s3 provides a simple AWS S3 downloader
package aws

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

// s3Downloader is every object which is capable of downloading files from S3.
type s3Downloader interface {
	Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (n int64, err error)
}

// Downloader is responsible to download files from S3.
type Downloader struct {
	ctx context.Context
	api s3Downloader
}

// NewDownloader creates a new Downloader.
// It may return an error if the AWS configuration cannot be created, however it will create
// the Downloader with incorrect configs (it does not test access key and secret),
// so it will only fail when Download() is called.
func NewDownloader(ctx context.Context) (*Downloader, error) {
	awsConf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create AWS Configs")
	}

	return &Downloader{
		ctx: ctx,
		api: manager.NewDownloader(s3.NewFromConfig(awsConf)),
	}, nil
}

// Download gets a file from the target S3 bucket with target key,
// and puts it into a new file defined by filePath. It may return an error
// if it fails to create the output file, or fails to download the target file from S3.
func (dl *Downloader) Download(bucket string, key string, filePath string) error {
	// create outputFile
	outputFile, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "Failed to create download file: %s", filePath)
	}

	// download a file from S5 into the outputFile
	if _, err := dl.api.Download(dl.ctx, outputFile, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return errors.Wrap(err, "Failed to download file")
	}

	return outputFile.Close()
}

func (dl *Downloader) DownloadToMemory(bucket string, key string) ([]byte, error) {
	buff := manager.NewWriteAtBuffer([]byte{})

	// download a file from S5 into the buffer
	if _, err := dl.api.Download(dl.ctx, buff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return nil, errors.Wrap(err, "Failed to download file")
	}

	return buff.Bytes(), nil
}
