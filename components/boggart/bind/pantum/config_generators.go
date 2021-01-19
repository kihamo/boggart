package pantum

import (
	"errors"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel("ProductID", openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicProductID.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"ProductID", openhab.ItemTypeString).
					WithLabel("Product ID"),
			),
		openhab.NewChannel("TonerRemain", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicTonerRemain.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"TonerRemain", openhab.ItemTypeNumber).
					WithLabel("Toner remain [%d %%]").
					WithIcon("humidity"),
			),
		openhab.NewChannel("PrinterStatus", openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicPrinterStatus.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"PrinterStatus", openhab.ItemTypeString).
					WithLabel("Printer status"),
			),
		openhab.NewChannel("CartridgeStatus", openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicCartridgeStatus.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"CartridgeStatus", openhab.ItemTypeString).
					WithLabel("Cartridge status"),
			),
		openhab.NewChannel("DrumStatus", openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicDrumStatus.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"DrumStatus", openhab.ItemTypeString).
					WithLabel("Drum status"),
			),
	)
}
