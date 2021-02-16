package pantum

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
	cfg := b.config()

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel("ProductID", openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicProductID.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"ProductID", openhab.ItemTypeString).
					WithLabel("Product ID"),
			),
		openhab.NewChannel("TonerRemain", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicTonerRemain.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"TonerRemain", openhab.ItemTypeNumber).
					WithLabel("Toner remain [%d %%]").
					WithIcon("humidity"),
			),
		openhab.NewChannel("PrinterStatus", openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicPrinterStatus.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"PrinterStatus", openhab.ItemTypeString).
					WithLabel("Printer status"),
			),
		openhab.NewChannel("CartridgeStatus", openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicCartridgeStatus.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"CartridgeStatus", openhab.ItemTypeString).
					WithLabel("Cartridge status"),
			),
		openhab.NewChannel("DrumStatus", openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicDrumStatus.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"DrumStatus", openhab.ItemTypeString).
					WithLabel("Drum status"),
			),
	)
}
