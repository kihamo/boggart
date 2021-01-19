package sp3s

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel("Status", openhab.ChannelTypeSwitch).
			WithStateTopic(b.config.TopicState).
			WithCommandTopic(b.config.TopicSet).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+"Status", openhab.ItemTypeSwitch).
					WithLabel("Status []"),
			),
		openhab.NewChannel("Power", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicPower).
			AddItems(
				openhab.NewItem(itemPrefix+"Power", openhab.ItemTypeNumber).
					WithLabel("Power [%.1f W]").
					WithIcon("energy"),
			),
	)
}
