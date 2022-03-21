package excel

import (
	"os"
	"strconv"
	"testing"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_getAlphaFromDigit(t *testing.T) {
	tests := []struct {
		name        string
		digit       int32
		wantLetters string
	}{
		{
			name:        "1 -> A",
			digit:       1,
			wantLetters: "A",
		},
		{
			name:        "123 -> DS",
			digit:       123,
			wantLetters: "DS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLetters := getAlphaFromDigit(tt.digit); gotLetters != tt.wantLetters {
				t.Errorf("getAlphaFromDigit() = %v, want %v", gotLetters, tt.wantLetters)
			}
		})
	}
}

func TestExcelFields_CreateAndSaveExcelFile(t *testing.T) {

	tests := []struct {
		Case                string
		fields              ExcelFields
		pathToSaveExcelFile string

		wantErr bool
	}{
		{
			Case: "happy case ;) when file name is correct or exist",
			fields: ExcelFields{
				Headers:  []string{"A", "B"},
				Rows:     [][]string{{"1", "2"}},
				FileName: "test-excel.xlsx",
			},
			pathToSaveExcelFile: "",
			wantErr:             false,
		},
		{
			Case: "happy case ;) if user forgot to add .xlsx then function add it by itself",
			fields: ExcelFields{
				Headers:          []string{"A", "B"},
				Rows:             [][]string{{"1", "2"}},
				FileName:         "test-excel",
				FileNameExist:    true,
				SheetName:        "test2",
				DeleteFromServer: true,
			},
			pathToSaveExcelFile: "",
			wantErr:             false,
		},
		{
			Case: "error of mkdir",
			fields: ExcelFields{
				Headers:          []string{"A", "B"},
				Rows:             [][]string{{"1", "2"}},
				FileName:         "test-excel",
				FileNameExist:    true,
				SheetName:        "test2",
				DeleteFromServer: true,
			},
			pathToSaveExcelFile: "qwe/123///",
			wantErr:             true,
		},
		{
			Case: "error while sheet name doesn't exist for already exists file",
			fields: ExcelFields{
				FileName:         "test-excel.xlsx",
				Headers:          []string{"A", "B"},
				Rows:             [][]string{{"1", "2"}},
				FileNameExist:    true,
				DeleteFromServer: true,
			},
			pathToSaveExcelFile: "",
			wantErr:             true,
		},
		{
			Case: "error case : when fileExist true but same file doesn't exist",
			fields: ExcelFields{
				Headers:       []string{"A", "B"},
				Rows:          [][]string{{"1", "2"}},
				FileName:      "test-excel2.xlsx",
				FileNameExist: true,
				SheetName:     "beta",
			},
			pathToSaveExcelFile: "",
			wantErr:             true,
		},
		{
			Case: "when no headers provided",
			fields: ExcelFields{
				FileName:         "test-excel.xlsx",
				DeleteFromServer: true,
			},
			pathToSaveExcelFile: "",
			wantErr:             true,
		},
		{
			Case: "when data-rows doesn't exists",
			fields: ExcelFields{
				Headers:          []string{"A", "B"},
				FileName:         "test-excel.xlsx",
				DeleteFromServer: true,
			},
			pathToSaveExcelFile: "",
			wantErr:             true,
		},
		{
			Case: "when file name doesn't exists",
			fields: ExcelFields{
				Headers:  []string{"A", "B"},
				Rows:     [][]string{{"1", "2"}},
				FileName: " ",
			},
			wantErr: true,
		},
		{
			Case: "when improper length given for rows or not according to headers",
			fields: ExcelFields{
				Headers:          []string{"A", "B"},
				Rows:             [][]string{{"1"}},
				FileName:         "test-excel.xlsx",
				DeleteFromServer: true,
			},
			pathToSaveExcelFile: "",
			wantErr:             true,
		},
		{
			Case: "error : file name doesn't exist",
			fields: ExcelFields{
				Headers: []string{"A", "B"},
				Rows:    [][]string{{"1", "2"}},
			},
			wantErr: true,
		},
		{
			Case: "happy case ;) with large no. of rows.",
			fields: ExcelFields{
				Headers:  []string{"A", "B"},
				Rows:     [][]string{{"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}, {"1", "2"}},
				FileName: "test-excel.xlsx",
			},
			pathToSaveExcelFile: "",
			wantErr:             false,
		},
		{
			Case: "while batch have more than 10K rows",
			fields: ExcelFields{
				Headers:  []string{"A", "B"},
				Rows:     [][]string{{"1", "2"}},
				FileName: "test-excel.xlsx",
			},
			pathToSaveExcelFile: "",
			wantErr:             false,
		},
		{
			Case: "error while output path in not there",
			fields: ExcelFields{
				Headers:  []string{"A", "B"},
				Rows:     [][]string{{"1", "2"}},
				FileName: "test-excel.xlsx",
			},
			pathToSaveExcelFile: "",
			wantErr:             true,
		},
	}
	path, err := os.Getwd()
	assert.NoError(t, err, "error shouldn't be occur")
	for _, tt := range tests {
		if tt.Case == "while batch have more than 10K rows" {
			for i := 0; i < 11000; i++ {
				tt.fields.Rows = append(tt.fields.Rows, []string{strconv.Itoa(i + 1), "test more than 10K"})
			}
		}
		if tt.Case != "error while output path in not there" {
			if tt.pathToSaveExcelFile == "" {
				tt.pathToSaveExcelFile = path
			}
		}
		t.Run(tt.Case, func(t *testing.T) {
			if err := tt.fields.CreateAndSaveExcelFile(tt.pathToSaveExcelFile); (err != nil) != tt.wantErr {
				t.Errorf("[%s]ExcelFields.CreateAndSaveExcelFile() error = %v, wantErr %v", tt.Case, err, tt.wantErr)
			}
		})
	}
}

func Test_excelizeGeneratorError(t *testing.T) {
	testCases := []struct {
		name string
		Case int
	}{
		{
			name: "styleErr",
			Case: 0,
		},

		{
			name: "SetCellStyleErr",
			Case: 1,
		},

		{
			name: "SetSheetRowError",
			Case: 2,
		},
	}
	for _, tt := range testCases {
		t.Run("Test_excelizeStyleError", func(t *testing.T) {
			v := excelizeErrDescriber{
				Case: tt.Case,
			}
			if err := generator([]string{"A", "B"}, [][]string{{"1", "2"}}, "", "", v); (err != nil) != true {
				t.Errorf("[%v]Test_excelizeGeneratorError error = %v", tt.name, err)
			}
		})
	}
}

type excelizeErrDescriber struct {
	Case int
}

func (e excelizeErrDescriber) NewStyle(style interface{}) (int, error) {
	switch e.Case {
	case 1, 2:
		return e.Case, nil
	}
	return 0, errors.New("style_error generated for test")
}
func (e excelizeErrDescriber) NewSheet(name string) int {
	return 1
}

func (e excelizeErrDescriber) SaveAs(name string, opt ...excelize.Options) error {
	return errors.New("error generated for SaveAs")
}
func (e excelizeErrDescriber) SetCellStyle(sheet string, hcell string, vcell string, styleID int) error {
	switch e.Case {
	case 2:
		return nil
	}
	return errors.New("error generated for SetCellStyle")

}
func (e excelizeErrDescriber) SetSheetRow(sheet string, axis string, slice interface{}) error {
	return errors.New("error generated for SetSheetRow")
}
func (e excelizeErrDescriber) SetActiveSheet(index int) {

}
