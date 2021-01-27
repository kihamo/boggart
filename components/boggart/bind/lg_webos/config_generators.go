package webos

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

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
			WithStateTopic(b.config.TopicStateApplication).
			WithCommandTopic(b.config.TopicApplication).
			AddItems(
				openhab.NewItem(itemPrefix+idApplication, openhab.ItemTypeString).
					WithLabel("Application").
					WithIcon("screen"),
			),
		openhab.NewChannel(idMute, openhab.ChannelTypeSwitch).
			WithStateTopic(b.config.TopicStateMute).
			WithCommandTopic(b.config.TopicMute).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idMute, openhab.ItemTypeSwitch).
					WithLabel("Mute").
					WithIcon("soundvolume_mute"),
			),
		openhab.NewChannel(idVolume, openhab.ChannelTypeDimmer).
			WithStateTopic(b.config.TopicStateVolume).
			WithCommandTopic(b.config.TopicVolume).
			WithMin(0).
			WithMax(100).
			WithStep(1).
			AddItems(
				openhab.NewItem(itemPrefix+idVolume, openhab.ItemTypeDimmer).
					WithLabel("Volume").
					WithIcon("soundvolume"),
			),
		openhab.NewChannel(idVolumeUp, openhab.ChannelTypeDimmer).
			WithCommandTopic(b.config.TopicVolumeUp),
		openhab.NewChannel(idVolumeDown, openhab.ChannelTypeDimmer).
			WithCommandTopic(b.config.TopicVolumeDown),
		openhab.NewChannel(idPower, openhab.ChannelTypeSwitch).
			WithStateTopic(b.config.TopicStatePower).
			WithCommandTopic(b.config.TopicPower).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idPower, openhab.ItemTypeSwitch).
					WithLabel("Power"),
			),
		openhab.NewChannel(idChannel, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicStateChannelNumber).
			AddItems(
				openhab.NewItem(itemPrefix+idChannel, openhab.ItemTypeNumber).
					WithLabel("Channel").
					WithIcon("video"),
			),
		openhab.NewChannel(idToast, openhab.ChannelTypeString).
			WithCommandTopic(b.config.TopicToast).
			AddItems(
				openhab.NewItem(itemPrefix+idToast, openhab.ItemTypeString).
					WithLabel("Toast").
					WithIcon("text"),
			),
	)
}
