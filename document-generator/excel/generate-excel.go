package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	slice "github.com/toolboxslice"

	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
)

type ExcelFields struct {
	Headers          []string
	Rows             [][]string
	FileName         string
	FileNameExist    bool
	SheetName        string
	DeleteFromServer bool
}

func (x ExcelFields) CreateAndSaveExcelFile(outputPath string) (err error) {

	if len(x.Headers) == 0 {
		return errors.New("error : headers length must be greater than zero (0)")
	} else if len(x.Rows) == 0 {
		return errors.New("error : rows must be must be greater than zero (0)")
	} else if strings.TrimSpace(x.FileName) == "" {
		return errors.New("error : file name doesn't exist")
	} else if len(x.FileName) < 5 || x.FileName[len(x.FileName)-5:] != ".xlsx" {
		x.FileName += ".xlsx"
	}
	if x.FileNameExist && x.SheetName == "" {
		return errors.New("error : SheetName can't be empty for already exist file")
	}

	if outputPath == "" {
		return errors.New("error : please specify path")
	}

	outputPath = filepath.FromSlash(outputPath)
	if outputPath[len(outputPath)-1:] != "/" {
		outputPath += "/"
	}

	outputPath = filepath.Join(outputPath + x.FileName)
	outputPath = filepath.FromSlash(outputPath)
	if x.SheetName == "" {
		x.SheetName = "Sheet"
	}

	var f *excelize.File
	if x.FileNameExist {
		if f, err = excelize.OpenFile(outputPath); err != nil {
			return
		}
	} else {
		f = excelize.NewFile()
	}

	if err = generator(x.Headers, x.Rows, x.SheetName, outputPath, f); err != nil {
		return
	}
	if x.DeleteFromServer {
		return os.Remove(outputPath)
	}
	return nil
}

func generator(headers []string, xRows [][]string, xSheetName, path string, f excelizeInterface) (err error) {
	style, styleErr := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Bold: true},
	})
	if styleErr != nil {
		return styleErr
	}
	allRows := slice.SplitToChunks(xRows, 9999).([][][]string)
	for sNo, rows := range allRows {
		sheetName := fmt.Sprintf(xSheetName+"%v", sNo+1)
		index := f.NewSheet(sheetName)
		row := 1
		err = f.SetCellStyle(sheetName, "A1", getAlphaFromDigit(int32(len(headers)))+"1", style)
		if err != nil {
			return
		}
		err = f.SetSheetRow(sheetName, "A1", &headers)
		if err != nil {
			return
		}

		for i := 2; row <= len(rows); i++ {
			if len(headers) != len(rows[row-1]) {
				return errors.New("error : number of headers not equals to no of row")
			}
			_ = f.SetSheetRow(sheetName, "A"+strconv.Itoa(i), &rows[row-1])
			row++
		}
		f.SetActiveSheet(index)
	}
	return f.SaveAs(path)
}

func getAlphaFromDigit(number int32) (letters string) {
	number--
	if firstLetter := number / 26; firstLetter > 0 {
		letters += getAlphaFromDigit(firstLetter)
		letters += string('A' + number%26)
	} else {
		letters += string('A' + number)
	}
	return
}

type excelizeInterface interface {
	NewSheet(name string) int
	NewStyle(style interface{}) (int, error)
	SaveAs(name string, opt ...excelize.Options) error
	SetCellStyle(sheet string, hcell string, vcell string, styleID int) error
	SetSheetRow(sheet string, axis string, slice interface{}) error
	SetActiveSheet(index int)
}
