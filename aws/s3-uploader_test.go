package aws

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	file "github.com/toolboxfile"
)

func TestNewUploader(t *testing.T) {
	assert := require.New(t)

	ul, err := NewUploader(context.TODO())
	assert.NoError(err)
	assert.IsType(&Uploader{}, ul)
}

func TestUploaderUpload(t *testing.T) {
	if os.Getenv("CI") != "true" {
		t.Skipf("Skipping S3 Upload test because not running in CI")
	}

	assert := require.New(t)

	ul, err := NewUploader(context.TODO())
	assert.NoError(err)

	testFile, rm := file.NewTestFile(t, os.ModePerm, "test")
	defer rm()

	assert.NoError(ul.Upload("appventurez-devops-test", "s3-upload-test/testfile", testFile))
}

func TestUploaderUploadFromMemory(t *testing.T) {
	if os.Getenv("CI") != "true" {
		t.Skipf("Skipping S3 Upload test because not running in CI")
	}

	assert := require.New(t)

	ul, err := NewUploader(context.TODO())
	assert.NoError(err)

	assert.NoError(ul.UploadFromMemory(
		"appventurez-devops-test",
		"s3-upload-test/testfile-from-memory",
		[]byte("memory-test"),
	))
}
