package ds18b20

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
