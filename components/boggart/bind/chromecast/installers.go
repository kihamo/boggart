package chromecast

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
	id := meta.ID()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()

	const (
		idMute   = "Mute"
		idVolume = "Volume"
		idPause  = "Pause"
		idStop   = "Stop"
		idResume = "Resume"
		idSeek   = "Seek"
		idPlay   = "Play"
		idAction = "Action"
		idStatus = "Status"
	)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idMute, openhab.ChannelTypeSwitch).
			WithStateTopic(cfg.TopicStateMute.Format(id)).
			WithCommandTopic(cfg.TopicMute.Format(id)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idMute, openhab.ItemTypeSwitch).
					WithLabel("Mute").
					WithIcon("soundvolume_mute"),
			),
		openhab.NewChannel(idVolume, openhab.ChannelTypeDimmer).
			WithStateTopic(cfg.TopicStateVolume.Format(id)).
			WithCommandTopic(cfg.TopicVolume.Format(id)).
			WithMin(0).
			WithMax(100).
			WithStep(1).
			AddItems(
				openhab.NewItem(itemPrefix+idVolume, openhab.ItemTypeDimmer).
					WithLabel("Volume").
					WithIcon("soundvolume"),
			),
		openhab.NewChannel(idPause, openhab.ChannelTypeSwitch).
			WithCommandTopic(cfg.TopicPause.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idPause, openhab.ItemTypeSwitch).
					WithLabel("Pause").
					WithIcon("soundvolume_mute"),
			),
		openhab.NewChannel(idStop, openhab.ChannelTypeSwitch).
			WithCommandTopic(cfg.TopicStop.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idStop, openhab.ItemTypeSwitch).
					WithLabel("Stop").
					WithIcon("soundvolume_mute"),
			),
		openhab.NewChannel(idResume, openhab.ChannelTypeSwitch).
			WithCommandTopic(cfg.TopicResume.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idResume, openhab.ItemTypeSwitch).
					WithLabel("Resume").
					WithIcon("soundvolume"),
			),
		openhab.NewChannel(idSeek, openhab.ChannelTypeNumber).
			WithCommandTopic(cfg.TopicSeek.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idSeek, openhab.ItemTypeNumber).
					WithLabel("Seek").
					WithIcon("time"),
			),
		openhab.NewChannel(idPlay, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicStateContent.Format(id)).
			WithCommandTopic(cfg.TopicPlay.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idPlay, openhab.ItemTypeString).
					WithLabel("Play").
					WithIcon("text"),
			),
		openhab.NewChannel(idAction, openhab.ChannelTypeString).
			WithCommandTopic(cfg.TopicAction.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idAction, openhab.ItemTypeString).
					WithLabel("Action").
					WithIcon("text"),
			),
		openhab.NewChannel(idStatus, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicStateStatus.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idStatus, openhab.ItemTypeString).
					WithLabel("Status").
					WithIcon("text"),
			),
	)
}
