package openweathermap

import (
	"strconv"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	meta := b.Meta()
	id := meta.ID()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	channels := make([]*openhab.Channel, 0, 1+7*6)

	channels = append(channels,
		openhab.NewChannel("CurrentTemp", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicCurrentTemp.Format(id)).
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
				WithStateTopic(b.config.TopicDailyTempMin.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"TempMin", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" min daily temperature [%.2f °C]").
						WithIcon("temperature"),
				),
			openhab.NewChannel(dayName+"TempMax", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicDailyTempMax.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"TempMax", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" max daily temperature [%.2f °C]").
						WithIcon("temperature"),
				),
			openhab.NewChannel(dayName+"TempDay", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicDailyTempDay.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"TempDay", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" day temperature [%.2f °C]").
						WithIcon("temperature"),
				),
			openhab.NewChannel(dayName+"TempNight", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicDailyTempNight.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"TempNight", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" night temperature [%.2f °C]").
						WithIcon("temperature"),
				),
			openhab.NewChannel(dayName+"TempMorn", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicDailyTempMorning.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"TempMorn", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" morning temperature [%.2f °C]").
						WithIcon("temperature"),
				),
			openhab.NewChannel(dayName+"WindSpeed", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicDailyWindSpeed.Format(id, i)).
				AddItems(
					openhab.NewItem(itemPrefix+dayName+"WindSpeed", openhab.ItemTypeNumber).
						WithLabel(dayLabel+" wind speed [%.2f m/s]").
						WithIcon("wind"),
				),
		)
	}

	return openhab.StepsByBind(b, nil, channels...)
}
