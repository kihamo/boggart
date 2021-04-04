package text2speech

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
		idBinaryResponse    = "BinaryResponse"
		idBinarySendText    = "BinarySendText"
		idBinarySendOptions = "BinarySendOptions"
		idURLResponse       = "URLResponse"
		idURLSendText       = "URLSendText"
		idURLSendOptions    = "URLSendOptions"
	)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idBinaryResponse, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicResponseBinary.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idBinaryResponse, openhab.ItemTypeString).
					WithLabel("Binary").
					WithIcon("text"),
			),
		openhab.NewChannel(idBinarySendText, openhab.ChannelTypeString).
			WithCommandTopic(cfg.TopicGenerateBinaryText.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idBinarySendText, openhab.ItemTypeString).
					WithLabel("Send text").
					WithIcon("text"),
			),
		openhab.NewChannel(idBinarySendOptions, openhab.ChannelTypeString).
			WithCommandTopic(cfg.TopicGenerateBinaryOptions.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idBinarySendOptions, openhab.ItemTypeString).
					WithLabel("Send JSON").
					WithIcon("text"),
			),
		openhab.NewChannel(idURLResponse, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicResponseURL.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idURLResponse, openhab.ItemTypeString).
					WithLabel("URL").
					WithIcon("text"),
			),
		openhab.NewChannel(idURLSendText, openhab.ChannelTypeString).
			WithCommandTopic(cfg.TopicGenerateURLText.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idURLSendText, openhab.ItemTypeString).
					WithLabel("Send text").
					WithIcon("text"),
			),
		openhab.NewChannel(idURLSendOptions, openhab.ChannelTypeString).
			WithCommandTopic(cfg.TopicGenerateURLOptions.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idURLSendOptions, openhab.ItemTypeString).
					WithLabel("Send JSON").
					WithIcon("text"),
			),
	)
}
