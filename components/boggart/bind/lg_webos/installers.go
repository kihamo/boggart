package webos

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
	mac := b.Meta().MACAsString()

	const (
		idApplication = "Application"
		idMute        = "Mute"
		idVolume      = "Volume"
		idVolumeUp    = "VolumeUp"
		idVolumeDown  = "VolumeDown"
		idPower       = "Power"
		idChannel     = "Channel"
		idToast       = "Toast"
	)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idApplication, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicStateApplication.Format(mac)).
			WithCommandTopic(cfg.TopicApplication.Format(mac)).
			AddItems(
				openhab.NewItem(itemPrefix+idApplication, openhab.ItemTypeString).
					WithLabel("Application").
					WithIcon("screen"),
			),
		openhab.NewChannel(idMute, openhab.ChannelTypeSwitch).
			WithStateTopic(cfg.TopicStateMute.Format(mac)).
			WithCommandTopic(cfg.TopicMute.Format(mac)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idMute, openhab.ItemTypeSwitch).
					WithLabel("Mute").
					WithIcon("soundvolume_mute"),
			),
		openhab.NewChannel(idVolume, openhab.ChannelTypeDimmer).
			WithStateTopic(cfg.TopicStateVolume.Format(mac)).
			WithCommandTopic(cfg.TopicVolume.Format(mac)).
			WithMin(0).
			WithMax(100).
			WithStep(1).
			AddItems(
				openhab.NewItem(itemPrefix+idVolume, openhab.ItemTypeDimmer).
					WithLabel("Volume").
					WithIcon("soundvolume"),
			),
		openhab.NewChannel(idVolumeUp, openhab.ChannelTypeDimmer).
			WithCommandTopic(cfg.TopicVolumeUp.Format(mac)),
		openhab.NewChannel(idVolumeDown, openhab.ChannelTypeDimmer).
			WithCommandTopic(cfg.TopicVolumeDown.Format(mac)),
		openhab.NewChannel(idPower, openhab.ChannelTypeSwitch).
			WithStateTopic(cfg.TopicStatePower.Format(mac)).
			WithCommandTopic(cfg.TopicPower.Format(mac)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idPower, openhab.ItemTypeSwitch).
					WithLabel("Power"),
			),
		openhab.NewChannel(idChannel, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicStateChannelNumber.Format(mac)).
			AddItems(
				openhab.NewItem(itemPrefix+idChannel, openhab.ItemTypeNumber).
					WithLabel("Channel").
					WithIcon("video"),
			),
		openhab.NewChannel(idToast, openhab.ChannelTypeString).
			WithCommandTopic(cfg.TopicToast.Format(mac)).
			AddItems(
				openhab.NewItem(itemPrefix+idToast, openhab.ItemTypeString).
					WithLabel("Toast").
					WithIcon("text"),
			),
	)
}
