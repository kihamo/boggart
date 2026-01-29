package modbus

import (
	"context"
	"fmt"
	"net"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemDevice,
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(ctx context.Context, system installer.System) ([]installer.Step, error) {
	cfg := b.config()

	if system == installer.SystemDevice {
		var exampleIP string
		if val := net.ParseIP(cfg.DSN.Hostname()); val != nil {
			exampleIP = " (current is " + val.String() + ")"
		}

		return []installer.Step{{
			Description: "Activate Modbus over TCP",
			Content: `1. Click Setup in menu
2. Select Screen Lock item
3. Enter pin 3036 and click lock icon
4. Setn Wifi Modbus TCP (activated read only registers) and Modbus RTU 485 (activated read and write registers)
5. Click Next button
6. Click Exit button and device reload`,
		}, {
			Description: "Set IP address (DHCP not supported)",
			Content: `1. Click Setup in menu
2. Select Network settings item
3. Select 01. Network settings item
4. Change Gateway, Net Mask and IP` + exampleIP + ` fields
5. Exit to main menu and reload device`,
		}}, nil
	}

	deviceType, err := b.DeviceType(ctx)
	if err != nil {
		return nil, fmt.Errorf("get device type failed: %w", err)
	}

	meta := b.Meta()
	id := meta.ID()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	const (
		temperatureUnit = "°C"

		idDeviceType          = "DeviceType"
		idRoomTemperature     = "RoomTemperature"
		idFloorTemperature    = "FloorTemperature"
		idHumidity            = "Humidity"
		idHeatingValve        = "HeatingValve"
		idCoolingValve        = "CoolingValve"
		idStatus              = "Power"
		idHeatingOutputStatus = "HeatingOutput"
		idHoldingFunction     = "HoldingFunction"
		idFloorOverheat       = "FloorOverheat"
		idFanSpeedNumbers     = "FanSpeedNumbers"
		idSystemMode          = "SystemMode"
		idFanSpeed            = "FanSpeed"
		idTargetTemperature   = "TargetTemperature"
		idAway                = "Away"
		idAwayTemperature     = "AwayTemperature"
		idHoldingTime         = "HoldingTime"
		idHoldingTemperature  = "HoldingTemperature"
	)

	channels := []*openhab.Channel{
		openhab.NewChannel(idDeviceType, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicDeviceType.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idDeviceType, openhab.ItemTypeString).
					WithLabel("Device type").
					WithIcon("text"),
			),
	}

	if deviceType.IsSupportedRoomTemperature() {
		channels = append(channels, openhab.NewChannel(idRoomTemperature, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicRoomTemperature.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idRoomTemperature, openhab.ItemTypeNumber).
					WithLabel("Room temperature [%.1f "+temperatureUnit+"]").
					WithIcon("temperature"),
			))
	}

	if deviceType.IsSupportedFloorTemperature() {
		channels = append(channels, openhab.NewChannel(idFloorTemperature, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicFloorTemperature.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idFloorTemperature, openhab.ItemTypeNumber).
					WithLabel("Floor temperature [%.1f "+temperatureUnit+"]").
					WithIcon("temperature"),
			))
	}

	if deviceType.IsSupportedHumidity() {
		channels = append(channels, openhab.NewChannel(idHumidity, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicHumidity.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idHumidity, openhab.ItemTypeNumber).
					WithLabel("Humidity [%.1f %%]").
					WithIcon("humidity"),
			))
	}

	if deviceType.IsSupportedHeatingValve() {
		channels = append(channels, openhab.NewChannel(idHeatingValve, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicHeatingValve.Format(id)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idHeatingValve, openhab.ItemTypeContact).
					WithLabel("Heating valve [%s]").
					WithIcon("fire"),
			))
	}

	if deviceType.IsSupportedCoolingValve() {
		channels = append(channels, openhab.NewChannel(idCoolingValve, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicCoolingValve.Format(id)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idCoolingValve, openhab.ItemTypeContact).
					WithLabel("Cooling valve [%s]").
					WithIcon("fire"),
			))
	}

	if deviceType.IsSupportedHeatingOutput() {
		channels = append(channels, openhab.NewChannel(idHeatingOutputStatus, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicHeatingOutputStatus.Format(id)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idHeatingOutputStatus, openhab.ItemTypeContact).
					WithLabel("Heating output [%s]").
					WithIcon("fire"),
			))
	}

	if deviceType.IsSupportedHoldingFunction() {
		channels = append(channels, openhab.NewChannel(idHoldingFunction, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicHoldingFunction.Format(id)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idHoldingFunction, openhab.ItemTypeContact).
					WithLabel("Holding function [%s]").
					WithIcon("fire"),
			))
	}

	if deviceType.IsSupportedFloorOverheat() {
		channels = append(channels, openhab.NewChannel(idFloorOverheat, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicFloorOverheat.Format(id)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idFloorOverheat, openhab.ItemTypeContact).
					WithLabel("Floor overheat [%s]").
					WithIcon("siren"),
			))
	}

	if deviceType.IsSupportedFanSpeedNumbers() {
		channels = append(channels, openhab.NewChannel(idFanSpeedNumbers, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicFanSpeedNumbers.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idFanSpeedNumbers, openhab.ItemTypeNumber).
					WithLabel("Fan speed numbers [%d]").
					WithIcon("fan"),
			))
	}

	if deviceType.IsSupportedTargetTemperature() {
		min := 5.0
		max := 35.0

		if val, e := b.Provider().TargetTemperatureMinimum(); e == nil {
			min = float64(val)
		}

		if val, e := b.Provider().TargetTemperatureMaximum(); e == nil {
			max = float64(val)
		}

		channels = append(channels,
			openhab.NewChannel(idTargetTemperature, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicTargetTemperatureState.Format(id)).
				WithCommandTopic(cfg.TopicTargetTemperature.Format(id)).
				WithMin(min).
				WithMax(max).
				WithStep(0.5).
				AddItems(
					openhab.NewItem(itemPrefix+idTargetTemperature, openhab.ItemTypeNumber).
						WithLabel("Set temperature [%.1f "+temperatureUnit+"]").
						WithIcon("temperature"),
				))
	}

	if deviceType.IsSupportedStatus() {
		channels = append(channels,
			openhab.NewChannel(idStatus, openhab.ChannelTypeSwitch).
				WithStateTopic(cfg.TopicStatusState.Format(id)).
				WithCommandTopic(cfg.TopicStatus.Format(id)).
				WithOn("true").
				WithOff("false").
				AddItems(
					openhab.NewItem(itemPrefix+idStatus, openhab.ItemTypeSwitch).
						WithLabel("Status []"),
				))
	}

	if deviceType.IsSupportedSystemMode() {
		channels = append(channels,
			openhab.NewChannel(idSystemMode, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicSystemModeState.Format(id)).
				WithCommandTopic(cfg.TopicSystemMode.Format(id)).
				WithMin(0).
				WithMax(4).
				AddItems(
					openhab.NewItem(itemPrefix+idSystemMode, openhab.ItemTypeNumber).
						WithLabel("System mode [%s]"),
				))
	}

	if deviceType.IsSupportedFanSpeed() {
		channels = append(channels,
			openhab.NewChannel(idFanSpeed, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicFanSpeedState.Format(id)).
				WithCommandTopic(cfg.TopicFanSpeed.Format(id)).
				WithMin(0).
				WithMax(3).
				AddItems(
					openhab.NewItem(itemPrefix+idFanSpeed, openhab.ItemTypeNumber).
						WithLabel("Fan speed [%s]").
						WithIcon("fan"),
				))
	}

	if deviceType.IsSupportedAway() {
		channels = append(channels,
			openhab.NewChannel(idAway, openhab.ChannelTypeSwitch).
				WithStateTopic(cfg.TopicAwayState.Format(id)).
				WithCommandTopic(cfg.TopicAway.Format(id)).
				WithOn("true").
				WithOff("false").
				AddItems(
					openhab.NewItem(itemPrefix+idAway, openhab.ItemTypeSwitch).
						WithLabel("Away []"),
				))
	}

	if deviceType.IsSupportedAwayTemperature() {
		channels = append(channels,
			openhab.NewChannel(idAwayTemperature, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicAwayTemperatureState.Format(id)).
				WithCommandTopic(cfg.TopicAwayTemperature.Format(id)).
				WithMin(7).
				WithMax(35).
				WithStep(1).
				AddItems(
					openhab.NewItem(itemPrefix+idAwayTemperature, openhab.ItemTypeNumber).
						WithLabel("Away temperature [%.1f "+temperatureUnit+"]").
						WithIcon("temperature"),
				))
	}

	if deviceType.IsSupportedHoldingTemperatureAndTime() || deviceType.IsSupportedHoldingTemperatureAndTime() {
		channels = append(channels,
			openhab.NewChannel(idHoldingTime, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicHoldingTimeState.Format(id)).
				WithCommandTopic(cfg.TopicHoldingTime.Format(id)).
				WithMax(1439).
				WithStep(1).
				AddItems(
					openhab.NewItem(itemPrefix+idHoldingTime, openhab.ItemTypeNumber).
						WithLabel("Holding time [%d minutes]").
						WithIcon("time"),
				))
	}

	if deviceType.IsSupportedHoldingTemperatureAndTime() || deviceType.IsSupportedHoldingTemperature() {
		channels = append(channels,
			openhab.NewChannel(idHoldingTemperature, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicHoldingTemperatureState.Format(id)).
				WithCommandTopic(cfg.TopicHoldingTemperature.Format(id)).
				WithMin(5).
				WithMax(40).
				WithStep(1).
				AddItems(
					openhab.NewItem(itemPrefix+idHoldingTemperature, openhab.ItemTypeNumber).
						WithLabel("Holding temperature [%.1f "+temperatureUnit+"]").
						WithIcon("temperature"),
				))
	}

	return openhab.StepsByBind(b, nil, channels...)
}
