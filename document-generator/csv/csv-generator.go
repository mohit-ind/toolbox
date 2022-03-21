package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	logger "github.com/toolboxlogger"
	slice "github.com/toolboxslice"
)

type CSVFields struct {
	Logger               *logger.Logger
	Filename, OutputPath string
	Headers              []string
	Data                 [][]string
	ChunckSize           int
}

func (c *CSVFields) Validate() error {
	var validationError []string
	if c.Logger == nil {
		validationError = append(validationError, "nil logger found")
	}
	if c.Filename == "" || filepath.Ext(c.Filename) != ".csv" {
		validationError = append(validationError, "empty/invalid Filename found")
	}
	if _, err := os.Stat(c.OutputPath); c.OutputPath != "" && err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(c.OutputPath, os.ModePerm); err != nil {
				validationError = append(validationError, "invalid OutputPath or unable to create the path")
			}
		}
	}
	if len(c.Headers) == 0 {
		validationError = append(validationError, "no Headers found")
	} else if len(c.Data) == 0 {
		validationError = append(validationError, "no Data found")
	} else {
		for count, dataSet := range c.Data {
			if len(dataSet) != len(c.Headers) {
				validationError = append(validationError, fmt.Sprintf("row [%v] fields are not equal to headers", count+1))
			}
		}
	}
	if c.ChunckSize < 10 || c.ChunckSize >= 10000 {
		c.ChunckSize = 9999
	}
	if len(validationError) != 0 {
		return errors.New("csv generator error : " + strings.Join(validationError, "; "))
	}
	return nil
}

func (c *CSVFields) GenerateCSVFile() (rootPath string, filenames []string, err error) {
	if c.OutputPath == "" {
		c.OutputPath = filepath.Join("documents", "csv")
	}
	if err = c.Validate(); err != nil {
		return
	}

	filenames = []string{}

	dataSets := slice.SplitToChunks(c.Data, c.ChunckSize).([][][]string)
	c.Logger.Entry().Infof("csv generator : working with %v data per file", c.ChunckSize)
	for fileNum, data := range dataSets {
		var filename string
		if len(dataSets) > 1 {
			file := strings.Split(c.Filename, ".csv")
			filename = strings.Join(file[:len(file)-1], "") + fmt.Sprintf("-%v.csv", fileNum+1)
		} else {
			filename = c.Filename
		}
		filename = filepath.Join(c.OutputPath, filename)

		file, createCSVErr := os.Create(filename)
		if createCSVErr != nil {
			c.Logger.WithField("file-name", filename).WithError(createCSVErr).Error("failed to create csv file")
			err = createCSVErr
			return
		}
		defer func() {
			if fileErr := file.Close(); fileErr != nil {
				c.Logger.WithField("file-name", filename).WithError(fileErr).Error("failed to close csv file")
			}
		}()

		w := csv.NewWriter(file)
		defer w.Flush()

		var csvData [][]string
		csvData = append(csvData, c.Headers)
		csvData = append(csvData, data...)

		if err = w.WriteAll(csvData); err != nil {
			c.Logger.WithField("file-name", filename).WithError(err).Error("failed to write in csv file")
			return
		} else {
			filenames = append(filenames, c.Filename)
		}
	}
	rootPath, err = filepath.Abs(c.OutputPath)
	return
}
