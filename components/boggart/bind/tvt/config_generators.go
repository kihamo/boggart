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

	channels = append(channels,
		openhab.NewChannel("Model", openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicStateModel.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"Model", openhab.ItemTypeString).
					WithLabel("Model"),
			),
		openhab.NewChannel("FirmwareVersion", openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicStateFirmwareVersion.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"FirmwareVersion", openhab.ItemTypeString).
					WithLabel("Firmware version"),
			),
	)

	for _, disk := range disks.Payload.Content.Disks {
		id := openhab.IDNormalizeCamelCase(disk.SerialNum)

		channels = append(channels,
			openhab.NewChannel(id+"_HDDCapacity", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicStateHDDCapacity.Format(sn, disk.SerialNum)).
				AddItems(
					openhab.NewItem(itemPrefix+id+"_HDDCapacity", openhab.ItemTypeNumber).
						WithLabel("HDD capacity [JS(human_bytes.js):%s]"),
				),
			openhab.NewChannel(id+"_HDDUsage", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicStateHDDUsage.Format(sn, disk.SerialNum)).
				AddItems(
					openhab.NewItem(itemPrefix+id+"_HDDUsage", openhab.ItemTypeNumber).
						WithLabel("HDD usage [JS(human_bytes.js):%s]"),
				),
			openhab.NewChannel(id+"_HDDFree", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicStateHDDFree.Format(sn, disk.SerialNum)).
				AddItems(
					openhab.NewItem(itemPrefix+id+"_HDDFree", openhab.ItemTypeNumber).
						WithLabel("HDD free [JS(human_bytes.js):%s]"),
				),
		)
	}

	return openhab.StepsByBind(b, []generators.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanBytes),
	}, channels...)
}
