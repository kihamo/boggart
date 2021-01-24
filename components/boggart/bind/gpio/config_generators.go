package gpio

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	meta := b.Meta()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	const id = "GPIO"

	var channel *openhab.Channel

	if b.Mode() != ModeOut {
		channel = openhab.NewChannel(id, openhab.ChannelTypeContact).
			AddItems(
				openhab.NewItem(itemPrefix+id, openhab.ItemTypeContact).
					WithLabel(meta.Description() + "[]"),
			)
	} else {
		channel = openhab.NewChannel(id, openhab.ChannelTypeSwitch).
			WithCommandTopic(b.config.TopicPinSet).
			AddItems(
				openhab.NewItem(itemPrefix+id, openhab.ItemTypeSwitch).
					WithLabel(meta.Description() + "[]"),
			)
	}

	channel.WithStateTopic(b.config.TopicPinState).
		WithLabel(meta.Description()).
		WithOn("true").
		WithOff("false")

	return openhab.StepsByBind(b, nil, channel)
}
