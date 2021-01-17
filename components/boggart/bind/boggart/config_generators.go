package boggart

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
			openhab.NewChannel("Name", openhab.ChannelTypeString).
				WithStateTopic(b.config.TopicName).
				AddItems(
					openhab.NewItem(itemPrefix+"Name", openhab.ItemTypeString).
						WithLabel("Name"),
				),
			openhab.NewChannel("Version", openhab.ChannelTypeString).
				WithStateTopic(b.config.TopicBuild).
				AddItems(
					openhab.NewItem(itemPrefix+"Version", openhab.ItemTypeString).
						WithLabel("Version"),
				),
			openhab.NewChannel("Build", openhab.ChannelTypeString).
				WithStateTopic(b.config.TopicVersion).
				AddItems(
					openhab.NewItem(itemPrefix+"Build", openhab.ItemTypeString).
						WithLabel("Build"),
				),
			openhab.NewChannel("Shutdown", openhab.ChannelTypeSwitch).
				WithStateTopic(b.config.TopicShutdown).
				WithCommandTopic(b.config.TopicShutdown).
				WithOn("true").
				WithOff("false").
				AddItems(
					openhab.NewItem(itemPrefix+"Shutdown", openhab.ItemTypeSwitch).
						WithLabel("Shutdown []"),
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
