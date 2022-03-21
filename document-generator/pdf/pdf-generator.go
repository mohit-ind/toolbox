package pdf

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"github.com/pkg/errors"
	config "github.com/toolboxconfig"
	logger "github.com/toolboxlogger"

	HTMLToPDF "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type NewPDFRequest struct {
	TemplateHTML     string
	TemplateData     interface{}
	FileName         string
	OutputPath       string
	DeleteFromServer bool
	Logger           *logger.Logger
	WKHTMLTOPDF_PATH string
}

func (pdf *NewPDFRequest) GenerateHTMLToPDF() (err error) {
	fullpath := filepath.Join(pdf.OutputPath, pdf.FileName)

	if pdf.Logger == nil {
		pdf.Logger = logger.NewCommonLogger("toolbox", "latest", "", config.GetHostName(), true)
	}
	if pdf.FileName == "" || filepath.Ext(pdf.FileName) != ".pdf" {
		pdf.Logger.WithField("path to save pdf", fullpath).Error("incorrect file-name provided")
		return errors.New("error : empty/invalid file name, please provide correct file-name with .pdf extension")
	} else if pdf.OutputPath == "" {
		pdf.Logger.WithField("path to save pdf", fullpath).Error("incorrect output path found")
		return errors.New("error : please provide outpath")
	}
	pdf.Logger.Entry().Infof("path to save pdf is : [%v]", fullpath)
	t, tmplErr := template.New(pdf.FileName).Parse(pdf.TemplateHTML)
	if tmplErr != nil {
		pdf.Logger.WithField("path to save pdf", fullpath).WithError(tmplErr).Error("error while creating template from data")
		return tmplErr
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, pdf.TemplateData); err != nil {
		pdf.Logger.WithField("path to save pdf", fullpath).WithError(err).Error("executing template to generate data")
		return err
	}
	HTMLContent := buf.String()
	htmlFile := strconv.FormatInt(int64(time.Now().UnixNano()), 10) + ".html"
	defer func() {
		if f, statErr := os.Stat(htmlFile); statErr == nil && htmlFile != "" && f.Name() == htmlFile {
			os.Remove(htmlFile)
		}
	}()
	if err = ioutil.WriteFile(htmlFile, []byte(HTMLContent), 0644); err != nil {
		pdf.Logger.WithField("path to save pdf", fullpath).WithError(err).Error("error while writing file")
		return
	}
	f, err := os.Open(htmlFile)
	if err != nil {
		pdf.Logger.WithField("path to save pdf", fullpath).WithError(err).Error("error while opening file")
		return
	}
	if f != nil {
		defer func() {
			if closeInstanceErr := f.Close(); closeInstanceErr != nil {
				pdf.Logger.WithField("path to save pdf", fullpath).WithError(closeInstanceErr).Error("error while closing file instance in defer")
			}
		}()
	}
	if pdf.WKHTMLTOPDF_PATH != "" {
		HTMLToPDF.SetPath(pdf.WKHTMLTOPDF_PATH)
	}
	pdfGenerator, err := HTMLToPDF.NewPDFGenerator()
	if err != nil {
		pdf.Logger.WithField("path to save pdf", fullpath).WithError(err).Error("error while calling wkhtmltopdf func")
		return
	}
	pdfGenerator.AddPage(HTMLToPDF.NewPageReader(f))
	pdfGenerator.PageSize.Set(HTMLToPDF.PageSizeA4)
	pdfGenerator.Dpi.Set(300)

	if err = pdfGenerator.Create(); err != nil {
		pdf.Logger.WithField("path to save pdf", fullpath).WithError(err).Error("error while calling wkhtmltopdf func")
		return
	}
	if err = pdfGenerator.WriteFile(fullpath); err != nil {
		pdf.Logger.WithField("path to save pdf", fullpath).WithError(err).Error("error while generating pdf")
		return
	}
	if pdf.DeleteFromServer {
		return os.Remove(fullpath)
	}
	return nil
}
