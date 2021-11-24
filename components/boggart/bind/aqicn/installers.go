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
	channels := make([]*openhab.Channel, 0, 11)

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
		openhab.NewChannel("CurrentPm25Value", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentPm25Value.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentPm25Value", openhab.ItemTypeNumber).
					WithLabel("Current PM25 level [%.1f µg/m³]").
					WithIcon("line"),
			),
		openhab.NewChannel("CurrentPm10Value", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentPm10Value.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentPm10Value", openhab.ItemTypeNumber).
					WithLabel("Current PM10 level [%.1f µg/m³]").
					WithIcon("line"),
			),
		openhab.NewChannel("CurrentO3Value", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentO3Value.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentO3Value", openhab.ItemTypeNumber).
					WithLabel("Current O3 level [%.1f ppm]").
					WithIcon("line"),
			),
		openhab.NewChannel("CurrentNO2Value", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentNO2Value.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentNO2Value", openhab.ItemTypeNumber).
					WithLabel("Current NO2 level [%.1f ppm]").
					WithIcon("line"),
			),
		openhab.NewChannel("CurrentCOValue", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentCOValue.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentCOValue", openhab.ItemTypeNumber).
					WithLabel("Current CO level [%.1f ppm]").
					WithIcon("line"),
			),
		openhab.NewChannel("CurrentSO2Value", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentSO2Value.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentSO2Value", openhab.ItemTypeNumber).
					WithLabel("Current SO2 level [%.1f ppm]").
					WithIcon("line"),
			),
	)

	return openhab.StepsByBind(b, nil, channels...)
}
