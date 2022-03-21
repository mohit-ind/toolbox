package csv

import (
	"testing"

	logger "github.com/toolboxlogger"
)

func TestCSVFields_Validate(t *testing.T) {
	log := logger.NewCommonLogger("test service", "csv.0.0.1", "test", "test", true)
	tests := []struct {
		name    string
		csvData CSVFields
		wantErr bool
	}{
		{
			name:    "when nothing is given",
			csvData: CSVFields{},
			wantErr: true,
		},
		{
			name: "when no data given",
			csvData: CSVFields{
				Logger:  log,
				Headers: []string{"A", "B"},
			},
			wantErr: true,
		},
		{
			name: "when invalid output path given",
			csvData: CSVFields{
				Logger:     log,
				Headers:    []string{"A", "B"},
				OutputPath: "////////../../....",
			},
			wantErr: true,
		},
		{
			name: "when invalid data given",
			csvData: CSVFields{
				Logger:  log,
				Headers: []string{"A", "B"},
				Data:    [][]string{{"1", "2", "3"}},
			},
			wantErr: true,
		},
		{
			name: "when everything is alright",
			csvData: CSVFields{
				Logger:   log,
				Filename: "validate-test-4.csv",
				Headers:  []string{"A", "B"},
				Data:     [][]string{{"1", "2"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.csvData.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("CSVFields.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCSVFields_GenerateCSVFile(t *testing.T) {
	log := logger.NewCommonLogger("test service", "csv.0.0.1", "test", "test", true)

	tests := []struct {
		name         string
		csvData      CSVFields
		wantRootPath bool
		wantErr      bool
	}{
		{
			name: "validation error",
			csvData: CSVFields{
				Logger:     log,
				Filename:   "test.go",
				OutputPath: "",
				Headers:    []string{"A"},
				Data:       [][]string{{"1", "2"}},
				ChunckSize: 2,
			},
			wantRootPath: false,
			wantErr:      true,
		},
		{
			name: "everything okay ;)",
			csvData: CSVFields{
				Logger:     log,
				Filename:   "test.csv",
				OutputPath: "",
				Headers:    []string{"A", "B"},
				Data:       [][]string{{"1", "2"}},
				ChunckSize: 2,
			},
			wantRootPath: true,
			wantErr:      false,
		},
		{
			name: "everything okay with bunch of data ;)",
			csvData: CSVFields{
				Logger:     log,
				Filename:   "test.csv",
				OutputPath: "",
				Headers:    []string{"A", "B"},
				Data:       [][]string{{"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}},
				ChunckSize: 10,
			},
			wantRootPath: true,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRootPath, _, err := tt.csvData.GenerateCSVFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("CSVFields.GenerateCSVFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotRootPath == "") == tt.wantRootPath {
				t.Errorf("CSVFields.GenerateCSVFile() gotRootPath = %v, want %v", gotRootPath, tt.wantRootPath)
			}
		})
	}
}
