package openweathermap

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/openweathermap"
	"github.com/kihamo/boggart/providers/openweathermap/client/weather"
	"github.com/kihamo/boggart/providers/openweathermap/models"
	"github.com/kihamo/shadow/components/i18n"
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.WidgetBind

	config *Config
	client *openweathermap.Client
}

func (b *Bind) lang(ctx context.Context) *string {
	var lang string
	locale := i18n.Locale(ctx).Locale()

	switch locale {
	case "ru":
		lang = locale
	}

	if lang != "" {
		return &lang
	}

	return nil
}

func (b *Bind) Current(ctx context.Context) (current *models.Current, err error) {
	switch {
	case b.config.CityID > 0:
		var response *weather.GetCurrentByCityIDOK

		params := weather.NewGetCurrentByCityIDParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&b.config.Units).
			WithID(b.config.CityID)

		response, err = b.client.Weather.GetCurrentByCityID(params, nil)
		if err == nil {
			current = response.Payload
		}
	case b.config.CityName != "":
		var response *weather.GetCurrentByCityNameOK

		params := weather.NewGetCurrentByCityNameParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&b.config.Units).
			WithQ(b.config.CityName)

		response, err = b.client.Weather.GetCurrentByCityName(params, nil)
		if err == nil {
			current = response.Payload
		}
	case b.config.Latitude != 0 && b.config.Longitude != 0:
		var response *weather.GetCurrentByGeographicCoordinatesOK

		params := weather.NewGetCurrentByGeographicCoordinatesParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&b.config.Units).
			WithLat(b.config.Latitude).
			WithLon(b.config.Longitude)

		response, err = b.client.Weather.GetCurrentByGeographicCoordinates(params, nil)
		if err == nil {
			current = response.Payload
		}
	case b.config.Zip != "":
		var response *weather.GetCurrentByZIPCodeOK

		params := weather.NewGetCurrentByZIPCodeParamsWithContext(ctx).
			WithLang(b.lang(ctx)).
			WithUnits(&b.config.Units).
			WithZip(b.config.Zip)

		response, err = b.client.Weather.GetCurrentByZIPCode(params, nil)
		if err == nil {
			current = response.Payload
		}
	default:
		err = errors.New("location is empty")
	}

	return current, err
}
