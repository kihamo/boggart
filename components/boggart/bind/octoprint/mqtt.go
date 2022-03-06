package octoprint

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/octoprint/client/system"
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

func (b *Bind) callbackMQTTExecuteCommand(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	parts := message.Topic().Split()

	if len(parts) < 2 {
		return errors.New("bad topic name")
	}

	params := system.NewExecuteCommandParamsWithContext(ctx).
		WithAction(parts[len(parts)-1]).
		WithSource(parts[len(parts)-2])

	_, err := b.provider.System.ExecuteCommand(params, nil)
	return err
}
