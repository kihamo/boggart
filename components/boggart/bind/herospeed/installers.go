package herospeed

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
		idModel           = "Model"
		idFirmwareVersion = "FirmwareVersion"
	)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idModel, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicStateModel.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idModel, openhab.ItemTypeString).
					WithLabel("Model"),
			),
		openhab.NewChannel(idFirmwareVersion, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicStateFirmwareVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareVersion, openhab.ItemTypeString).
					WithLabel("Firmware version"),
			),
	)
}
