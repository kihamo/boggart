package sun

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel("Nadir", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicNadir).
			AddItems(
				openhab.NewItem(itemPrefix+"Nadir", openhab.ItemTypeDateTime).
					WithLabel("Nadir [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel("NightBeforeStart", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicNightBeforeStart).
			AddItems(
				openhab.NewItem(itemPrefix+"NightBeforeStart", openhab.ItemTypeDateTime).
					WithLabel("Night before start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel("NightBeforeEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicNightBeforeEnd).
			AddItems(
				openhab.NewItem(itemPrefix+"NightBeforeEnd", openhab.ItemTypeDateTime).
					WithLabel("Night before end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel("NightBeforeDuration", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicNightBeforeDuration).
			AddItems(
				openhab.NewItem(itemPrefix+"NightBeforeDuration", openhab.ItemTypeNumber).
					WithLabel("Night before duration [%d s]").
					WithIcon("moon"),
			),
		openhab.NewChannel("AstronomicalDawnStart", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicAstronomicalDawnStart).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDawnStart", openhab.ItemTypeDateTime).
					WithLabel("Astronomical dawn start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("AstronomicalDawnEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicAstronomicalDawnEnd).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDawnEnd", openhab.ItemTypeDateTime).
					WithLabel("Astronomical dawn end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("AstronomicalDawnDuration", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicAstronomicalDawnDuration).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDawnDuration", openhab.ItemTypeNumber).
					WithLabel("Astronomical dawn duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("NauticalDawnStart", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicNauticalDawnStart).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDawnStart", openhab.ItemTypeDateTime).
					WithLabel("Nautical dawn start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("NauticalDawnEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicNauticalDawnEnd).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDawnEnd", openhab.ItemTypeDateTime).
					WithLabel("Nautical dawn end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("NauticalDawnDuration", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicNauticalDawnDuration).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDawnDuration", openhab.ItemTypeNumber).
					WithLabel("Nautical dawn duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("CivilDawnStart", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicCivilDawnStart).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDawnStart", openhab.ItemTypeDateTime).
					WithLabel("Civil dawn start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("CivilDawnEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicCivilDawnEnd).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDawnEnd", openhab.ItemTypeDateTime).
					WithLabel("Civil dawn end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("CivilDawnDuration", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicCivilDawnDuration).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDawnDuration", openhab.ItemTypeNumber).
					WithLabel("Civil dawn duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("RiseStart", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicRiseStart).
			AddItems(
				openhab.NewItem(itemPrefix+"RiseStart", openhab.ItemTypeDateTime).
					WithLabel("Sunrise start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("RiseEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicRiseEnd).
			AddItems(
				openhab.NewItem(itemPrefix+"RiseEnd", openhab.ItemTypeDateTime).
					WithLabel("Sunrise end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("RiseDuration", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicRiseDuration).
			AddItems(
				openhab.NewItem(itemPrefix+"RiseDuration", openhab.ItemTypeNumber).
					WithLabel("Sunrise duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("SolarNoon", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicSolarNoon).
			AddItems(
				openhab.NewItem(itemPrefix+"SolarNoon", openhab.ItemTypeDateTime).
					WithLabel("Solar noon [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sun"),
			),
		openhab.NewChannel("SetStart", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicSetStart).
			AddItems(
				openhab.NewItem(itemPrefix+"SetStart", openhab.ItemTypeDateTime).
					WithLabel("Sunset start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("SetEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicSetEnd).
			AddItems(
				openhab.NewItem(itemPrefix+"SetEnd", openhab.ItemTypeDateTime).
					WithLabel("Sunset end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("SetDuration", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicSetDuration).
			AddItems(
				openhab.NewItem(itemPrefix+"SetDuration", openhab.ItemTypeNumber).
					WithLabel("Sunset duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("CivilDuskStart", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicCivilDuskStart).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDuskStart", openhab.ItemTypeDateTime).
					WithLabel("Civil dusk start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("CivilDuskEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicCivilDuskEnd).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDuskEnd", openhab.ItemTypeDateTime).
					WithLabel("Civil dusk end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("CivilDuskDuration", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicCivilDuskDuration).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDuskDuration", openhab.ItemTypeNumber).
					WithLabel("Civil dusk duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("NauticalDuskStart", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicNauticalDuskStart).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDuskStart", openhab.ItemTypeDateTime).
					WithLabel("Nautical dusk start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("NauticalDuskEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicNauticalDuskEnd).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDuskEnd", openhab.ItemTypeDateTime).
					WithLabel("Nautical dusk end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("NauticalDuskDuration", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicNauticalDuskDuration).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDuskDuration", openhab.ItemTypeNumber).
					WithLabel("Nautical dusk duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("AstronomicalDuskStart", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicAstronomicalDuskStart).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDuskStart", openhab.ItemTypeDateTime).
					WithLabel("Astronomical dusk start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("AstronomicalDuskEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicAstronomicalDuskEnd).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDuskEnd", openhab.ItemTypeDateTime).
					WithLabel("Astronomical dusk end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("AstronomicalDuskDuration", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicAstronomicalDuskDuration).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDuskDuration", openhab.ItemTypeNumber).
					WithLabel("Astronomical dusk duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("NightAfterStart", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicNightAfterStart).
			AddItems(
				openhab.NewItem(itemPrefix+"NightAfterStart", openhab.ItemTypeDateTime).
					WithLabel("Night after start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel("NightAfterEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(b.config.TopicNightAfterEnd).
			AddItems(
				openhab.NewItem(itemPrefix+"NightAfterEnd", openhab.ItemTypeDateTime).
					WithLabel("Night after end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel("NightAfterDuration", openhab.ChannelTypeNumber).
			WithStateTopic(b.config.TopicNightAfterDuration).
			AddItems(
				openhab.NewItem(itemPrefix+"NightAfterDuration", openhab.ItemTypeNumber).
					WithLabel("Night after duration [%d s]").
					WithIcon("moon"),
			),
	)
}
