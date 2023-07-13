package logx

import (
	"crypto/tls"
	"io"

	"go.uber.org/zap"
	gomail "gopkg.in/gomail.v2"
)

// Mail 发送邮件
type Mail struct {
	Level    *Level
	From     string
	To       string
	Subject  string
	Stmp     string
	Port     int
	Password string
	fields   []zap.Field
}

// Write 实现io.Writer，发送邮件
func (lm *Mail) Write(p []byte) (int, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", lm.From)
	m.SetHeader("To", lm.To)
	m.SetHeader("Subject", lm.Subject)
	m.SetBody("text/plain", string(p))
	d := gomail.NewDialer(lm.Stmp, lm.Port, lm.From, lm.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	go func(m *gomail.Message, d *gomail.Dialer) {
		if err := d.DialAndSend(m); err != nil {
			return
		}
	}(m, d)
	return len(p), nil
}

func getMailWriter(level *Level, port int, from, to, subject, stmp, password string, fields ...zap.Field) (io.Writer, error) {
	lm := &Mail{
		Level:    level,
		From:     from,
		To:       to,
		Subject:  subject,
		Stmp:     stmp,
		Port:     port,
		Password: password,
		fields:   fields,
	}

	return lm, nil
}
