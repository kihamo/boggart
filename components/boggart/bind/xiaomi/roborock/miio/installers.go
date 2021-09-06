package miio

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

const (
	SystemXiaomi installer.System = "Xiaomi"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
		SystemXiaomi,
	}
}

func (b *Bind) InstallerSteps(_ context.Context, system installer.System) ([]installer.Step, error) {
	if system == SystemXiaomi {
		return []installer.Step{{
			Description: "Factory Reset Gen 2",
			Content: `1. Нажимаем и держим левую и правую кнопки (локальная уборка и возврат на базу).
2. Нажимаем кнопку сброса под крышкой (Reset) и держим 3-5 секунд.
3. Отпускаем кнопку сброса (Reset), продолжаем держать левую и правую кнопки, пока не заморгает быстро центральная кнопка
4. После некоторого ожидания робот заговорит на китайском.`,
		}}, nil
	}

	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()

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
			WithStateTopic(cfg.TopicBattery.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idBattery, openhab.ItemTypeNumber).
					WithLabel("Battery [%d %%]").
					WithIcon("batterylevel"),
			),
		openhab.NewChannel(idFanPower, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicFanPower.Format(sn)).
			WithCommandTopic(cfg.TopicSetFanPower.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idFanPower, openhab.ItemTypeNumber).
					WithLabel("Fan []").
					WithIcon("fan"),
			),
		openhab.NewChannel(idVolume, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicVolume.Format(sn)).
			WithCommandTopic(cfg.TopicSetVolume.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idVolume, openhab.ItemTypeNumber).
					WithLabel("Volume").
					WithIcon("soundvolume"),
			),
		openhab.NewChannel(idVolumeTest, openhab.ChannelTypeSwitch).
			WithStateTopic(cfg.TopicTestVolume.Format(sn)).
			WithCommandTopic(cfg.TopicTestVolume.Format(sn)).
			WithOn("true").
			WithOff(done).
			AddItems(
				openhab.NewItem(itemPrefix+idVolumeTest, openhab.ItemTypeSwitch).
					WithLabel("Volume test").
					WithIcon("soundvolume"),
			),
		openhab.NewChannel(idFindMe, openhab.ChannelTypeSwitch).
			WithStateTopic(cfg.TopicFind.Format(sn)).
			WithCommandTopic(cfg.TopicFind.Format(sn)).
			WithOn("true").
			WithOff(done).
			AddItems(
				openhab.NewItem(itemPrefix+idFindMe, openhab.ItemTypeSwitch).
					WithLabel("Find me").
					WithIcon("zoom"),
			),
		openhab.NewChannel(idStatus, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicState.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idStatus, openhab.ItemTypeString).
					WithLabel("Status").
					WithIcon("text"),
			),
		openhab.NewChannel(idError, openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicError.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idError, openhab.ItemTypeString).
					WithLabel("Error").
					WithIcon("error"),
			),
		openhab.NewChannel(idAction, openhab.ChannelTypeString).
			WithCommandTopic(cfg.TopicAction.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idAction, openhab.ItemTypeString).
					WithLabel("Action []").
					WithIcon("movecontrol"),
			),
		openhab.NewChannel(idLastCleanCompleted, openhab.ChannelTypeContact).
			WithStateTopic(cfg.TopicLastCleanCompleted.Format(sn)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+idLastCleanCompleted, openhab.ItemTypeContact).
					WithLabel("Last clean completed").
					WithIcon("contact"),
			),
		openhab.NewChannel(idLastCleanArea, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicLastCleanArea.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idLastCleanArea, openhab.ItemTypeNumber).
					WithLabel("Area").
					WithIcon("text"),
			),
		openhab.NewChannel(idLastCleanStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicLastCleanStartDateTime.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idLastCleanStart, openhab.ItemTypeDateTime).
					WithLabel("Start [%1$td.%1$tm.%1$tY at %1$tH:%1$tM:%1$tS]").
					WithIcon("calendar"),
			),
		openhab.NewChannel(idLastCleanEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicLastCleanEndDateTime.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idLastCleanEnd, openhab.ItemTypeDateTime).
					WithLabel("End [%1$td.%1$tm.%1$tY at %1$tH:%1$tM:%1$tS]").
					WithIcon("calendar"),
			),
		openhab.NewChannel(idLastCleanDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicLastCleanDuration.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idLastCleanDuration, openhab.ItemTypeNumber).
					WithLabel("Duration").
					WithIcon("time"),
			),
		openhab.NewChannel(idConsumableFilter, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicConsumableFilter.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idConsumableFilter, openhab.ItemTypeNumber).
					WithLabel("Filter").
					WithIcon("time"),
			),
		openhab.NewChannel(idConsumableBrushMain, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicConsumableBrushMain.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idConsumableBrushMain, openhab.ItemTypeNumber).
					WithLabel("Main brush").
					WithIcon("time"),
			),
		openhab.NewChannel(idConsumableBrushSide, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicConsumableBrushSide.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idConsumableBrushSide, openhab.ItemTypeNumber).
					WithLabel("Side brush").
					WithIcon("time"),
			),
		openhab.NewChannel(idConsumableSensor, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicConsumableSensor.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idConsumableSensor, openhab.ItemTypeNumber).
					WithLabel("Sensor").
					WithIcon("time"),
			),
	)
}
