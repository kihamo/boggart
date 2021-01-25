package pulsar

import (
	"errors"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	channels := make([]*openhab.Channel, 0, b.config.InputsCount*2+7)

	for i := int64(1); i <= b.config.InputsCount; i++ {
		index := strconv.FormatInt(i, 10)
		idPulses := openhab.IDNormalizeCamelCase("Input " + index + " pulses")
		idVolume := openhab.IDNormalizeCamelCase("Input " + index + " volume")

		channels = append(channels,
			openhab.NewChannel(idPulses, openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicInputPulses.Format(sn, index)).
				AddItems(
					openhab.NewItem(itemPrefix+idPulses, openhab.ItemTypeNumber).
						WithLabel("Input "+index+" pulses [%d]"),
				),
			openhab.NewChannel(idVolume, openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicInputVolume.Format(sn, index)).
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
			WithStateTopic(b.config.TopicTemperatureIn.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idTemperatureIn, openhab.ItemTypeNumber).
					WithLabel("Temperature in [%.2f °C]").
					WithIcon("temperature_hot"),
			),
		openhab.NewChannel(idTemperatureOut, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicTemperatureOut.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idTemperatureOut, openhab.ItemTypeNumber).
					WithLabel("Temperature out [%.2f °C]").
					WithIcon("temperature_cold"),
			),
		openhab.NewChannel(idTemperatureDelta, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicTemperatureDelta.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idTemperatureDelta, openhab.ItemTypeNumber).
					WithLabel("Temperature delta [%.2f °C]").
					WithIcon("temperature"),
			),
		openhab.NewChannel(idEnergy, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicEnergy.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idEnergy, openhab.ItemTypeNumber).
					WithLabel("Energy [%.2f Gcal]").
					WithIcon("heating"),
			),
		openhab.NewChannel(idConsumption, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicConsumption.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idConsumption, openhab.ItemTypeNumber).
					WithLabel("Consumption [%.2f m3/h]").
					WithIcon("heating"),
			),
		openhab.NewChannel(idCapacity, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicCapacity.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idCapacity, openhab.ItemTypeNumber).
					WithLabel("Capacity [%.2f m3]").
					WithIcon("heating"),
			),
		openhab.NewChannel(idPower, openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicPower.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idPower, openhab.ItemTypeNumber).
					WithLabel("Power [%.2f Gcal/h]").
					WithIcon("heating"),
			),
	)

	return openhab.StepsByBind(b, nil, channels...)
}
