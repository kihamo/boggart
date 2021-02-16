package miio

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

	const (
		idBattery             = "Battery"
		idFanPower            = "FanPower"
		idVolume              = "Volume"
		idVolumeTest          = "VolumeTest"
		idFindMe              = "FindMe"
		idStatus              = "Status"
		idError               = "Error"
		idAction              = "Action"
		idLastCleanCompleted  = "LastCleanCompleted"
		idLastCleanArea       = "LastCleanArea"
		idLastCleanStart      = "LastCleanStart"
		idLastCleanEnd        = "LastCleanEnd"
		idLastCleanDuration   = "LastCleanDuration"
		idConsumableFilter    = "ConsumableFilter"
		idConsumableBrushMain = "ConsumableBrushMain"
		idConsumableBrushSide = "ConsumableBrushSide"
		idConsumableSensor    = "ConsumableSensor"
	)

	done := string(payloadDone)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idBattery, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicBattery.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idBattery, openhab.ItemTypeNumber).
					WithLabel("Battery [%d %%]").
					WithIcon("batterylevel"),
			),
		openhab.NewChannel(idFanPower, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicFanPower.Format(sn)).
			WithCommandTopic(b.config.TopicSetFanPower.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idFanPower, openhab.ItemTypeNumber).
					WithLabel("Fan []").
					WithIcon("fan"),
			),
		openhab.NewChannel(idVolume, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicVolume.Format(sn)).
			WithCommandTopic(b.config.TopicSetVolume.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idVolume, openhab.ItemTypeNumber).
					WithLabel("Volume").
					WithIcon("soundvolume"),
			),
		openhab.NewChannel(idVolumeTest, openhab.ChannelTypeSwitch).
			WithStateTopic(b.config.TopicTestVolume.Format(sn)).
			WithCommandTopic(b.config.TopicTestVolume.Format(sn)).
			WithOn("true").
			WithOff(done).
			AddItems(
				openhab.NewItem(itemPrefix+idVolumeTest, openhab.ItemTypeSwitch).
					WithLabel("Volume test").
					WithIcon("soundvolume"),
			),
		openhab.NewChannel(idFindMe, openhab.ChannelTypeSwitch).
			WithStateTopic(b.config.TopicFind.Format(sn)).
			WithCommandTopic(b.config.TopicFind.Format(sn)).
			WithOn("true").
			WithOff(done).
			AddItems(
				openhab.NewItem(itemPrefix+idFindMe, openhab.ItemTypeSwitch).
					WithLabel("Find me").
					WithIcon("zoom"),
			),
		openhab.NewChannel(idStatus, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicState.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idStatus, openhab.ItemTypeString).
					WithLabel("Status").
					WithIcon("text"),
			),
		openhab.NewChannel(idError, openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicError.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idError, openhab.ItemTypeString).
					WithLabel("Error").
					WithIcon("error"),
			),
		openhab.NewChannel(idAction, openhab.ChannelTypeString).
			WithCommandTopic(b.config.TopicAction.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idAction, openhab.ItemTypeString).
					WithLabel("Action []").
					WithIcon("movecontrol"),
			),
		openhab.NewChannel(idLastCleanCompleted, openhab.ChannelTypeContact).
			WithStateTopic(b.config.TopicLastCleanCompleted.Format(sn)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idLastCleanCompleted, openhab.ItemTypeContact).
					WithLabel("Last clean completed").
					WithIcon("contact"),
			),
		openhab.NewChannel(idLastCleanArea, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicLastCleanArea.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idLastCleanArea, openhab.ItemTypeNumber).
					WithLabel("Area").
					WithIcon("text"),
			),
		openhab.NewChannel(idLastCleanStart, openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicLastCleanStartDateTime.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idLastCleanStart, openhab.ItemTypeDateTime).
					WithLabel("Start [%1$td.%1$tm.%1$tY at %1$tH:%1$tM:%1$tS]").
					WithIcon("calendar"),
			),
		openhab.NewChannel(idLastCleanEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicLastCleanEndDateTime.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idLastCleanEnd, openhab.ItemTypeDateTime).
					WithLabel("End [%1$td.%1$tm.%1$tY at %1$tH:%1$tM:%1$tS]").
					WithIcon("calendar"),
			),
		openhab.NewChannel(idLastCleanDuration, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicLastCleanDuration.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idLastCleanDuration, openhab.ItemTypeNumber).
					WithLabel("Duration").
					WithIcon("time"),
			),
		openhab.NewChannel(idConsumableFilter, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicConsumableFilter.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idConsumableFilter, openhab.ItemTypeNumber).
					WithLabel("Filter").
					WithIcon("time"),
			),
		openhab.NewChannel(idConsumableBrushMain, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicConsumableBrushMain.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idConsumableBrushMain, openhab.ItemTypeNumber).
					WithLabel("Main brush").
					WithIcon("time"),
			),
		openhab.NewChannel(idConsumableBrushSide, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicConsumableBrushSide.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idConsumableBrushSide, openhab.ItemTypeNumber).
					WithLabel("Side brush").
					WithIcon("time"),
			),
		openhab.NewChannel(idConsumableSensor, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicConsumableSensor.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idConsumableSensor, openhab.ItemTypeNumber).
					WithLabel("Sensor").
					WithIcon("time"),
			),
	)
}
