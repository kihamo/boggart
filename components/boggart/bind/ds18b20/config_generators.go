package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	sensors, err := b.Sensors()
	if err != nil {
		return nil, err
	}

	channels := make([]*openhab.Channel, 0, len(sensors))

	for _, sensor := range sensors {
		id := openhab.IDNormalizeCamelCase("Sensor " + sensor)

		channels = append(channels,
			openhab.NewChannel(id, openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicValue.Format(sensor)).
				AddItems(
					openhab.NewItem(itemPrefix+id, openhab.ItemTypeNumber).
						WithLabel("Temperature [%.2f °C]").
						WithIcon("temperature"),
				),
		)
	}

	return openhab.StepsByBind(b, nil, channels...)
}
