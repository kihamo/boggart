package openweathermap

import (
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/providers/openweathermap"
	"github.com/kihamo/boggart/providers/openweathermap/models"
	"github.com/kihamo/shadow/components/dashboard"
)

type ForecastView struct {
	Min *models.ForecastListItem
	Max *models.ForecastListItem
}

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	ctx := r.Context()

	current, err := b.Current(ctx)
	if err != nil {
		widget.FlashError(r, "Get current weather failed with error %v", "", err)
	}

	forecast, err := b.Forecast(ctx)
	if err != nil {
		widget.FlashError(r, "Get forecast failed with error %v", "", err)
	}

	now := time.Now()
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	forecastDays := make(map[time.Time]ForecastView, 5)

	for _, item := range forecast.List {
		dt := time.Unix(int64(item.Dt), 0)
		if dt.Before(todayEnd) {
			continue
		}

		var f ForecastView

		key := time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, now.Location())
		if _, ok := forecastDays[key]; ok {
			f = forecastDays[key]
		}

		if f.Min == nil || f.Min.Main.Temp > item.Main.Temp {
			f.Min = item
		}

		if f.Max == nil || f.Max.Main.Temp < item.Main.Temp {
			f.Max = item
		}

		forecastDays[key] = f
	}

	widget.Render(r.Context(), "widget", map[string]interface{}{
		"current":          current,
		"current_datetime": time.Unix(int64(current.Dt), 0),
		"forecast":         forecastDays,
		"icon":             openweathermap.Icon,
	})
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
