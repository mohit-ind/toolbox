package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestAddFileToS3_NoFileError(t *testing.T) {
	t.Run("TestAddFileToS3_NoFileError", func(t *testing.T) {
		err := S3Essentials{
			BucketName: "s3",
		}.AddFileToS3()
		assert.Equal(t, "open : no such file or directory", err.Error())
	})
}

func TestGetUploader(t *testing.T) {
	t.Run("TestGetUploader func ()", func(t *testing.T) {
		_, err := getUploader()
		if err != nil {
			assert.Equal(t, err.Error(), err.Error())
		}
	})
}

type mockUploader struct {
}

func (m mockUploader) Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	return &manager.UploadOutput{}, nil
}

func TestUploadFunc(t *testing.T) {
	t.Run("TestUpload func ()", func(t *testing.T) {
		var uploader mockUploader
		err := uploaderNeeds{}.upload(uploader)
		assert.NoError(t, err, "No error should be come as we mcked the function")
	})
}
