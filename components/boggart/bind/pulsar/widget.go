package pulsar

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/pulsar"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	vars := map[string]interface{}{
		"action": r.URL().Query().Get("action"),
	}

	switch vars["action"] {
	case "archive":
		type stat struct {
			Date   time.Time
			Energy float32
		}

		stats := make([]stat, 0)
		end := time.Now()

		var (
			period pulsar.ArchiveType
			start  time.Time
		)

		switch r.URL().Query().Get("period") {
		case "daily":
			period = pulsar.ArchiveTypeDaily
			start = end.AddDate(0, -1, 0)
		case "hourly":
			period = pulsar.ArchiveTypeHourly
			start = end.AddDate(0, 0, -1)
		default:
			period = pulsar.ArchiveTypeMonthly
			start = end.AddDate(-1, 0, 0)
		}

		// energy
		date, values, err := bind.provider.EnergyArchive(start, end, period)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(r.Context(), "Get archive failed with error %s", "", err.Error()))
		} else {
			for _, value := range values {
				stats = append(stats, stat{
					Date:   date,
					Energy: value,
				})

				switch period {
				case pulsar.ArchiveTypeMonthly:
					date = date.AddDate(0, 1, 0)

				case pulsar.ArchiveTypeDaily:
					date = date.AddDate(0, 0, 1)

				case pulsar.ArchiveTypeHourly:
					date = date.Add(time.Hour)
				}
			}
		}

		// pulse input 1

		// pulse input 2

		// pulse input 3

		// pulse input 4

		vars["stats"] = stats

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
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
