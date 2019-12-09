package sun

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers/task"
	"github.com/mourner/suncalc-go"
)

const (
	dayDuration = time.Hour * 24
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config           *Config
	taskStateUpdater *task.FunctionTask
}

type Time struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

type Times struct {
	NightBefore      Time
	NightAfter       Time
	Nadir            time.Time
	AstronomicalDawn Time
	NauticalDawn     Time
	CivilDawn        Time
	Sunrise          Time
	SolarNoon        time.Time
	Sunset           Time
	CivilDusk        Time
	NauticalDusk     Time
	AstronomicalDusk Time
}

func (b *Bind) Times() Times {
	t := Times{}

	now := time.Now()

	// для 00:00:00 почему-то считает предыдущий день, поэтому берем полдень
	todaySolarNoon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	timesYesterday := suncalc.SunTimes(todaySolarNoon.Add(-dayDuration), b.config.Lat, b.config.Lon)
	timesToday := suncalc.SunTimes(todaySolarNoon, b.config.Lat, b.config.Lon)
	timesTomorrow := suncalc.SunTimes(todaySolarNoon.Add(dayDuration), b.config.Lat, b.config.Lon)

	t.NightBefore.Start = timesYesterday["night"]
	t.NightBefore.End = timesToday["nightEnd"]
	t.NightBefore.Duration = t.NightBefore.End.Sub(t.NightBefore.Start)

	t.NightAfter.Start = timesToday["night"]
	t.NightAfter.End = timesTomorrow["nightEnd"]
	t.NightAfter.Duration = t.NightAfter.End.Sub(t.NightAfter.Start)

	t.Nadir = timesToday["nadir"]

	t.AstronomicalDawn.Start = timesToday["nightEnd"]
	t.AstronomicalDawn.End = timesToday["nauticalDawn"]
	t.AstronomicalDawn.Duration = t.AstronomicalDawn.End.Sub(t.AstronomicalDawn.Start)

	t.NauticalDawn.Start = timesToday["nauticalDawn"]
	t.NauticalDawn.End = timesToday["dawn"]
	t.NauticalDawn.Duration = t.NauticalDawn.End.Sub(t.NauticalDawn.Start)

	t.CivilDawn.Start = timesToday["dawn"]
	t.CivilDawn.End = timesToday["sunrise"]
	t.CivilDawn.Duration = t.CivilDawn.End.Sub(t.CivilDawn.Start)

	t.Sunrise.Start = timesToday["sunrise"]
	t.Sunrise.End = timesToday["sunriseEnd"]
	t.Sunrise.Duration = t.Sunrise.End.Sub(t.Sunrise.Start)

	t.SolarNoon = timesToday["solarNoon"]

	t.Sunset.Start = timesToday["sunsetStart"]
	t.Sunset.End = timesToday["sunset"]
	t.Sunset.Duration = t.Sunset.End.Sub(t.Sunset.Start)

	t.CivilDusk.Start = timesToday["sunset"]
	t.CivilDusk.End = timesToday["dusk"]
	t.CivilDusk.Duration = t.CivilDusk.End.Sub(t.CivilDusk.Start)

	t.NauticalDusk.Start = timesToday["dusk"]
	t.NauticalDusk.End = timesToday["nauticalDusk"]
	t.NauticalDusk.Duration = t.NauticalDusk.End.Sub(t.NauticalDusk.Start)

	t.AstronomicalDusk.Start = timesToday["nauticalDusk"]
	t.AstronomicalDusk.End = timesToday["night"]
	t.AstronomicalDusk.Duration = t.AstronomicalDusk.End.Sub(t.AstronomicalDusk.Start)

	return t
}
