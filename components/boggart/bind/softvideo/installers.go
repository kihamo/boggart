package softvideo

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
	cfg := b.config()

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel("Balance", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicBalance.Format(cfg.Login)).
			AddItems(
				openhab.NewItem(itemPrefix+"Balance", openhab.ItemTypeNumber).
					WithLabel("Balance [%.2f ₽]").
					WithIcon("price"),
			),
		openhab.NewChannel("Promise", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicPromise.Format(cfg.Login)).
			AddItems(
				openhab.NewItem(itemPrefix+"Promise", openhab.ItemTypeNumber).
					WithLabel("Promise [%.2f ₽]").
					WithIcon("price"),
			),
	)
}
