package pulsar

import (
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	vars := map[string]interface{}{}

	action := r.URL().Query().Get("action")
	vars["action"] = action

	switch action {

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
