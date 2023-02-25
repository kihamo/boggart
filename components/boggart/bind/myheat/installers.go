package myheat

import (
	"context"
	"errors"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
	"github.com/kihamo/boggart/providers/myheat/device/client/sensors"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(ctx context.Context, _ installer.System) ([]installer.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	sensorsResponse, err := b.client.Sensors.GetSensors(sensors.NewGetSensorsParamsWithContext(ctx), nil)
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()
	channels := make([]*openhab.Channel, 0, len(sensorsResponse.GetPayload()))

	const (
		idSensor = "Sensor"
	)

	for _, sensor := range sensorsResponse.GetPayload() {
		id := strconv.FormatInt(sensor.ID, 10)

		switch sensor.Type {
		case 201: // Проводной датчик температуры
			channels = append(channels,
				openhab.NewChannel(idSensor+id, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicSensorValue.Format(sn, sensor.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idSensor+id, openhab.ItemTypeNumber).
							WithLabel(sensor.Name+" [%.2f °C]").
							WithIcon("temperature"),
					),
			)
		case 205: // Дискретный вход
			channels = append(channels,
				openhab.NewChannel(idSensor+id, openhab.ChannelTypeContact).
					WithStateTopic(cfg.TopicSensorValue.Format(sn, sensor.ID)).
					WithOn("1").
					WithOff("0").
					AddItems(
						openhab.NewItem(itemPrefix+idSensor+id, openhab.ItemTypeContact).
							WithLabel(sensor.Name),
					),
			)

		default:
			channels = append(channels,
				openhab.NewChannel(idSensor+id, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicSensorValue.Format(sn, sensor.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idSensor+id, openhab.ItemTypeNumber).
							WithLabel(sensor.Name).
							WithIcon("chart"),
					),
			)
		}
	}

	return openhab.StepsByBind(b, nil, channels...)
}
