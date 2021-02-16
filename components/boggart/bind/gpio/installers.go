package gpio

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
	cfg := b.config()
	pinNumber := b.pin.Number()
	const id = "GPIO"

	var channel *openhab.Channel

	if b.Mode() != ModeOut {
		channel = openhab.NewChannel(id, openhab.ChannelTypeContact).
			AddItems(
				openhab.NewItem(itemPrefix+id, openhab.ItemTypeContact).
					WithLabel(meta.Description() + "[]"),
			)
	} else {
		channel = openhab.NewChannel(id, openhab.ChannelTypeSwitch).
			WithCommandTopic(cfg.TopicPinSet.Format(pinNumber)).
			AddItems(
				openhab.NewItem(itemPrefix+id, openhab.ItemTypeSwitch).
					WithLabel(meta.Description() + "[]"),
			)
	}

	channel.WithStateTopic(cfg.TopicPinState.Format(pinNumber)).
		WithLabel(meta.Description()).
		WithOn("true").
		WithOff("false")

	return openhab.StepsByBind(b, nil, channel)
}
