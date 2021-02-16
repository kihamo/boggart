package smtp

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
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	const idSend = "SendJSON"

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idSend, openhab.ChannelTypeString).
			WithCommandTopic(b.config().TopicSend.Format(meta.ID())).
			AddItems(
				openhab.NewItem(itemPrefix+idSend, openhab.ItemTypeString),
			),
	)
}
