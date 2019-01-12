package hikvision

import (
	"net/url"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Config)

	u, _ := url.Parse(config.Address)
	port, _ := strconv.ParseInt(u.Port(), 10, 64)
	password, _ := u.User.Password()

	device := &Bind{
		isapi:                 hikvision.NewISAPI(u.Hostname(), port, u.User.Username(), password),
		alertStreamingHistory: make(map[string]time.Time),
	}

	device.Init()

	return device, nil
}
