package modbus

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemDevice,
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(_ context.Context, system installer.System) ([]installer.Step, error) {
	if system == installer.SystemDevice {
		return []installer.Step{{
			Description: "Activate Modbus over TCP",
			Content: `1. Click Setup in menu
2. Select Screen Lock item
3. Enter pin 3036 and click lock icon
4. Wifi Modbus TCP
5. Click Next button
6. Click Exit button and device reload`,
		}}, nil
	}

	meta := b.Meta()
	id := meta.ID()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()

	const (
		idDeviceType          = "DeviceType"
		idHeatingOutputStatus = "HeatingOutputStatus"
		idRoomTemperature     = "RoomTemperature"
		idFloorTemperature    = "FloorTemperature"
		idHumidity            = "Humidity"
		idSetTemperature      = "SetTemperature"
	)

	channels := []*openhab.Channel{
		openhab.NewChannel(idDeviceType, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicDeviceType.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idDeviceType, openhab.ItemTypeNumber).
					WithLabel("Device type").
					WithIcon("text"),
			),
		openhab.NewChannel(idHeatingOutputStatus, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicHeatingOutputStatus.Format(id)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idHeatingOutputStatus, openhab.ItemTypeContact).
					WithLabel("Heating output [%s]").
					WithIcon("fire"),
			),
		openhab.NewChannel(idRoomTemperature, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicRoomTemperature.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idRoomTemperature, openhab.ItemTypeNumber).
					WithLabel("Room temperature [%.2f °C]").
					WithIcon("temperature"),
			),
		openhab.NewChannel(idFloorTemperature, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicFloorTemperature.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idFloorTemperature, openhab.ItemTypeNumber).
					WithLabel("Floor temperature [%.2f °C]").
					WithIcon("temperature"),
			),
		openhab.NewChannel(idHumidity, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicHumidity.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idHumidity, openhab.ItemTypeNumber).
					WithLabel("Humidity [%.1f %%]").
					WithIcon("humidity"),
			),
		openhab.NewChannel(idSetTemperature, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicSetTemperatureState.Format(id)).
			WithCommandTopic(cfg.TopicSetTemperature.Format(id)).
			WithMin(5).
			WithMax(35).
			WithStep(1).
			AddItems(
				openhab.NewItem(itemPrefix+idSetTemperature, openhab.ItemTypeDimmer).
					WithLabel("Set temperature [%.2f °C]").
					WithIcon("temperature"),
			),
	}

	return openhab.StepsByBind(b, nil, channels...)
}
