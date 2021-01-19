package boggart

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	return openhab.StepsByBind(b, nil,
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
}
