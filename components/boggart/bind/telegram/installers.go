package telegram

import (
	"context"
	"errors"

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
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	const (
		idReceiveMessage = "ReceiveMessage"
		idReceiveAudio   = "ReceiveAudio"
		idReceiveVoice   = "ReceiveVoice"
	)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idReceiveMessage, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicReceiveMessage.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idReceiveMessage, openhab.ItemTypeString).
					WithLabel("Text message [%s]").
					WithIcon("text"),
			),
		openhab.NewChannel(idReceiveAudio, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicReceiveAudio.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idReceiveAudio, openhab.ItemTypeString).
					WithLabel("Audio file [%s]").
					WithIcon("soundvolume"),
			),
		openhab.NewChannel(idReceiveVoice, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicReceiveVoice.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idReceiveVoice, openhab.ItemTypeString).
					WithLabel("Voice file [%s]").
					WithIcon("recorder"),
			),
	)
}
