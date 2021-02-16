package openweathermap

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/openweathermap"
	"github.com/kihamo/boggart/providers/openweathermap/client/forecast"
	"github.com/kihamo/boggart/providers/openweathermap/client/onecall"
	"github.com/kihamo/boggart/providers/openweathermap/client/weather"
	"github.com/kihamo/boggart/providers/openweathermap/models"
	"github.com/kihamo/shadow/components/i18n"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.WidgetBind
	di.WorkersBind

	client *openweathermap.Client

	locationName  string
	locationCoord *models.Coord
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() (err error) {
	ctx := context.Background()
	cfg := b.config()

	b.client = openweathermap.New(cfg.APIKey, cfg.Price, cfg.Debug, swagger.NewLogger(
		func(message string) {
			b.Logger().Info(message)
		},
		func(message string) {
			b.Logger().Debug(message)
		}))

	switch {
	case cfg.CityID > 0:
		var response *weather.GetCurrentByCityIDOK

		params := weather.NewGetCurrentByCityIDParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&cfg.Units).
			WithID(cfg.CityID)

		response, err = b.client.Weather.GetCurrentByCityID(params, nil)
		if err == nil {
			b.locationName = response.Payload.Name
			b.locationCoord = response.Payload.Coord
		}
	case cfg.CityName != "":
		var response *weather.GetCurrentByCityNameOK

		params := weather.NewGetCurrentByCityNameParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&cfg.Units).
			WithQ(cfg.CityName)

		response, err = b.client.Weather.GetCurrentByCityName(params, nil)
		if err == nil {
			b.locationName = response.Payload.Name
			b.locationCoord = response.Payload.Coord
		}
	case cfg.Latitude != 0 && cfg.Longitude != 0:
		var response *weather.GetCurrentByGeographicCoordinatesOK

		params := weather.NewGetCurrentByGeographicCoordinatesParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&cfg.Units).
			WithLat(cfg.Latitude).
			WithLon(cfg.Longitude)

		response, err = b.client.Weather.GetCurrentByGeographicCoordinates(params, nil)
		if err == nil {
			b.locationName = response.Payload.Name
			b.locationCoord = response.Payload.Coord
		}
	case cfg.Zip != "":
		var response *weather.GetCurrentByZIPCodeOK

		params := weather.NewGetCurrentByZIPCodeParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&cfg.Units).
			WithZip(cfg.Zip)

		response, err = b.client.Weather.GetCurrentByZIPCode(params, nil)
		if err == nil {
			b.locationName = response.Payload.Name
			b.locationCoord = response.Payload.Coord
		}
	default:
		err = errors.New("location is empty")
	}

	return err
}

func (b *Bind) lang(ctx context.Context) *string {
	var lang string

	if locale := i18n.Locale(ctx).Locale(); locale == "ru" {
		lang = locale
	}

	if lang != "" {
		return &lang
	}

	return nil
}

func (b *Bind) Current(ctx context.Context) (current *models.Current, err error) {
	if b.locationCoord == nil {
		return nil, errors.New("location is empty")
	}

	params := weather.NewGetCurrentByGeographicCoordinatesParamsWithContext(ctx).
		WithLang(b.lang(ctx)).
		WithUnits(&b.config().Units).
		WithLat(b.locationCoord.Lat).
		WithLon(b.locationCoord.Lon)

	response, err := b.client.Weather.GetCurrentByGeographicCoordinates(params, nil)
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

func (b *Bind) Forecast(ctx context.Context) (current *models.Forecast, err error) {
	if b.locationCoord == nil {
		return nil, errors.New("location is empty")
	}

	params := forecast.NewGetForecastByGeographicCoordinatesParamsWithContext(ctx).
		WithLang(b.lang(ctx)).
		WithUnits(&b.config().Units).
		WithLat(b.locationCoord.Lat).
		WithLon(b.locationCoord.Lon)

	response, err := b.client.Forecast.GetForecastByGeographicCoordinates(params, nil)
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

func (b *Bind) OneCall(ctx context.Context, include []string) (*models.OneCall, error) {
	if b.locationCoord == nil {
		return nil, errors.New("location is empty")
	}

	params := onecall.NewGetOneCallParamsWithContext(ctx).
		WithLang(b.lang(ctx)).
		WithUnits(&b.config().Units).
		WithLat(b.locationCoord.Lat).
		WithLon(b.locationCoord.Lon)

	if len(include) > 0 {
		excludeMap := map[string]struct{}{
			"lat":             {},
			"lon":             {},
			"timezone":        {},
			"timezone_offset": {},
			"current":         {},
			"minutely":        {},
			"hourly":          {},
			"daily":           {},
		}

		for _, field := range include {
			delete(excludeMap, field)
		}

		exclude := make([]string, len(excludeMap))
		for field := range excludeMap {
			exclude = append(exclude, field)
		}

		params.SetExclude(exclude)
	}

	response, err := b.client.Onecall.GetOneCall(params, nil)
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}
