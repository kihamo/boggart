package premiergc

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

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel("Balance", openhab.ChannelTypeNumber).
			WithStateTopic(b.config().TopicBalance).
			AddItems(
				openhab.NewItem(itemPrefix+"Balance", openhab.ItemTypeNumber).
					WithLabel("Balance [%.2f ₽]").
					WithIcon("price"),
			),
	)
}
