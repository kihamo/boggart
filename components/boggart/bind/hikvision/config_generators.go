package hikvision

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
	"github.com/kihamo/boggart/providers/hikvision/client/content_manager"
	"github.com/kihamo/boggart/providers/hikvision/client/ptz"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	const (
		idModel                = "Model"
		idFirmwareVersion      = "FirmwareVersion"
		idFirmwareReleasedDate = "FirmwareReleasedDate"
		idUpTime               = "UpTime"
		idMemoryUsage          = "MemoryUsage"
		idMemoryAvailable      = "MemoryAvailable"
	)

	channels := []*openhab.Channel{
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
		openhab.NewChannel(idFirmwareReleasedDate, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicStateFirmwareReleasedDate.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idFirmwareReleasedDate, openhab.ItemTypeString).
					WithLabel("Firmware release date"),
			),
		openhab.NewChannel(idUpTime, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicStateUpTime.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idUpTime, openhab.ItemTypeString).
					WithLabel("Uptime [%d s]").
					WithIcon("time"),
			),
		openhab.NewChannel(idMemoryUsage, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicStateMemoryUsage.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idMemoryUsage, openhab.ItemTypeNumber).
					WithLabel("Memory usage [JS(human_bytes.js):%s]").
					WithIcon("chart"),
			),
		openhab.NewChannel(idMemoryAvailable, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicStateMemoryAvailable.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idMemoryAvailable, openhab.ItemTypeNumber).
					WithLabel("Memory available [JS(human_bytes.js):%s]").
					WithIcon("chart"),
			),
	}

	if b.config.EventsEnabled {
		const idEvent = "Event"

		channels = append(channels,
			openhab.NewChannel(idEvent, openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicEvent.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idEvent, openhab.ItemTypeNumber).
						WithLabel("Event [%d]").
						WithIcon("motion"),
				),
		)
	}

	ctx := context.Background()

	if storage, err := b.client.ContentManager.GetStorage(content_manager.NewGetStorageParamsWithContext(ctx), nil); err == nil {
		const (
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
				openhab.NewChannel(id+idHDDCapacity, openhab.ChannelTypeNumber).
					WithStateTopic(b.config.TopicStateHDDCapacity.Format(sn, hdd.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idHDDCapacity, openhab.ItemTypeNumber).
							WithLabel("HDD capacity [JS(human_bytes.js):%s]").
							WithIcon("chart"),
					),
				openhab.NewChannel(id+idHDDUsage, openhab.ChannelTypeNumber).
					WithStateTopic(b.config.TopicStateHDDUsage.Format(sn, hdd.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idHDDUsage, openhab.ItemTypeNumber).
							WithLabel("HDD usage [JS(human_bytes.js):%s]").
							WithIcon("chart"),
					),
				openhab.NewChannel(id+idHDDFree, openhab.ChannelTypeNumber).
					WithStateTopic(b.config.TopicStateHDDFree.Format(sn, hdd.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idHDDFree, openhab.ItemTypeNumber).
							WithLabel("HDD free [JS(human_bytes.js):%s]").
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
			channels = append(channels,
				openhab.NewChannel(idPTZAbsolute, openhab.ChannelTypeString).
					WithCommandTopic(b.config.TopicPTZAbsolute.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZAbsolute, openhab.ItemTypeString).
							WithLabel("Absolute [%s]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZContinuous, openhab.ChannelTypeString).
					WithCommandTopic(b.config.TopicPTZContinuous.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZContinuous, openhab.ItemTypeString).
							WithLabel("Continuous [%s]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZRelative, openhab.ChannelTypeString).
					WithCommandTopic(b.config.TopicPTZRelative.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZRelative, openhab.ItemTypeString).
							WithLabel("Relative [%s]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZPreset, openhab.ChannelTypeString).
					WithCommandTopic(b.config.TopicPTZPreset.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZPreset, openhab.ItemTypeString).
							WithLabel("Preset [%s]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZMomentary, openhab.ChannelTypeString).
					WithCommandTopic(b.config.TopicPTZMomentary.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZMomentary, openhab.ItemTypeString).
							WithLabel("Momentary [%d]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZElevation, openhab.ChannelTypeNumber).
					WithStateTopic(b.config.TopicPTZStatusElevation.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZElevation, openhab.ItemTypeNumber).
							WithLabel("Elevation [%d]").
							WithIcon("movecontrol"),
					),
				openhab.NewChannel(idPTZAzimuth, openhab.ChannelTypeNumber).
					WithStateTopic(b.config.TopicPTZStatusAzimuth.Format(sn, ch.ID)).
					AddItems(
						openhab.NewItem(itemPrefix+idPTZAzimuth, openhab.ItemTypeNumber).
							WithLabel("Azimuth [%d]").
							WithIcon("movecontrol"),
					),
			)

			if ch.ZoomSupport {
				channels = append(channels,
					openhab.NewChannel(idPTZZoom, openhab.ChannelTypeNumber).
						WithStateTopic(b.config.TopicPTZStatusZoom.Format(sn, ch.ID)).
						AddItems(
							openhab.NewItem(itemPrefix+idPTZZoom, openhab.ItemTypeNumber).
								WithLabel("Zoom [%d]").
								WithIcon("zoom"),
						),
				)
			}
		}
	}

	return openhab.StepsByBind(b, []generators.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanBytes),
	}, channels...)
}
