package sun

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
	id := b.Meta().ID()

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel("Nadir", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNadir.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"Nadir", openhab.ItemTypeDateTime).
					WithLabel("Nadir [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel("NightBeforeStart", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNightBeforeStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NightBeforeStart", openhab.ItemTypeDateTime).
					WithLabel("Night before start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel("NightBeforeEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNightBeforeEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NightBeforeEnd", openhab.ItemTypeDateTime).
					WithLabel("Night before end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel("NightBeforeDuration", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicNightBeforeDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NightBeforeDuration", openhab.ItemTypeNumber).
					WithLabel("Night before duration [%d s]").
					WithIcon("moon"),
			),
		openhab.NewChannel("AstronomicalDawnStart", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicAstronomicalDawnStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDawnStart", openhab.ItemTypeDateTime).
					WithLabel("Astronomical dawn start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("AstronomicalDawnEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicAstronomicalDawnEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDawnEnd", openhab.ItemTypeDateTime).
					WithLabel("Astronomical dawn end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("AstronomicalDawnDuration", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicAstronomicalDawnDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDawnDuration", openhab.ItemTypeNumber).
					WithLabel("Astronomical dawn duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("NauticalDawnStart", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNauticalDawnStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDawnStart", openhab.ItemTypeDateTime).
					WithLabel("Nautical dawn start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("NauticalDawnEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNauticalDawnEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDawnEnd", openhab.ItemTypeDateTime).
					WithLabel("Nautical dawn end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("NauticalDawnDuration", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicNauticalDawnDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDawnDuration", openhab.ItemTypeNumber).
					WithLabel("Nautical dawn duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("CivilDawnStart", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicCivilDawnStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDawnStart", openhab.ItemTypeDateTime).
					WithLabel("Civil dawn start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("CivilDawnEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicCivilDawnEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDawnEnd", openhab.ItemTypeDateTime).
					WithLabel("Civil dawn end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("CivilDawnDuration", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCivilDawnDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDawnDuration", openhab.ItemTypeNumber).
					WithLabel("Civil dawn duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("RiseStart", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicRiseStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"RiseStart", openhab.ItemTypeDateTime).
					WithLabel("Sunrise start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("RiseEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicRiseEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"RiseEnd", openhab.ItemTypeDateTime).
					WithLabel("Sunrise end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("RiseDuration", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicRiseDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"RiseDuration", openhab.ItemTypeNumber).
					WithLabel("Sunrise duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel("SolarNoon", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicSolarNoon.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"SolarNoon", openhab.ItemTypeDateTime).
					WithLabel("Solar noon [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sun"),
			),
		openhab.NewChannel("SetStart", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicSetStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"SetStart", openhab.ItemTypeDateTime).
					WithLabel("Sunset start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("SetEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicSetEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"SetEnd", openhab.ItemTypeDateTime).
					WithLabel("Sunset end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("SetDuration", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicSetDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"SetDuration", openhab.ItemTypeNumber).
					WithLabel("Sunset duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("CivilDuskStart", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicCivilDuskStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDuskStart", openhab.ItemTypeDateTime).
					WithLabel("Civil dusk start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("CivilDuskEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicCivilDuskEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDuskEnd", openhab.ItemTypeDateTime).
					WithLabel("Civil dusk end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("CivilDuskDuration", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCivilDuskDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"CivilDuskDuration", openhab.ItemTypeNumber).
					WithLabel("Civil dusk duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("NauticalDuskStart", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNauticalDuskStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDuskStart", openhab.ItemTypeDateTime).
					WithLabel("Nautical dusk start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("NauticalDuskEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNauticalDuskEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDuskEnd", openhab.ItemTypeDateTime).
					WithLabel("Nautical dusk end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("NauticalDuskDuration", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicNauticalDuskDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NauticalDuskDuration", openhab.ItemTypeNumber).
					WithLabel("Nautical dusk duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("AstronomicalDuskStart", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicAstronomicalDuskStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDuskStart", openhab.ItemTypeDateTime).
					WithLabel("Astronomical dusk start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("AstronomicalDuskEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicAstronomicalDuskEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDuskEnd", openhab.ItemTypeDateTime).
					WithLabel("Astronomical dusk end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("AstronomicalDuskDuration", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicAstronomicalDuskDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"AstronomicalDuskDuration", openhab.ItemTypeNumber).
					WithLabel("Astronomical dusk duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel("NightAfterStart", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNightAfterStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NightAfterStart", openhab.ItemTypeDateTime).
					WithLabel("Night after start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel("NightAfterEnd", openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNightAfterEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NightAfterEnd", openhab.ItemTypeDateTime).
					WithLabel("Night after end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel("NightAfterDuration", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicNightAfterDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+"NightAfterDuration", openhab.ItemTypeNumber).
					WithLabel("Night after duration [%d s]").
					WithIcon("moon"),
			),
	)
}
