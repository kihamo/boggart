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
	di.MetaBind
	di.MQTTBind

	config *Config
	auth   smtp.Auth
}

func (b *Bind) Send(to []string, subject string, body []byte) error {
	payload := bytes.NewBuffer(body)
	mimeType := http.DetectContentType(body)

	message := bytes.NewBuffer(nil)

	message.WriteString("From: ")
	message.WriteString(b.config.Sender)
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

	return smtp.SendMail(b.config.DSN.Host, b.auth, b.config.Sender, to, message.Bytes())
}
