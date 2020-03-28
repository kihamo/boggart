package pulsar

import (
	"sort"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/pulsar"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	config := b.Config().(*Config)
	q := r.URL().Query()
	vars := map[string]interface{}{
		"action": q.Get("action"),
	}

	switch vars["action"] {
	case "archive":
		type stat struct {
			Date time.Time
			Energy, EnergyDelta, EnergyTrend,
			Pulse1, Pulse1Volume, Pulse1Delta, Pulse1Trend,
			Pulse2, Pulse2Volume, Pulse2Delta, Pulse2Trend,
			Pulse3, Pulse3Volume, Pulse3Delta, Pulse3Trend,
			Pulse4, Pulse4Volume, Pulse4Delta, Pulse4Trend float32
		}

		stats := make([]*stat, 0)
		statsByDate := make(map[int]*stat)
		statsKey := make([]int, 0)
		end := time.Now()

		var (
			period pulsar.ArchiveType
			start  time.Time

			date   time.Time
			values []float32
			err    error
		)

		switch q.Get("period") {
		case "daily":
			period = pulsar.ArchiveTypeDaily
			start = end.AddDate(0, -1, 0)
			vars["period"] = "daily"
		case "hourly":
			period = pulsar.ArchiveTypeHourly
			start = end.AddDate(0, 0, -1)
			vars["period"] = "hourly"
		default:
			period = pulsar.ArchiveTypeMonthly
			start = end.AddDate(-1, 0, 0)
			vars["period"] = "monthly"
		}

		if queryTime := q.Get("from"); queryTime != "" {
			if tm, err := time.Parse(time.RFC3339, queryTime); err == nil {
				start = tm
			} else {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Parse date from failed with error %s", "", err.Error()))
			}
		}

		if queryTime := q.Get("to"); queryTime != "" {
			if tm, err := time.Parse(time.RFC3339, queryTime); err == nil {
				end = tm
			} else {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Parse date to failed with error %s", "", err.Error()))
			}
		}

		vars["date_from"] = start
		vars["date_to"] = end

		// energy
		date, values, err = bind.provider.EnergyArchive(start, end, period)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get energy archive failed with error %s", "", err.Error()))
		} else {
			for _, value := range values {
				key := int(date.Unix())
				statsKey = append(statsKey, key)
				statsByDate[key] = &stat{
					Date:   date,
					Energy: value,
				}

				date = nextData(period, date)
			}
		}

		// pulse input 1
		date, values, err = bind.provider.PulseInput1Archive(start, end, period)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get pulse %d archive failed with error %s", "", 1, err.Error()))
		} else {
			for _, value := range values {
				row, ok := statsByDate[int(date.Unix())]
				if ok {
					row.Pulse1 = value
					row.Pulse1Volume = bind.inputVolume(value, config.Input1Offset)
				}

				date = nextData(period, date)
			}
		}

		// pulse input 2
		date, values, err = bind.provider.PulseInput2Archive(start, end, period)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get pulse %d archive failed with error %s", "", 2, err.Error()))
		} else {
			for _, value := range values {
				row, ok := statsByDate[int(date.Unix())]
				if ok {
					row.Pulse2 = value
					row.Pulse2Volume = bind.inputVolume(value, config.Input2Offset)
				}

				date = nextData(period, date)
			}
		}

		// pulse input 3
		date, values, err = bind.provider.PulseInput3Archive(start, end, period)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get pulse %d archive failed with error %s", "", 3, err.Error()))
		} else {
			for _, value := range values {
				row, ok := statsByDate[int(date.Unix())]
				if ok {
					row.Pulse3 = value
					row.Pulse3Volume = bind.inputVolume(value, config.Input3Offset)
				}

				date = nextData(period, date)
			}
		}

		// pulse input 4
		date, values, err = bind.provider.PulseInput4Archive(start, end, period)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get pulse %d archive failed with error %s", "", 4, err.Error()))
		} else {
			for _, value := range values {
				row, ok := statsByDate[int(date.Unix())]
				if ok {
					row.Pulse4 = value
					row.Pulse4Volume = bind.inputVolume(value, config.Input4Offset)
				}

				date = nextData(period, date)
			}
		}

		sort.Ints(statsKey)

		for _, k := range statsKey {
			stats = append(stats, statsByDate[k])
		}

		// deltas
		for i, current := range stats {
			if i == 0 {
				continue
			}

			current.EnergyDelta = current.Energy - stats[i-1].Energy
			current.Pulse1Delta = current.Pulse1Volume - stats[i-1].Pulse1Volume
			current.Pulse2Delta = current.Pulse2Volume - stats[i-1].Pulse2Volume
			current.Pulse3Delta = current.Pulse3Volume - stats[i-1].Pulse3Volume
			current.Pulse4Delta = current.Pulse4Volume - stats[i-1].Pulse4Volume
		}

		// trends
		for i, current := range stats {
			if i < 2 {
				continue
			}

			current.EnergyTrend = current.EnergyDelta - stats[i-1].EnergyDelta
			current.Pulse1Trend = current.Pulse1Delta - stats[i-1].Pulse1Delta
			current.Pulse2Trend = current.Pulse2Delta - stats[i-1].Pulse2Delta
			current.Pulse3Trend = current.Pulse3Delta - stats[i-1].Pulse3Delta
			current.Pulse4Trend = current.Pulse4Delta - stats[i-1].Pulse4Delta
		}

		vars["stats"] = statsByDate

	default:
		type metricView struct {
			Value interface{}
			Delta int64
			Error error
		}

		// date time
		now := time.Now()
		timeValue, err := bind.provider.Datetime()
		variable := metricView{
			Value: timeValue,
			Error: err,
		}

		if err == nil {
			variable.Delta = int64(now.Sub(timeValue).Seconds())
		}

		vars["datetime"] = variable

		floatValue, err := bind.provider.TemperatureIn()
		vars["temperature_in"] = metricView{
			Value: floatValue,
			Error: err,
		}

		floatValue, err = bind.provider.TemperatureOut()
		vars["temperature_out"] = metricView{
			Value: floatValue,
			Error: err,
		}

		floatValue, err = bind.provider.TemperatureDelta()
		vars["temperature_delta"] = metricView{
			Value: floatValue,
			Error: err,
		}

		floatValue, err = bind.provider.Power()
		vars["power"] = metricView{
			Value: floatValue,
			Error: err,
		}

		floatValue, err = bind.provider.Energy()
		vars["energy"] = metricView{
			Value: floatValue,
			Error: err,
		}

		floatValue, err = bind.provider.Capacity()
		vars["capacity"] = metricView{
			Value: floatValue,
			Error: err,
		}

		floatValue, err = bind.provider.Consumption()
		vars["consumption"] = metricView{
			Value: floatValue,
			Error: err,
		}

		floatValue, err = bind.provider.PulseInput1()
		vars["pusle_input_1"] = metricView{
			Value: floatValue,
			Error: err,
		}

		floatValue, err = bind.provider.PulseInput2()
		vars["pusle_input_2"] = metricView{
			Value: floatValue,
			Error: err,
		}

		floatValue, err = bind.provider.PulseInput3()
		vars["pusle_input_3"] = metricView{
			Value: floatValue,
			Error: err,
		}

		floatValue, err = bind.provider.PulseInput4()
		vars["pusle_input_4"] = metricView{
			Value: floatValue,
			Error: err,
		}

		durationValue, err := bind.provider.OperatingTime()
		vars["operating_time"] = metricView{
			Value: durationValue,
			Error: err,
		}
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func nextData(period pulsar.ArchiveType, date time.Time) time.Time {
	switch period {
	case pulsar.ArchiveTypeMonthly:
		date = date.AddDate(0, 1, 0)

	case pulsar.ArchiveTypeDaily:
		date = date.AddDate(0, 0, 1)

	case pulsar.ArchiveTypeHourly:
		date = date.Add(time.Hour)
	}

	return date
}
