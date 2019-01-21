package hikvision

import (
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	port, _ := strconv.ParseInt(config.Address.Port(), 10, 64)
	password, _ := config.Address.User.Password()

	device := &Bind{
		isapi: hikvision.NewISAPI(config.Address.Hostname(), port, config.Address.User.Username(), password),
		alertStreamingHistory: make(map[string]time.Time),
		address:               config.Address.URL,
		livenessInterval:      config.LivenessInterval,
		livenessTimeout:       config.LivenessTimeout,
		updaterInterval:       config.UpdaterInterval,
		updaterTimeout:        config.UpdaterTimeout,
		ptzInterval:           config.PTZInterval,
		ptzTimeout:            config.PTZTimeout,
		ptzEnabled:            config.PTZEnabled,
		eventsEnabled:         config.EventsEnabled,
		eventsIgnoreInterval:  config.EventsIgnoreInterval,
	}

	return device, nil
}
