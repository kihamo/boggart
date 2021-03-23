package ledwifi

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
	id := meta.ID()

	const (
		idPower = "Power"
		idColor = "ColorHSB"
		idMode  = "Mode"
		idSpeed = "Speed"
	)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idPower, openhab.ChannelTypeSwitch).
			WithStateTopic(cfg.TopicStatePower.Format(id)).
			WithCommandTopic(cfg.TopicPower.Format(id)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idPower, openhab.ItemTypeSwitch).
					WithLabel("Power"),
			),
		openhab.NewChannel(idColor, openhab.ChannelTypeColor).
			WithStateTopic(cfg.TopicStateColorHSV.Format(id)).
			WithCommandTopic(cfg.TopicColor.Format(id)).
			WithColorMode("hsb").
			AddItems(
				openhab.NewItem(itemPrefix+idColor, openhab.ItemTypeColor).
					WithLabel("Color").
					WithIcon("rgb"),
			),
		openhab.NewChannel(idMode, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicStateMode.Format(id)).
			WithCommandTopic(cfg.TopicMode.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idMode, openhab.ItemTypeString).
					WithLabel("Mode").
					WithIcon("text"),
			),
		openhab.NewChannel(idSpeed, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicStateSpeed.Format(id)).
			WithCommandTopic(cfg.TopicSpeed.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idSpeed, openhab.ItemTypeNumber).
					WithLabel("Speed").
					WithIcon("heating"),
			),
	)
}
