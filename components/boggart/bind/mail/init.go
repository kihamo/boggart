package mail

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/mail/smtp"
)

func init() {
	boggart.RegisterBindType("mail:smtp", smtp.Type{})
}
