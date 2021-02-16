package tvt

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
	"github.com/kihamo/boggart/providers/tvt/client/storage"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(ctx context.Context, _ installer.System) ([]installer.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	disks, err := b.client.Storage.GetStorageInfo(storage.NewGetStorageInfoParamsWithContext(ctx), nil)
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()
	channels := make([]*openhab.Channel, 0, 2+len(disks.Payload.Content.Disks)*3)

	const (
		idModel           = "Model"
		idFirmwareVersion = "FirmwareVersion"
		idHDDCapacity     = "HDDCapacity"
		idHDDUsage        = "HDDUsage"
		idHDDFree         = "HDDFree"
	)

	channels = append(channels,
		openhab.NewChannel(idModel, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicStateModel.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idModel, openhab.ItemTypeString).
					WithLabel("Model"),
			),
		openhab.NewChannel(idFirmwareVersion, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicStateFirmwareVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareVersion, openhab.ItemTypeString).
					WithLabel("Firmware version"),
			),
	)

	for _, disk := range disks.Payload.Content.Disks {
		id := openhab.IDNormalizeCamelCase(disk.SerialNum)

		channels = append(channels,
			openhab.NewChannel(id+idHDDCapacity, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicStateHDDCapacity.Format(sn, disk.SerialNum)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idHDDCapacity, openhab.ItemTypeNumber).
						WithLabel("HDD capacity [JS(human_bytes.js):%s]").
						WithIcon("chart"),
				),
			openhab.NewChannel(id+idHDDUsage, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicStateHDDUsage.Format(sn, disk.SerialNum)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idHDDUsage, openhab.ItemTypeNumber).
						WithLabel("HDD usage [JS(human_bytes.js):%s]").
						WithIcon("chart"),
				),
			openhab.NewChannel(id+idHDDFree, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicStateHDDFree.Format(sn, disk.SerialNum)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idHDDFree, openhab.ItemTypeNumber).
						WithLabel("HDD free [JS(human_bytes.js):%s]").
						WithIcon("chart"),
				),
		)
	}

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanBytes),
	}, channels...)
}
