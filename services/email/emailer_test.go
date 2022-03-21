package email

import (
	"net/textproto"
	"reflect"
	"testing"

	"github.com/jordan-wright/email"
)

func TestEmailConfigs_Send(t *testing.T) {
	tests := []struct {
		name        string
		Host        string
		Port        string
		From        string
		Password    string
		To          []string
		CC          []string
		BCC         []string
		Subject     string
		BodyText    string
		HTML        string
		Attachments []Attachment
		Describer   EmailDescriber
		wantErr     bool
	}{
		{
			Describer: mock{},
		},
		{
			Describer: mock{c: 1},
			wantErr:   true,
		},
		{
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := EmailConfigs{
				Host:        tt.Host,
				Port:        tt.Port,
				From:        tt.From,
				Password:    tt.Password,
				To:          tt.To,
				CC:          tt.CC,
				BCC:         tt.BCC,
				Subject:     tt.Subject,
				BodyText:    tt.BodyText,
				HTML:        tt.HTML,
				Attachments: tt.Attachments,
				describer:   tt.Describer,
			}
			if err := config.Send(); (err != nil) != tt.wantErr {
				t.Errorf("EmailConfigs.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmailConfigs_getDescriber(t *testing.T) {
	var fields = EmailConfigs{
		Attachments: []Attachment{
			{
				Filename:    "",
				ContentType: "",
			},
		},
	}
	tests := []struct {
		name   string
		config EmailConfigs
		want   EmailDescriber
	}{
		{
			config: fields,
			want: &email.Email{
				From:        fields.From,
				To:          fields.To,
				Cc:          fields.CC,
				Bcc:         fields.BCC,
				Subject:     fields.Subject,
				Text:        []byte(fields.BodyText),
				HTML:        []byte(fields.HTML),
				Headers:     textproto.MIMEHeader{},
				Attachments: fields.mapAttachments(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.config
			if got := config.getDescriber(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailConfigs.getDescriber() = %v, want %v", got, tt.want)
			}
		})
	}
}
