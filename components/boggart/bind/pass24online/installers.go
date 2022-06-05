package pass24online

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	cfg := b.config()

	const (
		idModelName   = "ModelName"
		idPlateNumber = "PlateNumber"
		idMessage     = "Message"
		idStatus      = "Status"
		idDatetime    = "Datetime"
	)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idModelName, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicFeedEvent.Format(cfg.Phone)).
			WithTransformationPattern("JSONPATH:$.model_name").
			AddItems(
				openhab.NewItem(itemPrefix+idModelName, openhab.ItemTypeString).
					WithLabel("Model name").
					WithIcon("text"),
			),
		openhab.NewChannel(idPlateNumber, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicFeedEvent.Format(cfg.Phone)).
			WithTransformationPattern("JSONPATH:$.plate_number").
			AddItems(
				openhab.NewItem(itemPrefix+idPlateNumber, openhab.ItemTypeString).
					WithLabel("Plate number").
					WithIcon("text"),
			),
		openhab.NewChannel(idMessage, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicFeedEvent.Format(cfg.Phone)).
			WithTransformationPattern("JSONPATH:$.message").
			AddItems(
				openhab.NewItem(itemPrefix+idMessage, openhab.ItemTypeString).
					WithLabel("Message").
					WithIcon("garagedoor"),
			),
		openhab.NewChannel(idStatus, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicFeedEvent.Format(cfg.Phone)).
			WithTransformationPattern("JSONPATH:$.status").
			AddItems(
				openhab.NewItem(itemPrefix+idStatus, openhab.ItemTypeString).
					WithLabel("Status").
					WithIcon("text"),
			),
		openhab.NewChannel(idDatetime, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicFeedEvent.Format(cfg.Phone)).
			WithTransformationPattern("JSONPATH:$.datetime").
			AddItems(
				openhab.NewItem(itemPrefix+idDatetime, openhab.ItemTypeDateTime).
					WithLabel("Date [%1$td.%1$tm.%1$tY %1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("time"),
			),
	)
}
