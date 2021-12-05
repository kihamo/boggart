package hikvision

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
	"github.com/kihamo/boggart/providers/hikvision/client/content_manager"
	"github.com/kihamo/boggart/providers/hikvision/client/ptz"
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

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()

	const (
		idModel                = "Model"
		idFirmwareVersion      = "FirmwareVersion"
		idFirmwareReleasedDate = "FirmwareReleasedDate"
		idUpTime               = "UpTime"
		idMemoryUsage          = "MemoryUsage"
		idMemoryAvailable      = "MemoryAvailable"
	)

	transformHumanSeconds := openhab.StepDefaultTransformHumanSeconds.Base()
	transformHumanBytes := openhab.StepDefaultTransformHumanBytes.Base()

	channels := []*openhab.Channel{
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
					WithLabel("Uptime [JS(" + transformHumanSeconds + "):%s]").
					WithIcon("time"),
			),
		openhab.NewChannel(idMemoryUsage, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicStateMemoryUsage.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idMemoryUsage, openhab.ItemTypeNumber).
					WithLabel("Memory usage [JS(" + transformHumanBytes + "):%s]").
					WithIcon("chart"),
			),
		openhab.NewChannel(idMemoryAvailable, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicStateMemoryAvailable.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idMemoryAvailable, openhab.ItemTypeNumber).
					WithLabel("Memory available [JS(" + transformHumanBytes + "):%s]").
					WithIcon("chart"),
			),
	}

	if cfg.EventsEnabled {
		const idEvent = "Event"

		channels = append(channels,
			openhab.NewChannel(idEvent, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicEvent.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idEvent, openhab.ItemTypeNumber).
						WithLabel("Event [%d]").
						WithIcon("motion"),
				),
		)
	}

	if storage, err := b.client.ContentManager.GetStorage(content_manager.NewGetStorageParamsWithContext(ctx), nil); err == nil {
		const (
			idHDDStatus   = "HDDStatus"
			idHDDCapacity = "HDDCapacity"
			idHDDUsage    = "HDDUsage"
			idHDDFree     = "HDDFree"
		)

		for _, hdd := range storage.Payload.HddList {
			if hdd.Name == "" {
				continue
			}

			id := openhab.IDNormalizeCamelCase(hdd.Name)

			channels = append(channels,
				openhab.NewChannel(id+idHDDStatus, openhab.ChannelTypeString).
					WithStateTopic(cfg.TopicStateHDDStatus.Format(sn, hdd.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idHDDStatus, openhab.ItemTypeString).
							WithLabel("HDD status"),
					),
				openhab.NewChannel(id+idHDDCapacity, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicStateHDDCapacity.Format(sn, hdd.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idHDDCapacity, openhab.ItemTypeNumber).
							WithLabel("HDD capacity [JS("+transformHumanBytes+"):%s]").
							WithIcon("chart"),
					),
				openhab.NewChannel(id+idHDDUsage, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicStateHDDUsage.Format(sn, hdd.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idHDDUsage, openhab.ItemTypeNumber).
							WithLabel("HDD usage [JS("+transformHumanBytes+"):%s]").
							WithIcon("chart"),
					),
				openhab.NewChannel(id+idHDDFree, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicStateHDDFree.Format(sn, hdd.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idHDDFree, openhab.ItemTypeNumber).
							WithLabel("HDD free [JS("+transformHumanBytes+"):%s]").
							WithIcon("chart"),
					),
			)
		}
	}

	if ptzChannels, err := b.client.Ptz.GetPtzChannels(ptz.NewGetPtzChannelsParamsWithContext(ctx), nil); err == nil {
		const (
			idPTZAbsolute   = "PTZAbsolute"
			idPTZContinuous = "PTZContinuous"
			idPTZRelative   = "PTZRelative"
			idPTZPreset     = "PTZPreset"
			idPTZMomentary  = "PTZMomentary"
			idPTZElevation  = "PTZElevation"
			idPTZAzimuth    = "PTZAzimuth"
			idPTZZoom       = "PTZZoom"
		)

		for _, ch := range ptzChannels.Payload {
			if ch.ID < 1 {
				continue
			}

			channels = append(channels,
				openhab.NewChannel(idPTZAbsolute, openhab.ChannelTypeString).
					WithCommandTopic(cfg.TopicPTZAbsolute.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZAbsolute, openhab.ItemTypeString).
							WithLabel("Absolute [%s]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZContinuous, openhab.ChannelTypeString).
					WithCommandTopic(cfg.TopicPTZContinuous.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZContinuous, openhab.ItemTypeString).
							WithLabel("Continuous [%s]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZRelative, openhab.ChannelTypeString).
					WithCommandTopic(cfg.TopicPTZRelative.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZRelative, openhab.ItemTypeString).
							WithLabel("Relative [%s]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZPreset, openhab.ChannelTypeString).
					WithCommandTopic(cfg.TopicPTZPreset.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZPreset, openhab.ItemTypeString).
							WithLabel("Preset [%s]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZMomentary, openhab.ChannelTypeString).
					WithCommandTopic(cfg.TopicPTZMomentary.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZMomentary, openhab.ItemTypeString).
							WithLabel("Momentary [%d]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZElevation, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicPTZStatusElevation.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZElevation, openhab.ItemTypeNumber).
							WithLabel("Elevation [%d]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZAzimuth, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicPTZStatusAzimuth.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZAzimuth, openhab.ItemTypeNumber).
							WithLabel("Azimuth [%d]").
							WithIcon("movecontrol"),
					),
			)

			if ch.ZoomSupport {
				channels = append(channels,
					openhab.NewChannel(idPTZZoom, openhab.ChannelTypeNumber).
						WithStateTopic(cfg.TopicPTZStatusZoom.Format(sn, ch.ID)).
						AddItems(
							openhab.NewItem(itemPrefix+idPTZZoom, openhab.ItemTypeNumber).
								WithLabel("Zoom [%d]").
								WithIcon("zoom"),
						),
				)
			}
		}
	}

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanBytes),
		openhab.StepDefault(openhab.StepDefaultTransformHumanSeconds),
	}, channels...)
}
