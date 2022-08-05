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

	const (
		idNadir                    = "Nadir"
		idNightBeforeStart         = "NightBeforeStart"
		idNightBeforeEnd           = "NightBeforeEnd"
		idNightBeforeDuration      = "NightBeforeDuration"
		idAstronomicalDawnStart    = "AstronomicalDawnStart"
		idAstronomicalDawnEnd      = "AstronomicalDawnEnd"
		idAstronomicalDawnDuration = "AstronomicalDawnDuration"
		idNauticalDawnStart        = "NauticalDawnStart"
		idNauticalDawnEnd          = "NauticalDawnEnd"
		idNauticalDawnDuration     = "NauticalDawnDuration"
		idCivilDawnStart           = "CivilDawnStart"
		idCivilDawnEnd             = "CivilDawnEnd"
		idCivilDawnDuration        = "CivilDawnDuration"
		idRiseStart                = "RiseStart"
		idRiseEnd                  = "RiseEnd"
		idRiseDuration             = "RiseDuration"
		idSolarNoon                = "SolarNoon"
		idSetStart                 = "SetStart"
		idSetEnd                   = "SetEnd"
		idSetDuration              = "SetDuration"
		idCivilDuskStart           = "CivilDuskStart"
		idCivilDuskEnd             = "CivilDuskEnd"
		idCivilDuskDuration        = "CivilDuskDuration"
		idNauticalDuskStart        = "NauticalDuskStart"
		idNauticalDuskEnd          = "NauticalDuskEnd"
		idNauticalDuskDuration     = "NauticalDuskDuration"
		idAstronomicalDuskStart    = "AstronomicalDuskStart"
		idAstronomicalDuskEnd      = "AstronomicalDuskEnd"
		idAstronomicalDuskDuration = "AstronomicalDuskDuration"
		idNightAfterStart          = "NightAfterStart"
		idNightAfterEnd            = "NightAfterEnd"
		idNightAfterDuration       = "NightAfterDuration"
	)

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel(idNadir, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNadir.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNadir, openhab.ItemTypeDateTime).
					WithLabel("Nadir [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel(idNightBeforeStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNightBeforeStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNightBeforeStart, openhab.ItemTypeDateTime).
					WithLabel("Night before start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel(idNightBeforeEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNightBeforeEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNightBeforeEnd, openhab.ItemTypeDateTime).
					WithLabel("Night before end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel(idNightBeforeDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicNightBeforeDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNightBeforeDuration, openhab.ItemTypeNumber).
					WithLabel("Night before duration [%d s]").
					WithIcon("moon"),
			),
		openhab.NewChannel(idAstronomicalDawnStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicAstronomicalDawnStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idAstronomicalDawnStart, openhab.ItemTypeDateTime).
					WithLabel("Astronomical dawn start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idAstronomicalDawnEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicAstronomicalDawnEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idAstronomicalDawnEnd, openhab.ItemTypeDateTime).
					WithLabel("Astronomical dawn end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idAstronomicalDawnDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicAstronomicalDawnDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idAstronomicalDawnDuration, openhab.ItemTypeNumber).
					WithLabel("Astronomical dawn duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idNauticalDawnStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNauticalDawnStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNauticalDawnStart, openhab.ItemTypeDateTime).
					WithLabel("Nautical dawn start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idNauticalDawnEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNauticalDawnEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNauticalDawnEnd, openhab.ItemTypeDateTime).
					WithLabel("Nautical dawn end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idNauticalDawnDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicNauticalDawnDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNauticalDawnDuration, openhab.ItemTypeNumber).
					WithLabel("Nautical dawn duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idCivilDawnStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicCivilDawnStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idCivilDawnStart, openhab.ItemTypeDateTime).
					WithLabel("Civil dawn start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idCivilDawnEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicCivilDawnEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idCivilDawnEnd, openhab.ItemTypeDateTime).
					WithLabel("Civil dawn end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idCivilDawnDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCivilDawnDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idCivilDawnDuration, openhab.ItemTypeNumber).
					WithLabel("Civil dawn duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idRiseStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicRiseStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idRiseStart, openhab.ItemTypeDateTime).
					WithLabel("Sunrise start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idRiseEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicRiseEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idRiseEnd, openhab.ItemTypeDateTime).
					WithLabel("Sunrise end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idRiseDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicRiseDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idRiseDuration, openhab.ItemTypeNumber).
					WithLabel("Sunrise duration [%d s]").
					WithIcon("sunrise"),
			),
		openhab.NewChannel(idSolarNoon, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicSolarNoon.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idSolarNoon, openhab.ItemTypeDateTime).
					WithLabel("Solar noon [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sun"),
			),
		openhab.NewChannel(idSetStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicSetStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idSetStart, openhab.ItemTypeDateTime).
					WithLabel("Sunset start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idSetEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicSetEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idSetEnd, openhab.ItemTypeDateTime).
					WithLabel("Sunset end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idSetDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicSetDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idSetDuration, openhab.ItemTypeNumber).
					WithLabel("Sunset duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idCivilDuskStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicCivilDuskStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idCivilDuskStart, openhab.ItemTypeDateTime).
					WithLabel("Civil dusk start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idCivilDuskEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicCivilDuskEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idCivilDuskEnd, openhab.ItemTypeDateTime).
					WithLabel("Civil dusk end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idCivilDuskDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicCivilDuskDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idCivilDuskDuration, openhab.ItemTypeNumber).
					WithLabel("Civil dusk duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idNauticalDuskStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNauticalDuskStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNauticalDuskStart, openhab.ItemTypeDateTime).
					WithLabel("Nautical dusk start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idNauticalDuskEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNauticalDuskEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNauticalDuskEnd, openhab.ItemTypeDateTime).
					WithLabel("Nautical dusk end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idNauticalDuskDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicNauticalDuskDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNauticalDuskDuration, openhab.ItemTypeNumber).
					WithLabel("Nautical dusk duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idAstronomicalDuskStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicAstronomicalDuskStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idAstronomicalDuskStart, openhab.ItemTypeDateTime).
					WithLabel("Astronomical dusk start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idAstronomicalDuskEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicAstronomicalDuskEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idAstronomicalDuskEnd, openhab.ItemTypeDateTime).
					WithLabel("Astronomical dusk end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idAstronomicalDuskDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicAstronomicalDuskDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idAstronomicalDuskDuration, openhab.ItemTypeNumber).
					WithLabel("Astronomical dusk duration [%d s]").
					WithIcon("sunset"),
			),
		openhab.NewChannel(idNightAfterStart, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNightAfterStart.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNightAfterStart, openhab.ItemTypeDateTime).
					WithLabel("Night after start [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel(idNightAfterEnd, openhab.ChannelTypeDateTime).
			WithStateTopic(cfg.TopicNightAfterEnd.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNightAfterEnd, openhab.ItemTypeDateTime).
					WithLabel("Night after end [%1$tH:%1$tM:%1$tS %1$tz]").
					WithIcon("moon"),
			),
		openhab.NewChannel(idNightAfterDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicNightAfterDuration.Format(id)).
			AddItems(
				openhab.NewItem(itemPrefix+idNightAfterDuration, openhab.ItemTypeNumber).
					WithLabel("Night after duration [%d s]").
					WithIcon("moon"),
			),
	)
}
