package file

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filePaths []string

	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer zipReader.Close()

	for _, zipFile := range zipReader.File {

		// Store filename/path for returning and using later on
		zipFilePath := filepath.Join(dest, zipFile.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(zipFilePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return nil, errors.Errorf("%s: illegal file path", zipFilePath)
		}

		filePaths = append(filePaths, zipFilePath)

		if zipFile.FileInfo().IsDir() {
			// Make Folder
			if err := os.MkdirAll(zipFilePath, os.ModePerm); err != nil {
				return nil, err
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(zipFilePath), os.ModePerm); err != nil {
			return nil, err
		}

		outFile, err := os.OpenFile(zipFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
		if err != nil {
			return nil, err
		}

		zipFileReader, err := zipFile.Open()
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(outFile, zipFileReader); err != nil {
			return nil, err
		}

		// Close the file without defer to close before next iteration of loop
		if err := outFile.Close(); err != nil {
			return nil, err
		}

		if err := zipFileReader.Close(); err != nil {
			return nil, err
		}

	}

	return filePaths, nil
}

// UnzipAndRemoveOrigin will call Unzip, and after it finishes, removes the original zip file.
func UnzipAndRemoveOrigin(src string, dest string) ([]string, error) {
	files, err := Unzip(src, dest)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unzip file(s)")
	}
	if err := os.Remove(src); err != nil {
		return nil, errors.Wrapf(err, "Failed to remove source zip file: %s", src)
	}
	return files, nil
}
