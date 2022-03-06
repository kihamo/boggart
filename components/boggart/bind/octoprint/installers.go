package octoprint

import (
	"context"
	"sort"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(_ context.Context, _ installer.System) ([]installer.Step, error) {
	meta := b.Meta()
	bindID := meta.ID()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()
	channels := make([]*openhab.Channel, 0)

	const (
		idTemperatureActual = "TemperatureActual"
		idTemperatureTarget = "TemperatureTarget"
		idTemperatureOffset = "TemperatureOffset"
		idStatus            = "Status"
		idJobFileName       = "JobFileName"
		idJobFileSize       = "JobFileSize"
		idJobProgress       = "JobProgress"
		idJobTime           = "JobTime"
		idJobTimeLeft       = "JobTimeLeft"
		idLayerTotal        = "LayerTotal"
		idLayerCurrent      = "LayerCurrent"
		idHeightTotal       = "HeightTotal"
		idHeightCurrent     = "HeightCurrent"
		idCommand           = "Command"
	)

	b.devicesMutex.RLock()
	devices := make([]string, 0, len(b.devices))
	for device := range b.devices {
		devices = append(devices, device)
	}
	b.devicesMutex.RUnlock()
	sort.Strings(devices)

	for _, device := range devices {
		deviceID := openhab.IDNormalizeCamelCase(device) + "_"

		if b.TemperatureFromMQTT() {
			topic := b.TemperatureTopic().Format(device)

			channels = append(channels,
				openhab.NewChannel(deviceID+idTemperatureActual, openhab.ChannelTypeNumber).
					WithStateTopic(topic).
					WithTransformationPattern("JSONPATH:$.actual").
					AddItems(
						openhab.NewItem(itemPrefix+deviceID+idTemperatureActual, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" actual [%.2f °C]").
							WithIcon("temperature_cold"),
					),
				openhab.NewChannel(deviceID+idTemperatureTarget, openhab.ChannelTypeNumber).
					WithStateTopic(topic).
					WithTransformationPattern("JSONPATH:$.target").
					AddItems(
						openhab.NewItem(itemPrefix+deviceID+idTemperatureTarget, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" target [%.2f °C]").
							WithIcon("temperature_hot"),
					),
			)
		} else {
			channels = append(channels,
				openhab.NewChannel(deviceID+idTemperatureActual, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicTemperatureActual.Format(bindID, device)).
					AddItems(
						openhab.NewItem(itemPrefix+deviceID+idTemperatureActual, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" actual [%.2f °C]").
							WithIcon("temperature_cold"),
					),
				openhab.NewChannel(deviceID+idTemperatureTarget, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicTemperatureTarget.Format(bindID, device)).
					AddItems(
						openhab.NewItem(itemPrefix+deviceID+idTemperatureTarget, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" target [%.2f °C]").
							WithIcon("temperature_hot"),
					),
				openhab.NewChannel(deviceID+idTemperatureOffset, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicTemperatureOffset.Format(bindID, device)).
					AddItems(
						openhab.NewItem(itemPrefix+deviceID+idTemperatureOffset, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" offset [%.2f °C]").
							WithIcon("temperature"),
					),
			)
		}
	}

	transformHumanSeconds := openhab.StepDefaultTransformHumanSeconds.Base()
	transformHumanBytes := openhab.StepDefaultTransformHumanBytes.Base()

	// Job
	if b.JobFromMQTT() {
		topic := b.JobTopic()

		channels = append(channels,
			openhab.NewChannel(idStatus, openhab.ChannelTypeString).
				WithStateTopic(topic).
				WithTransformationPattern("JSONPATH:$.printer_data.state.text").
				AddItems(
					openhab.NewItem(itemPrefix+idStatus, openhab.ItemTypeString).
						WithLabel("Status").
						WithIcon("text"),
				),
			openhab.NewChannel(idJobFileName, openhab.ChannelTypeString).
				WithStateTopic(topic).
				WithTransformationPattern("JSONPATH:$.printer_data.job.file.name").
				AddItems(
					openhab.NewItem(itemPrefix+idJobFileName, openhab.ItemTypeString).
						WithLabel("File name [%s]").
						WithIcon("text"),
				),
			openhab.NewChannel(idJobFileSize, openhab.ChannelTypeNumber).
				WithStateTopic(topic).
				WithTransformationPattern("JSONPATH:$.printer_data.job.file.size").
				AddItems(
					openhab.NewItem(itemPrefix+idJobFileSize, openhab.ItemTypeNumber).
						WithLabel("File size [JS("+transformHumanBytes+"):%s]").
						WithIcon("chart"),
				),
			openhab.NewChannel(idJobProgress, openhab.ChannelTypeNumber).
				WithStateTopic(topic).
				WithTransformationPattern("JSONPATH:$.progress").
				AddItems(
					openhab.NewItem(itemPrefix+idJobProgress, openhab.ItemTypeNumber).
						WithLabel("Progress [%d %%]").
						WithIcon("humidity"),
				),
			openhab.NewChannel(idJobTime, openhab.ChannelTypeNumber).
				WithStateTopic(topic).
				WithTransformationPattern("JSONPATH:$.printer_data.progress.printTime").
				AddItems(
					openhab.NewItem(itemPrefix+idJobTime, openhab.ItemTypeNumber).
						WithLabel("Time [JS("+transformHumanSeconds+"):%s]").
						WithIcon("time"),
				),
			openhab.NewChannel(idJobTimeLeft, openhab.ChannelTypeNumber).
				WithStateTopic(topic).
				WithTransformationPattern("JSONPATH:$.printer_data.progress.printTimeLeft").
				AddItems(
					openhab.NewItem(itemPrefix+idJobTimeLeft, openhab.ItemTypeNumber).
						WithLabel("Time left [JS("+transformHumanSeconds+"):%s]").
						WithIcon("time"),
				),
		)
	} else {
		channels = append(channels,
			openhab.NewChannel(idStatus, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicState.Format(bindID)).
				AddItems(
					openhab.NewItem(itemPrefix+idStatus, openhab.ItemTypeString).
						WithLabel("Status").
						WithIcon("text"),
				),
			openhab.NewChannel(idJobFileName, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicJobFileName.Format(bindID)).
				AddItems(
					openhab.NewItem(itemPrefix+idJobFileName, openhab.ItemTypeString).
						WithLabel("File name [%s]").
						WithIcon("text"),
				),
			openhab.NewChannel(idJobFileSize, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicJobFileSize.Format(bindID)).
				AddItems(
					openhab.NewItem(itemPrefix+idJobFileSize, openhab.ItemTypeNumber).
						WithLabel("File size [JS("+transformHumanBytes+"):%s]").
						WithIcon("chart"),
				),
			openhab.NewChannel(idJobProgress, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicJobProgress.Format(bindID)).
				AddItems(
					openhab.NewItem(itemPrefix+idJobProgress, openhab.ItemTypeNumber).
						WithLabel("Progress [%d %%]").
						WithIcon("humidity"),
				),
			openhab.NewChannel(idJobTime, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicJobTime.Format(bindID)).
				AddItems(
					openhab.NewItem(itemPrefix+idJobTime, openhab.ItemTypeNumber).
						WithLabel("Time [JS("+transformHumanSeconds+"):%s]").
						WithIcon("time"),
				),
			openhab.NewChannel(idJobTimeLeft, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicJobTimeLeft.Format(bindID)).
				AddItems(
					openhab.NewItem(itemPrefix+idJobTimeLeft, openhab.ItemTypeNumber).
						WithLabel("Time left [JS("+transformHumanSeconds+"):%s]").
						WithIcon("time"),
				),
		)
	}

	// Layer & Height
	if b.DisplayLayerProgressEnabled() {
		channels = append(channels,
			openhab.NewChannel(idLayerTotal, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicLayerTotal.Format(bindID)).
				AddItems(
					openhab.NewItem(itemPrefix+idLayerTotal, openhab.ItemTypeNumber).
						WithLabel("Total layers [%d]").
						WithIcon("niveau"),
				),
			openhab.NewChannel(idLayerCurrent, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicLayerCurrent.Format(bindID)).
				AddItems(
					openhab.NewItem(itemPrefix+idLayerCurrent, openhab.ItemTypeNumber).
						WithLabel("Current layer [%d]").
						WithIcon("niveau"),
				),
			openhab.NewChannel(idHeightTotal, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicHeightTotal.Format(bindID)).
				AddItems(
					openhab.NewItem(itemPrefix+idHeightTotal, openhab.ItemTypeNumber).
						WithLabel("Total height [%.2f mm]").
						WithIcon("niveau"),
				),
			openhab.NewChannel(idHeightCurrent, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicHeightCurrent.Format(bindID)).
				AddItems(
					openhab.NewItem(itemPrefix+idHeightCurrent, openhab.ItemTypeNumber).
						WithLabel("Current height [%.2f mm]").
						WithIcon("niveau"),
				),
		)
	}

	// Commands
	for _, command := range b.Commands() {
		commandID := idCommand +
			openhab.IDNormalizeCamelCase(command.Source) + "_" +
			openhab.IDNormalizeCamelCase(command.Action)

		channels = append(channels,
			openhab.NewChannel(commandID, openhab.ChannelTypeSwitch).
				WithCommandTopic(cfg.TopicCommand.Format(bindID, command.Source, command.Action)).
				WithOn("true").
				WithOff("false").
				WithTrigger(true).
				AddItems(
					openhab.NewItem(itemPrefix+commandID, openhab.ItemTypeSwitch).
						WithLabel(command.Name),
				),
		)
	}

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanBytes),
		openhab.StepDefault(openhab.StepDefaultTransformHumanSeconds),
	}, channels...)
}
