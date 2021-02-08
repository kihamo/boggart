package sp3s

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

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanWatts),
	},
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
					WithLabel("Power [JS(human_watts.js):%s]").
					WithIcon("energy"),
			),
	)
}
