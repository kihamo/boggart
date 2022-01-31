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
	)

	b.devicesMutex.RLock()
	for device := range b.devices {
		id := openhab.IDNormalizeCamelCase(device) + "_"

		if b.TemperatureFromMQTT() {
			channels = append(channels,
				openhab.NewChannel(id+idTemperatureActual, openhab.ChannelTypeNumber).
					WithStateTopic(b.TemperatureTopic().Format(device)).
					WithTransformationPattern("JSONPATH:$.actual").
					AddItems(
						openhab.NewItem(itemPrefix+id+idTemperatureActual, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" actual [%.2f °C]").
							WithIcon("temperature"),
					),
				openhab.NewChannel(id+idTemperatureTarget, openhab.ChannelTypeNumber).
					WithStateTopic(b.TemperatureTopic().Format(device)).
					WithTransformationPattern("JSONPATH:$.target").
					AddItems(
						openhab.NewItem(itemPrefix+id+idTemperatureTarget, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" target [%.2f °C]").
							WithIcon("temperature"),
					),
			)
		} else {
			channels = append(channels,
				openhab.NewChannel(id+idTemperatureActual, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicStateTemperatureActual.Format(device)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idTemperatureActual, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" actual [%.2f °C]").
							WithIcon("temperature"),
					),
				openhab.NewChannel(id+idTemperatureTarget, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicStateTemperatureTarget.Format(device)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idTemperatureTarget, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" target [%.2f °C]").
							WithIcon("temperature"),
					),
				openhab.NewChannel(id+idTemperatureOffset, openhab.ChannelTypeNumber).
					WithStateTopic(cfg.TopicStateTemperatureOffset.Format(device)).
					AddItems(
						openhab.NewItem(itemPrefix+id+idTemperatureOffset, openhab.ItemTypeNumber).
							WithLabel("Temperature "+device+" offset [%.2f °C]").
							WithIcon("temperature"),
					),
			)
		}
	}
	b.devicesMutex.RUnlock()

	return openhab.StepsByBind(b, nil, channels...)
}
