package emailx

import (
	"fmt"
	"log"
	"net"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type Email struct {
	Host string //smtp.gmail.com"
	Port string //587
	Auth smtp.Auth
}

type Mail struct {
	From    string
	To      []string
	Subject string
	Text    []byte
	HTML    []byte
}

// xxx@gmail.com
// xxxx
func NewEmail(host, port, username, password string) *Email {
	return &Email{
		Port: port,
		Host: host,
		Auth: smtp.PlainAuth("", username, password, host),
	}
}

func (t *Email) SendEmail(mail Mail) (_err error) {
	e := &email.Email{
		From:    fmt.Sprintf("发送者名字<%s>", mail.From),
		To:      mail.To,
		Subject: mail.Subject,
		Text:    mail.Text,
		HTML:    mail.HTML,
	}
	_err = e.Send(net.JoinHostPort(t.Host, t.Port), t.Auth)
	if _err != nil {
		log.Fatal(_err)
	}
	return
}
