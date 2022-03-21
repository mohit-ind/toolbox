package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/c2fo/testify/require"
)

type mockApi struct {
	bytes int64
	err   string
}

func (m *mockApi) Download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*manager.Downloader)) (n int64, err error) {
	if m.err != "" {
		err = errors.New(m.err)
	}
	return m.bytes, err
}

func TestNewDownloader(t *testing.T) {
	assert := require.New(t)

	dl, err := NewDownloader(context.TODO())
	assert.NoError(err)
	assert.NotNil(dl)
}

func TestMockDownload(t *testing.T) {
	assert := require.New(t)

	dl := &Downloader{
		ctx: context.TODO(),
		api: &mockApi{},
	}

	testFile := filepath.Join(os.TempDir(), time.Now().Format("20060102150405")+"-s3-dl-test")

	err := dl.Download("", "", testFile)

	assert.NoError(err)

	_, err = os.Stat(testFile)

	assert.NoError(err)

	assert.NoError(os.Remove(testFile))
}

func TestMockDownload_fail(t *testing.T) {
	assert := require.New(t)

	dl := &Downloader{
		ctx: context.TODO(),
		api: &mockApi{
			err: "S3 Service is down!",
		},
	}

	testFile := filepath.Join(os.TempDir(), time.Now().Format("20060102150405")+"-s3-dl-test")

	err := dl.Download("", "", testFile)

	assert.EqualError(err, "Failed to download file: S3 Service is down!")

	_, err = os.Stat(testFile)

	assert.NoError(err)

	assert.NoError(os.Remove(testFile))
}

func TestMockDownload_cannot_create_file(t *testing.T) {
	assert := require.New(t)

	dl := &Downloader{
		ctx: context.TODO(),
		api: &mockApi{},
	}

	testFile := filepath.Join(os.TempDir(), time.Now().Format("20060102150405")+"-s3-dl-test")

	assert.NoError(os.MkdirAll(testFile, os.ModePerm))

	err := dl.Download("", "", testFile)

	assert.EqualError(err, fmt.Sprintf(
		"Failed to create download file: %s: open %s: is a directory",
		testFile,
		testFile))

	assert.NoError(os.Remove(testFile))
}

func TestCIDownload(t *testing.T) {
	if os.Getenv("CI") != "true" {
		t.Skipf("Skipping S3 Download test because not running in CI")
	}

	assert := require.New(t)

	dl, err := NewDownloader(context.TODO())

	assert.NoError(err)

	testFile := filepath.Join(os.TempDir(), time.Now().Format("20060102150405")+"-s3-dl-test")

	assert.NoError(dl.Download("appventurez-images", "beer_parrot.gif", testFile))

	_, err = os.Stat(testFile)

	assert.NoError(err)

	assert.NoError(os.Remove(testFile))
}

func TestCIDownloadToMemory(t *testing.T) {
	// if os.Getenv("CI") != "true" {
	// 	t.Skipf("Skipping S3 Download test because not running in CI")
	// }

	assert := require.New(t)

	dl, err := NewDownloader(context.TODO())

	assert.NoError(err)

	data, err := dl.DownloadToMemory("appventurez-images", "beer_parrot.gif")
	assert.NoError(err)
	assert.True(len(data) > 100)

	data, err = dl.DownloadToMemory("invalid-bucket-name", "beer_parrot.gif")
	assert.Contains(err.Error(), "Failed to download file: operation error S3: GetObject, https response error StatusCode: 301")
	assert.Nil(data)
}
