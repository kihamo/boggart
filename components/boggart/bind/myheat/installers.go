package myheat

import (
	"context"
	"errors"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
	"github.com/kihamo/boggart/providers/myheat/device/client/sensors"
	"github.com/kihamo/boggart/providers/myheat/device/client/state"
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

	stateObjResponse, err := b.client.State.GetObjState(state.NewGetObjStateParamsWithContext(ctx), nil)
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()
	channels := make([]*openhab.Channel, 0, len(sensorsResponse.GetPayload())+1+2)

	const (
		idSensor               = "Sensor"
		idSensorSecurityArmed  = "SecurityArmed"
		idSensorGSMSignalLevel = "GSMSignalLevel"
		idSensorGSMBalance     = "GSMBalance"
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

	if v := stateObjResponse.Payload.SecurityArmed; v != nil {
		channels = append(channels,
			openhab.NewChannel(idSensorSecurityArmed, openhab.ChannelTypeContact).
				WithStateTopic(cfg.TopicSecurityArmedState.Format(sn)).
				WithOn("true").
				WithOff("false").
				AddItems(
					openhab.NewItem(itemPrefix+idSensorSecurityArmed, openhab.ItemTypeContact).
						WithLabel("Security armed").
						WithIcon("shield"),
				),
		)
	}

	channels = append(channels,
		openhab.NewChannel(idSensorGSMSignalLevel, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicGSMSignalLevel.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idSensorGSMSignalLevel, openhab.ItemTypeNumber).
					WithLabel("GSM signal level").
					WithIcon("qualityofservice"),
			),
		openhab.NewChannel(idSensorGSMBalance, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicGSMBalance.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idSensorGSMBalance, openhab.ItemTypeNumber).
					WithLabel("GSM balance [%.2f ₽]").
					WithIcon("price"),
			),
	)

	return openhab.StepsByBind(b, nil, channels...)
}
