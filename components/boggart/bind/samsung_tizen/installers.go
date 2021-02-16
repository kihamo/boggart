package tizen

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
	mac := meta.MACAsString()
	if mac == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	cfg := b.config()

	const (
		idPower = "Power"
		idKey   = "Key"
		idID    = "ID"
		idModel = "Model"
	)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idPower, openhab.ChannelTypeSwitch).
			WithCommandTopic(cfg.TopicPower.Format(mac)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idPower, openhab.ItemTypeSwitch).
					WithLabel("Power"),
			),
		openhab.NewChannel(idKey, openhab.ChannelTypeString).
			WithCommandTopic(cfg.TopicKey.Format(mac)).
			AddItems(
				openhab.NewItem(itemPrefix+idKey, openhab.ItemTypeString).
					WithLabel("Key").
					WithIcon("network"),
			),
		openhab.NewChannel(idID, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicDeviceID.Format(mac)).
			AddItems(
				openhab.NewItem(itemPrefix+idID, openhab.ItemTypeString).
					WithLabel("ID").
					WithIcon("text"),
			),
		openhab.NewChannel(idModel, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicDeviceModelName.Format(mac)).
			AddItems(
				openhab.NewItem(itemPrefix+idModel, openhab.ItemTypeString).
					WithLabel("Model").
					WithIcon("text"),
			),
	)
}
