package pdf

import "testing"

func TestNewPDFRequest_GenerateHTMLToPDF(t *testing.T) {

	tests := []struct {
		name    string
		request *NewPDFRequest
		wantErr bool
	}{
		{
			name: "Successfully without delete",
			request: &NewPDFRequest{
				TemplateHTML: `<html><body><h1 style="color:red;">This is an html
				from pdf to test color</h1></img></body></html>`,
				TemplateData: map[string]interface{}{
					"Name": "Testing",
				},
				FileName:         "Test.pdf",
				OutputPath:       "./../../",
				DeleteFromServer: false,
			},
			wantErr: false,
		},
		{
			name: "Successfully created",
			request: &NewPDFRequest{
				TemplateHTML: `<html><body><h1 style="color:red;">This is an html
				from pdf to test color</h1></img></body></html>`,
				TemplateData: map[string]interface{}{
					"Name": "Testing",
				},
				FileName:         "Test.pdf",
				OutputPath:       "./../../",
				DeleteFromServer: true,
			},
			wantErr: false,
		},
		{
			name: "Error",
			request: &NewPDFRequest{
				TemplateHTML: "<html>{{.Name}}</html>",
				TemplateData: map[string]interface{}{
					"Name": "Testing",
				},
				FileName:         "TestPDF",
				DeleteFromServer: true,
			},
			wantErr: true,
		},
		{
			name: "Error",
			request: &NewPDFRequest{
				TemplateHTML: "<html>{{.Name}}</html>",
				TemplateData: map[string]interface{}{
					"Name": "Testing",
				},
				FileName:         "TestPDF.pdf",
				DeleteFromServer: true,
			},
			wantErr: true,
		},
		{
			name: "Empty FileName",
			request: &NewPDFRequest{
				TemplateHTML: "<html><h1>{{.Name}}</h1></html>",
				TemplateData: map[string]interface{}{
					"Name": "Testing",
				},
				DeleteFromServer: true,
			},
			wantErr: true,
		},
		{
			name: "parsing error",
			request: &NewPDFRequest{
				TemplateHTML:     `</html></html></html>`,
				TemplateData:     map[string]interface{}{"Title": "My test title"},
				FileName:         "TestPDF.pdf",
				OutputPath:       "./../../",
				DeleteFromServer: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.request.GenerateHTMLToPDF(); (err != nil) != tt.wantErr {
				t.Errorf("[Test- case : %v]NewPDFRequest.GenerateHTMLToPDF() error : %v", tt.name, err)
			}
		})
	}
}
