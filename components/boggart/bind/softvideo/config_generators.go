package softvideo

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() []generators.Step {
	opts, err := b.MQTT().ClientOptions()
	if err != nil {
		return nil
	}

	meta := b.Meta()
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
			openhab.BindSerialNumberChannel(meta),
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
