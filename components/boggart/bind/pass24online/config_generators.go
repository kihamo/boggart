package pass24online

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	const idFeedEvent = "FeedEvent"

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idFeedEvent, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicFeedEvent).
			AddItems(
				openhab.NewItem(itemPrefix+idFeedEvent, openhab.ItemTypeString).
					WithLabel("Feed event").
					WithIcon("text"),
			),
	)
}
