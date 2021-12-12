package pulsar

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

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()

	channels := make([]*openhab.Channel, 0, cfg.InputsCount*2+7)

	for i := int64(1); i <= cfg.InputsCount; i++ {
		index := strconv.FormatInt(i, 10)
		idPulses := openhab.IDNormalizeCamelCase("Input " + index + " pulses")
		idVolume := openhab.IDNormalizeCamelCase("Input " + index + " volume")

		channels = append(channels,
			openhab.NewChannel(idPulses, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicInputPulses.Format(sn, index)).
				AddItems(
					openhab.NewItem(itemPrefix+idPulses, openhab.ItemTypeNumber).
						WithLabel("Input "+index+" pulses [%d]"),
				),
			openhab.NewChannel(idVolume, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicInputVolume.Format(sn, index)).
				AddItems(
					openhab.NewItem(itemPrefix+idVolume, openhab.ItemTypeNumber).
						WithLabel("Input "+index+" volume [%.2f m3]"),
				),
		)
	}

	const (
		idTemperatureIn    = "TemperatureIn"
		idTemperatureOut   = "TemperatureOut"
		idTemperatureDelta = "TemperatureDelta"
		idEnergy           = "Energy"
		idConsumption      = "Consumption"
		idCapacity         = "Capacity"
		idPower            = "Power"
	)

	channels = append(channels,
		openhab.NewChannel(idTemperatureIn, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicTemperatureIn.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idTemperatureIn, openhab.ItemTypeNumber).
					WithLabel("Temperature in [%.2f °C]").
					WithIcon("temperature_hot"),
			),
		openhab.NewChannel(idTemperatureOut, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicTemperatureOut.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idTemperatureOut, openhab.ItemTypeNumber).
					WithLabel("Temperature out [%.2f °C]").
					WithIcon("temperature_cold"),
			),
		openhab.NewChannel(idTemperatureDelta, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicTemperatureDelta.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idTemperatureDelta, openhab.ItemTypeNumber).
					WithLabel("Temperature delta [%.2f °C]").
					WithIcon("temperature"),
			),
		openhab.NewChannel(idEnergy, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicEnergy.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idEnergy, openhab.ItemTypeNumber).
					WithLabel("Energy [%.3f Gcal]").
					WithIcon("heating"),
			),
		openhab.NewChannel(idConsumption, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicConsumption.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idConsumption, openhab.ItemTypeNumber).
					WithLabel("Consumption [%.3f m3/h]").
					WithIcon("heating"),
			),
		openhab.NewChannel(idCapacity, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCapacity.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idCapacity, openhab.ItemTypeNumber).
					WithLabel("Capacity [%.3f m3]").
					WithIcon("heating"),
			),
		openhab.NewChannel(idPower, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicPower.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPower, openhab.ItemTypeNumber).
					WithLabel("Power [%.4f Gcal/h]").
					WithIcon("heating"),
			),
	)

	return openhab.StepsByBind(b, nil, channels...)
}
