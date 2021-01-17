package sp3s

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
			openhab.BindMACChannel(meta),
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
