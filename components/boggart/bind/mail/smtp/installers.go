package smtp

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

const (
	SystemGMail installer.System = "GMail"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
		SystemGMail,
	}
}

func (b *Bind) InstallerSteps(_ context.Context, s installer.System) ([]installer.Step, error) {
	if s == SystemGMail {
		return []installer.Step{{
			Description: "Create password for apps",
			Content: `1. Open https://myaccount.google.com/apppasswords
2. Select app Mail and device Other
3. Copy created password and use dsn in format smtp://myemail%40gmail.com:[created_password]@smtp.gmail.com:587`,
		}}, nil
	}

	meta := b.Meta()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	const idSend = "SendJSON"

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idSend, openhab.ChannelTypeString).
			WithCommandTopic(b.config().TopicSend.Format(meta.ID())).
			AddItems(
				openhab.NewItem(itemPrefix+idSend, openhab.ItemTypeString),
			),
	)
}
