package softvideo

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel("Balance", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicBalance).
			AddItems(
				openhab.NewItem(itemPrefix+"Balance", openhab.ItemTypeNumber).
					WithLabel("Balance [%.2f ₽]").
					WithIcon("price"),
			),
		openhab.NewChannel("Promise", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicPromise).
			AddItems(
				openhab.NewItem(itemPrefix+"Promise", openhab.ItemTypeNumber).
					WithLabel("Promise [%.2f ₽]").
					WithIcon("price"),
			),
	)
}
