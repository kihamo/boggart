package apcupsd

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

func (b *Bind) InstallerSteps(ctx context.Context, system installer.System) ([]installer.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()
	channels := make([]*openhab.Channel, 0, 0)

	channels = append(channels,
		openhab.NewChannel(VariableUPSName, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicVariable.Format(sn, VariableUPSName)).
			AddItems(
				openhab.NewItem(itemPrefix+VariableUPSName, openhab.ItemTypeString).
					WithLabel("UPS name"),
			),
	)

	return openhab.StepsByBind(b, nil, channels...)
}
