package service

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
	meta := b.Meta()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	const idLatency = "Latency"

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idLatency, openhab.ChannelTypeNumber).
			WithStateTopic(b.config().TopicLatency.Format(meta.ID())).
			AddItems(
				openhab.NewItem(itemPrefix+idLatency, openhab.ItemTypeNumber).
					WithLabel("Address "+b.address+" ping latency [%d ms]").
					WithIcon("clock"),
			),
	)
}
