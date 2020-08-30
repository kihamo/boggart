package smtp

import (
	"net/smtp"

	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	var (
		username string
		password string
	)

	if user := config.DSN.User; user != nil {
		username = user.Username()
		password, _ = user.Password()
	}

	bind := &Bind{
		config: config,
		auth:   smtp.PlainAuth("", username, password, config.DSN.Hostname()),
	}

	return bind, nil
}
