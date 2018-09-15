package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/internal/subscribes"
	"github.com/kihamo/boggart/components/mqtt"
)

func (c *Component) initMQTT() {
	m := c.application.GetComponent(mqtt.ComponentName).(mqtt.Component)

	if c.config.Bool(boggart.ConfigOwnTracksEnabled) {
		m.Subscribe(subscribes.NewOwnTracksSubscribe())
	}
}
