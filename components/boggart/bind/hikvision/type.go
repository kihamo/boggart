package hikvision

import (
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	port, _ := strconv.ParseInt(config.Address.Port(), 10, 64)
	password, _ := config.Address.User.Password()

	bind := &Bind{
		isapi:                 hikvision.NewISAPI(config.Address.Hostname(), port, config.Address.User.Username(), password),
		alertStreamingHistory: make(map[string]time.Time),
		address:               config.Address.URL,
		config:                config,
	}

	return bind, nil
}
