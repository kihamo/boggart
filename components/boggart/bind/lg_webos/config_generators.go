package webos

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel("Application", openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicStateApplication).
			WithCommandTopic(b.config.TopicApplication).
			AddItems(
				openhab.NewItem(itemPrefix+"Application", openhab.ItemTypeString).
					WithLabel("Application").
					WithIcon("screen"),
			),
		openhab.NewChannel("Mute", openhab.ChannelTypeSwitch).
			WithStateTopic(b.config.TopicStateMute).
			WithCommandTopic(b.config.TopicMute).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+"Mute", openhab.ItemTypeSwitch).
					WithLabel("Mute").
					WithIcon("soundvolume_mute"),
			),
		openhab.NewChannel("Volume", openhab.ChannelTypeDimmer).
			WithStateTopic(b.config.TopicStateVolume).
			WithCommandTopic(b.config.TopicVolume).
			WithMin(0).
			WithMax(100).
			WithStep(1).
			AddItems(
				openhab.NewItem(itemPrefix+"Volume", openhab.ItemTypeDimmer).
					WithLabel("Volume").
					WithIcon("soundvolume"),
			),
		openhab.NewChannel("VolumeUp", openhab.ChannelTypeDimmer).
			WithCommandTopic(b.config.TopicVolumeUp),
		openhab.NewChannel("VolumeDown", openhab.ChannelTypeDimmer).
			WithCommandTopic(b.config.TopicVolumeDown),
		openhab.NewChannel("Power", openhab.ChannelTypeSwitch).
			WithStateTopic(b.config.TopicStatePower).
			WithCommandTopic(b.config.TopicPower).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+"Power", openhab.ItemTypeSwitch).
					WithLabel("Power"),
			),
		openhab.NewChannel("Channel", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicStateChannelNumber).
			AddItems(
				openhab.NewItem(itemPrefix+"Channel", openhab.ItemTypeNumber).
					WithLabel("Channel").
					WithIcon("video"),
			),
		openhab.NewChannel("Toast", openhab.ChannelTypeString).
			WithCommandTopic(b.config.TopicToast),
	)
}
