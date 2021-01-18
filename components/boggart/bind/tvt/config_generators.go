package tvt

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
	"github.com/kihamo/boggart/providers/tvt/client/storage"
)

func (b *Bind) GenerateConfigOpenHab() []generators.Step {
	opts, err := b.MQTT().ClientOptions()
	if err != nil {
		return nil
	}

	meta := b.Meta()
	filePrefix := openhab.FilePrefixFromBindMeta(meta)
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	sn := meta.SerialNumber()
	if sn == "" {
		sn = "+"
	}

	broker := openhab.BrokerFromClientOptionsReader(opts)

	steps := []generators.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanBytes),
		{
			FilePath: openhab.DirectoryThings + "broker.things",
			Content:  broker.String(),
		},
	}

	thing := openhab.GenericThingFromBindMeta(meta).
		WithBroker(broker).
		AddChannels(
			openhab.BindStatusChannel(meta),
			openhab.BindSerialNumberChannel(meta),
			openhab.BindMACChannel(meta),
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

	disks, err := b.client.Storage.GetStorageInfo(storage.NewGetStorageInfoParamsWithContext(context.Background()), nil)
	if err == nil {
		for _, disk := range disks.Payload.Content.Disks {
			id := openhab.IDNormalizeCamelCase(disk.SerialNum)

			thing.AddChannels(
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
	}

	if content := thing.String(); content != "" {
		steps = append(steps, generators.Step{
			FilePath: openhab.DirectoryThings + filePrefix + ".things",
			Content:  content,
		})
	}

	if content := thing.Items().String(); content != "" {
		steps = append(steps, generators.Step{
			FilePath: openhab.DirectoryItems + filePrefix + ".items",
			Content:  content,
		})
	}

	return steps
}
