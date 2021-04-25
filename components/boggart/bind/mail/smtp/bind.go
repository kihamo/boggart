package smtp

import (
	"bytes"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/kihamo/boggart/components/boggart/di"
)

var (
	crlf = []byte("\r\n")
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind

	auth smtp.Auth
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	var (
		username string
		password string
	)

	cfg := b.config()

	if user := cfg.DSN.User; user != nil {
		username = user.Username()
		password, _ = user.Password()
	}

	b.auth = smtp.PlainAuth("", username, password, cfg.DSN.Hostname())

	return nil
}

func (b *Bind) Send(to []string, subject string, body []byte) error {
	cfg := b.config()
	payload := bytes.NewBuffer(body)
	mimeType := http.DetectContentType(body)

	message := bytes.NewBuffer(nil)

	message.WriteString("From: ")
	message.WriteString(cfg.Sender)
	message.Write(crlf)

	message.WriteString("To: ")
	message.WriteString(strings.Join(to, ","))
	message.Write(crlf)

	message.WriteString("Subject: ")
	message.WriteString(subject)
	message.Write(crlf)

	if mimeType != "" && mimeType != "application/octet-stream" {
		message.WriteString("MIME-version: 1.0")
		message.Write(crlf)
		message.WriteString("Content-Type: ")
		message.WriteString(mimeType)
		message.Write(crlf)
	}

	message.Write(crlf)
	message.ReadFrom(payload)

	return smtp.SendMail(cfg.DSN.Host, b.auth, cfg.Sender, to, message.Bytes())
}
