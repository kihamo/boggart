package internal

import (
	"github.com/kihamo/boggart/components/boggart/internal/subscribes"
	"github.com/kihamo/boggart/components/mqtt"
)

func (c *Component) initMQTT() {
	m := c.application.GetComponent(mqtt.ComponentName).(mqtt.Component)

	m.Subscribe(subscribes.NewOwnTracksSubscribe())
	m.Subscribe(subscribes.NewGPIOSubscribe(c.devicesManager))
}
