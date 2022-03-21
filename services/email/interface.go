package email

import (
	"net/smtp"

	"github.com/pkg/errors"
)

type EmailDescriber interface {
	Send(addr string, a smtp.Auth) error
}

type mock struct {
	c byte
}

func (m mock) Send(addr string, a smtp.Auth) error {
	if m.c == 0 {
		return nil
	}
	return errors.New("mock email func error")
}
