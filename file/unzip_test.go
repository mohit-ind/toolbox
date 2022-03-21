package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/c2fo/testify/require"
)

func copyFile(assert *require.Assertions, src, dest string) {
	input, err := ioutil.ReadFile(src)
	assert.NoError(err)

	err = ioutil.WriteFile(dest, input, 0644)
	assert.NoError(err)
}

func TestUnzip(t *testing.T) {
	assert := require.New(t)

	workingDir, err := ioutil.TempDir("", "zip-test")
	fmt.Println(workingDir)
	assert.NoError(err)
	defer func() {
		assert.NoError(os.RemoveAll(workingDir))
	}()

	zipFile := filepath.Join(workingDir, "test.zip")

	copyFile(assert, filepath.Join("test-data", "test.zip"), zipFile)

	files, err := Unzip(zipFile, workingDir)
	assert.NoError(err)
	assert.Equal([]string{
		filepath.Join(workingDir, "test-folder"),
		filepath.Join(workingDir, "test-folder", "test.txt")},
		files)
}

func TestUnzip_fail(t *testing.T) {
	assert := require.New(t)

	workingDir, err := ioutil.TempDir("", "zip-test")
	fmt.Println(workingDir)
	assert.NoError(err)
	defer func() {
		assert.NoError(os.RemoveAll(workingDir))
	}()

	zipFile := filepath.Join(workingDir, "test.zip")

	copyFile(assert, filepath.Join("test-data", "wrong.zip"), zipFile)

	files, err := Unzip(zipFile, workingDir)
	assert.EqualError(err, "zip: not a valid zip file")
	assert.Nil(files)
}

func TestUnzipAndRemoveOrigin(t *testing.T) {
	assert := require.New(t)

	workingDir, err := ioutil.TempDir("", "zip-test")
	fmt.Println(workingDir)
	assert.NoError(err)
	defer func() {
		assert.NoError(os.RemoveAll(workingDir))
	}()

	zipFile := filepath.Join(workingDir, "test.zip")

	copyFile(assert, filepath.Join("test-data", "test.zip"), zipFile)

	files, err := UnzipAndRemoveOrigin(zipFile, workingDir)
	assert.NoError(err)
	assert.Equal([]string{
		filepath.Join(workingDir, "test-folder"),
		filepath.Join(workingDir, "test-folder", "test.txt")},
		files)

	_, err = os.Stat(zipFile)
	assert.EqualError(err, fmt.Sprintf("stat %s: no such file or directory", zipFile))
}

func TestUnzipErrors(t *testing.T) {
	assert := require.New(t)

	files, err := Unzip("a", "b")
	assert.EqualError(err, "open a: no such file or directory")
	assert.Nil(files)

	files, err = UnzipAndRemoveOrigin("a", "b")
	assert.EqualError(err, "Failed to unzip file(s): open a: no such file or directory")
	assert.Nil(files)
}
