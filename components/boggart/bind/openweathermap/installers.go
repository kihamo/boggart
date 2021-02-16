package openweathermap

import (
	"context"
	"strconv"

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
	channels := make([]*openhab.Channel, 0, 1+7*6)

	channels = append(channels,
		openhab.NewChannel("CurrentTemp", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCurrentTemp.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CurrentTemp", openhab.ItemTypeNumber).
					WithLabel("Current temperature [%.2f °C]").
					WithIcon("temperature"),
			),
	)

	var (
		dayName  string
		dayLabel string
	)

	for i := 0; i <= 7; i++ {
		if i == 0 {
			dayName = "Today"
			dayLabel = "Today"
		} else {
			dayName = "Next" + strconv.Itoa(i)
			dayLabel = "Next " + strconv.Itoa(i)
		}

		channels = append(channels,
			openhab.NewChannel(dayName+"TempMin", openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicDailyTempMin.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"TempMin", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" min daily temperature [%.2f °C]").
						WithIcon("temperature"),
				),
			openhab.NewChannel(dayName+"TempMax", openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicDailyTempMax.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"TempMax", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" max daily temperature [%.2f °C]").
						WithIcon("temperature"),
				),
			openhab.NewChannel(dayName+"TempDay", openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicDailyTempDay.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"TempDay", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" day temperature [%.2f °C]").
						WithIcon("temperature"),
				),
			openhab.NewChannel(dayName+"TempNight", openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicDailyTempNight.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"TempNight", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" night temperature [%.2f °C]").
						WithIcon("temperature"),
				),
			openhab.NewChannel(dayName+"TempMorn", openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicDailyTempMorning.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"TempMorn", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" morning temperature [%.2f °C]").
						WithIcon("temperature"),
				),
			openhab.NewChannel(dayName+"WindSpeed", openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicDailyWindSpeed.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"WindSpeed", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" wind speed [%.2f m/s]").
						WithIcon("wind"),
				),
		)
	}

	return openhab.StepsByBind(b, nil, channels...)
}
