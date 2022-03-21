package email

import (
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
)

type Attachment struct {
	Filename    string
	ContentType string
	Header      textproto.MIMEHeader
	Content     []byte
	HTMLRelated bool
}

type EmailConfigs struct {
	Host, Port, From, Password string
	To, CC, BCC                []string
	Subject, BodyText, HTML    string
	Attachments                []Attachment
	describer                  EmailDescriber
}

func (conf EmailConfigs) mapAttachments() []*email.Attachment {
	var attachments = []*email.Attachment{}
	for _, attach := range conf.Attachments {
		attachments = append(attachments, &email.Attachment{
			Filename:    attach.Filename,
			ContentType: attach.ContentType,
			Header:      attach.Header,
			Content:     attach.Content,
			HTMLRelated: attach.HTMLRelated,
		})
	}
	return attachments
}
func (config EmailConfigs) getDescriber() EmailDescriber {
	return (&email.Email{
		From:        config.From,
		To:          config.To,
		Cc:          config.CC,
		Bcc:         config.BCC,
		Subject:     config.Subject,
		Text:        []byte(config.BodyText),
		HTML:        []byte(config.HTML),
		Headers:     textproto.MIMEHeader{},
		Attachments: config.mapAttachments(),
	})
}
func (config EmailConfigs) Send() error {
	if config.describer == nil {
		config.describer = config.getDescriber()
	}
	return config.describer.Send(config.Host+":"+config.Port, smtp.PlainAuth("", config.From, config.Password, config.Host))
}
