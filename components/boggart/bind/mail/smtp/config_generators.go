package smtp

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	meta := b.Meta()
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	const idSend = "SendJSON"

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idSend, openhab.ChannelTypeString).
			WithCommandTopic(b.config.TopicSend.Format(meta.ID())).
			AddItems(
				openhab.NewItem(itemPrefix+idSend, openhab.ItemTypeString),
			),
	)
}
