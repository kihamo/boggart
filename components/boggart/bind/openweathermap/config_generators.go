package openweathermap

import (
	"strconv"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() []generators.Step {
	opts, err := b.MQTT().ClientOptions()
	if err != nil {
		return nil
	}

	meta := b.Meta()
	id := meta.ID()
	filePrefix := openhab.FilePrefixFromBindMeta(meta)
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	broker := openhab.BrokerFromClientOptionsReader(opts)

	steps := []generators.Step{
		{
			FilePath: openhab.DirectoryThings + "broker.things",
			Content:  broker.String(),
		},
	}

	thing := openhab.GenericThingFromBindMeta(meta).
		WithBroker(broker).
		AddChannels(
			openhab.BindStatusChannel(meta),
			openhab.NewChannel("CurrentTemp", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicCurrentTemp.Format(id)).
				AddItems(
					openhab.NewItem(itemPrefix+"CurrentTemp", openhab.ItemTypeNumber).
						WithLabel("Current temperature").
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

		thing.AddChannels(
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

	if content := thing.String(); content != "" {
		steps = append(steps, generators.Step{
			FilePath: openhab.DirectoryThings + filePrefix + ".things",
			Content:  content,
		})
	}

	if content := thing.Items().String(); content != "" {
		steps = append(steps, generators.Step{
			FilePath: openhab.DirectoryItems + filePrefix + ".items",
			Content:  content,
		})
	}

	return steps
}
