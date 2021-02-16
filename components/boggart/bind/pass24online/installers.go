package pass24online

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

	const idFeedEvent = "FeedEvent"

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idFeedEvent, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicFeedEvent.Format(cfg.Phone)).
			AddItems(
				openhab.NewItem(itemPrefix+idFeedEvent, openhab.ItemTypeString).
					WithLabel("Feed event").
					WithIcon("text"),
			),
	)
}
