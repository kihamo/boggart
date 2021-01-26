package tvt

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
	"github.com/kihamo/boggart/providers/tvt/client/storage"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	disks, err := b.client.Storage.GetStorageInfo(storage.NewGetStorageInfoParamsWithContext(context.Background()), nil)
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
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

	for _, disk := range disks.Payload.Content.Disks {
		id := openhab.IDNormalizeCamelCase(disk.SerialNum)

		channels = append(channels,
			openhab.NewChannel(id+idHDDCapacity, openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicStateHDDCapacity.Format(sn, disk.SerialNum)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idHDDCapacity, openhab.ItemTypeNumber).
						WithLabel("HDD capacity [JS(human_bytes.js):%s]").
						WithIcon("chart"),
				),
			openhab.NewChannel(id+idHDDUsage, openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicStateHDDUsage.Format(sn, disk.SerialNum)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idHDDUsage, openhab.ItemTypeNumber).
						WithLabel("HDD usage [JS(human_bytes.js):%s]").
						WithIcon("chart"),
				),
			openhab.NewChannel(id+idHDDFree, openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicStateHDDFree.Format(sn, disk.SerialNum)).
				AddItems(
					openhab.NewItem(itemPrefix+id+idHDDFree, openhab.ItemTypeNumber).
						WithLabel("HDD free [JS(human_bytes.js):%s]").
						WithIcon("chart"),
				),
		)
	}

	return openhab.StepsByBind(b, []generators.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanBytes),
	}, channels...)
}
