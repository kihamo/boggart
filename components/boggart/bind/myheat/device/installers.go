package device

import (
	"context"
	"errors"
	"strconv"

	"github.com/kihamo/boggart/providers/myheat"

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

	sensorsResponse, err := b.client.Sensors.GetSensors(sensors.NewGetSensorsParamsWithContext(ctx))
	if err != nil {
		return nil, err
	}

	stateObjResponse, err := b.client.State.GetObjState(state.NewGetObjStateParamsWithContext(ctx))
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()
	channels := make([]*openhab.Channel, 0, len(sensorsResponse.Payload)+1+7)

	const (
		idSensor                         = "Sensor"
		idSecurityArmed                  = "SecurityArmed"
		idDeviceSeverity                 = "DeviceSeverity"
		idInternetConnected              = "InternetConnected"
		idGSMSignalLevel                 = "GSMSignalLevel"
		idGSMBalance                     = "GSMBalance"
		idAlarmPowerSupply               = "AlarmPowerSupply"
		idAlarmReplaceBattery            = "AlarmReplaceBattery"
		idAlarmGSMBalance                = "AlarmGSMBalance"
		idEnvironmentStateTemperature    = "EnvironmentStateTemperature"
		idHeaterHeatingFeedTemperature   = "HeaterHeatingFeedTemperature"
		idHeaterHeatingReturnTemperature = "HeaterHeatingReturnTemperature"
		idHeaterHeatingCircuitPressure   = "HeaterHeatingCircuitPressure"
	)

	for _, sensor := range sensorsResponse.Payload {
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

	for _, heater := range stateObjResponse.Payload.Heaters {
		id := strconv.FormatInt(heater.ID, 10)

		if _, ok := heater.State[myheat.HeaterHeatingFeedTemperatureCelsius]; ok {
			channels = append(channels,
				openhab.NewChannel(idHeaterHeatingFeedTemperature+id, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicHeaterHeatingFeedTemperature.Format(sn, heater.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idHeaterHeatingFeedTemperature+id, openhab.ItemTypeNumber).
							WithLabel(heater.Name+" heating feed [%.2f °C]").
							WithIcon("temperature"),
					),
			)
		}

		if _, ok := heater.State[myheat.HeaterHeatingReturnTemperatureCelsius]; ok {
			channels = append(channels,
				openhab.NewChannel(idHeaterHeatingReturnTemperature+id, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicHeaterHeatingReturnTemperature.Format(sn, heater.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idHeaterHeatingReturnTemperature+id, openhab.ItemTypeNumber).
							WithLabel(heater.Name+" heating return [%.2f °C]").
							WithIcon("temperature"),
					),
			)
		}

		if _, ok := heater.State[myheat.HeaterHeatingCircuitPressureBar]; ok {
			channels = append(channels,
				openhab.NewChannel(idHeaterHeatingCircuitPressure+id, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicHeaterHeatingCircuitPressure.Format(sn, heater.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idHeaterHeatingCircuitPressure+id, openhab.ItemTypeNumber).
							WithLabel(heater.Name+" heating circuit pressure [%.2f bar]").
							WithIcon("pressure"),
					),
			)
		}
	}

	for _, env := range stateObjResponse.Payload.Environments {
		id := strconv.FormatInt(env.ID, 10)

		if _, ok := env.State[myheat.EnvironmentStateTemperatureValue]; ok {
			channels = append(channels,
				openhab.NewChannel(idEnvironmentStateTemperature+id, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicEnvironmentStateTemperature.Format(sn, env.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idEnvironmentStateTemperature+id, openhab.ItemTypeNumber).
							WithLabel(env.Name+" state [%.2f °C]").
							WithIcon("temperature"),
					),
			)
		}
	}

	if v := stateObjResponse.Payload.SecurityArmed; v != nil {
		channels = append(channels,
			openhab.NewChannel(idSecurityArmed, openhab.ChannelTypeSwitch).
				WithStateTopic(cfg.TopicSecurityArmedState.Format(sn)).
				WithCommandTopic(cfg.TopicSecurityArmed.Format(sn)).
				WithOn("true").
				WithOff("false").
				AddItems(
					openhab.NewItem(itemPrefix+idSecurityArmed, openhab.ItemTypeSwitch).
						WithLabel("Security armed").
						WithIcon("shield"),
				),
		)
	}

	channels = append(channels,
		openhab.NewChannel(idDeviceSeverity, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicDeviceSeverity.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idDeviceSeverity, openhab.ItemTypeNumber).
					WithLabel("Device severity").
					WithIcon("error"),
			),
		openhab.NewChannel(idInternetConnected, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicInternetConnected.Format(sn)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idInternetConnected, openhab.ItemTypeSwitch).
					WithLabel("Internet connected [%s]").
					WithIcon("network"),
			),
		openhab.NewChannel(idGSMSignalLevel, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicGSMSignalLevel.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idGSMSignalLevel, openhab.ItemTypeNumber).
					WithLabel("GSM signal level").
					WithIcon("qualityofservice"),
			),
		openhab.NewChannel(idGSMBalance, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicGSMBalance.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idGSMBalance, openhab.ItemTypeNumber).
					WithLabel("GSM balance [%.2f ₽]").
					WithIcon("price"),
			),
		openhab.NewChannel(idAlarmPowerSupply, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicAlarmPowerSupply.Format(sn)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idAlarmPowerSupply, openhab.ItemTypeSwitch).
					WithLabel("Alarm power supply [%s]").
					WithIcon("siren"),
			),
		openhab.NewChannel(idAlarmReplaceBattery, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicAlarmReplaceBattery.Format(sn)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idAlarmReplaceBattery, openhab.ItemTypeSwitch).
					WithLabel("Alarm replace battery [%s]").
					WithIcon("siren"),
			),
		openhab.NewChannel(idAlarmGSMBalance, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicAlarmGSMBalance.Format(sn)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idAlarmGSMBalance, openhab.ItemTypeSwitch).
					WithLabel("Alarm GSM balance [%s]").
					WithIcon("siren"),
			),
	)

	return openhab.StepsByBind(b, nil, channels...)
}
