package octoprint

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

func (b *Bind) InstallerSteps(ctx context.Context, _ installer.System) ([]installer.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
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
	)

	b.devicesMutex.RLock()
	for device := range b.devices {
		id := openhab.IDNormalizeCamelCase(device) + "_"

		if b.TemperatureFromMQTT() {
			topic := b.TemperatureTopic().Format(device)

			channels = append(channels,
				openhab.NewChannel(id+idTemperatureActual, openhab.ChannelTypeNumber).
					WithStateTopic(topic).
					WithTransformationPattern("JSONPATH:$.actual").
					AddItems(
						openhab.NewItem(itemPrefix+id+idTemperatureActual, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" actual [%.2f °C]").
							WithIcon("temperature_cold"),
					),
				openhab.NewChannel(id+idTemperatureTarget, openhab.ChannelTypeNumber).
					WithStateTopic(topic).
					WithTransformationPattern("JSONPATH:$.target").
					AddItems(
						openhab.NewItem(itemPrefix+id+idTemperatureTarget, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" target [%.2f °C]").
							WithIcon("temperature_hot"),
					),
			)
		} else {
			channels = append(channels,
				openhab.NewChannel(id+idTemperatureActual, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicTemperatureActual.Format(device)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idTemperatureActual, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" actual [%.2f °C]").
							WithIcon("temperature_cold"),
					),
				openhab.NewChannel(id+idTemperatureTarget, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicTemperatureTarget.Format(device)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idTemperatureTarget, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" target [%.2f °C]").
							WithIcon("temperature_hot"),
					),
				openhab.NewChannel(id+idTemperatureOffset, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicTemperatureOffset.Format(device)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idTemperatureOffset, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" offset [%.2f °C]").
							WithIcon("temperature"),
					),
			)
		}
	}
	b.devicesMutex.RUnlock()

	transformHumanSeconds := openhab.StepDefaultTransformHumanSeconds.Base()
	transformHumanBytes := openhab.StepDefaultTransformHumanBytes.Base()

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
				WithStateTopic(cfg.TopicState).
				AddItems(
					openhab.NewItem(itemPrefix+idStatus, openhab.ItemTypeString).
						WithLabel("Status").
						WithIcon("text"),
				),
			openhab.NewChannel(idJobFileName, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicJobFileName).
				AddItems(
					openhab.NewItem(itemPrefix+idJobFileName, openhab.ItemTypeString).
						WithLabel("File name [%s]").
						WithIcon("text"),
				),
			openhab.NewChannel(idJobFileSize, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicJobFileSize).
				AddItems(
					openhab.NewItem(itemPrefix+idJobFileSize, openhab.ItemTypeNumber).
						WithLabel("File size [JS("+transformHumanBytes+"):%s]").
						WithIcon("chart"),
				),
			openhab.NewChannel(idJobProgress, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicJobProgress).
				AddItems(
					openhab.NewItem(itemPrefix+idJobProgress, openhab.ItemTypeNumber).
						WithLabel("Progress [%d %%]").
						WithIcon("humidity"),
				),
			openhab.NewChannel(idJobTime, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicJobTime).
				AddItems(
					openhab.NewItem(itemPrefix+idJobTime, openhab.ItemTypeNumber).
						WithLabel("Time [JS("+transformHumanSeconds+"):%s]").
						WithIcon("time"),
				),
			openhab.NewChannel(idJobTimeLeft, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicJobTimeLeft).
				AddItems(
					openhab.NewItem(itemPrefix+idJobTimeLeft, openhab.ItemTypeNumber).
						WithLabel("Time left [JS("+transformHumanSeconds+"):%s]").
						WithIcon("time"),
				),
		)
	}

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanBytes),
		openhab.StepDefault(openhab.StepDefaultTransformHumanSeconds),
	}, channels...)
}
