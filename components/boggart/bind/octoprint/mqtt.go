package octoprint

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) callbackMQTTTemperature(message mqtt.Message, offset int) error {
	temperature := NewTemperature(b.PluginMQTTSettings().TimestampFieldname)

	if err := message.JSONUnmarshal(temperature); err != nil {
		return err
	}

	parts := message.Topic().Split()
	id := b.Meta().ID()

	metricDeviceTemperatureActual.With("id", id).With("device", parts[offset]).Set(temperature.Actual)
	metricDeviceTemperatureTarget.With("id", id).With("device", parts[offset]).Set(temperature.Target)

	b.devicesMutex.Lock()
	b.devices[parts[offset]] = true
	b.devicesMutex.Unlock()

	return nil
}

//func (b *Bind) callbackMQTTJob(message mqtt.Message) error {
//	fmt.Println(message)
//
//	return nil
//}
