package aqicn

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	meta := b.Meta()
	id := meta.ID()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()
	channels := make([]*openhab.Channel, 0, 5)

	channels = append(channels,
		openhab.NewChannel("CurrentTemperature", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentTemperature.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentTemperature", openhab.ItemTypeNumber).
					WithLabel("Current temperature [%.2f °C]").
					WithIcon("temperature"),
			),
		openhab.NewChannel("CurrentPressure", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentPressure.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentPressure", openhab.ItemTypeNumber).
					WithLabel("Current pressure [%.1f hPa]").
					WithIcon("pressure"),
			),
		openhab.NewChannel("CurrentHumidity", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentHumidity.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentHumidity", openhab.ItemTypeNumber).
					WithLabel("Current humidity [%.1f %%]").
					WithIcon("humidity"),
			),
		openhab.NewChannel("CurrentDewPoint", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentDewPoint.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentDewPoint", openhab.ItemTypeNumber).
					WithLabel("Current dew point [%.1f °C]").
					WithIcon("temperature"),
			),
		openhab.NewChannel("CurrentWindSpeed", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentWindSpeed.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentWindSpeed", openhab.ItemTypeNumber).
					WithLabel("Current wind speed [%.1f m/s]").
					WithIcon("wind"),
			),
	)

	return openhab.StepsByBind(b, nil, channels...)
}
