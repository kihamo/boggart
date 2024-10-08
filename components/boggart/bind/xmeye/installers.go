package xmeye

import (
	"context"
	"errors"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
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

	client, err := b.client(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()

	const (
		idEvent                = "Event"
		idModel                = "Model"
		idFirmwareVersion      = "FirmwareVersion"
		idFirmwareReleasedDate = "FirmwareReleasedDate"
		idUpTime               = "UpTime"
	)

	channels := []*openhab.Channel{
		openhab.NewChannel(idEvent, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicEvent.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idEvent, openhab.ItemTypeString).
					WithLabel("Annotation").
					WithIcon("text"),
			),
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
		openhab.NewChannel(idFirmwareReleasedDate, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicStateFirmwareReleasedDate.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareReleasedDate, openhab.ItemTypeString).
					WithLabel("Firmware release date"),
			),
		openhab.NewChannel(idUpTime, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicStateUpTime.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idUpTime, openhab.ItemTypeString).
					WithLabel("Uptime [%d s]").
					WithIcon("time"),
			),
	}

	if storage, err := client.StorageInfo(ctx); err == nil {
		const (
			idHDDCapacity = "HDDCapacity"
			idHDDUsage    = "HDDUsage"
			idHDDFree     = "HDDFree"
		)

		transformHumanBytes := openhab.StepDefaultTransformHumanBytes.Base()

		for _, s := range storage {
			for _, p := range s.Partition {
				if !p.IsCurrent {
					continue
				}

				id := openhab.IDNormalizeCamelCase(strconv.FormatUint(p.LogicSerialNo, 10))

				channels = append(channels,
					openhab.NewChannel(id+idHDDCapacity, openhab.ChannelTypeNumber).
						WithStateTopic(cfg.TopicStateHDDCapacity.Format(sn, p.LogicSerialNo)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idHDDCapacity, openhab.ItemTypeNumber).
								WithLabel("HDD capacity [JS("+transformHumanBytes+"):%s]").
								WithIcon("chart"),
						),
					openhab.NewChannel(id+idHDDUsage, openhab.ChannelTypeNumber).
						WithStateTopic(cfg.TopicStateHDDUsage.Format(sn, p.LogicSerialNo)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idHDDUsage, openhab.ItemTypeNumber).
								WithLabel("HDD usage [JS("+transformHumanBytes+"):%s]").
								WithIcon("chart"),
						),
					openhab.NewChannel(id+idHDDFree, openhab.ChannelTypeNumber).
						WithStateTopic(cfg.TopicStateHDDFree.Format(sn, p.LogicSerialNo)).
						AddItems(
							openhab.NewItem(itemPrefix+id+idHDDFree, openhab.ItemTypeNumber).
								WithLabel("HDD free [JS("+transformHumanBytes+"):%s]").
								WithIcon("chart"),
						),
				)
			}
		}
	}

	return openhab.StepsByBind(b, nil, channels...)
}
