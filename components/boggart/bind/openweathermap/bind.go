package openweathermap

import (
	"context"
	"errors"
	"sync"
	"time"

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

	locationMutex sync.RWMutex
	locationName  string
	locationCoord *models.Coord
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	b.client = openweathermap.New(cfg.APIKey, cfg.Price, cfg.Debug, swagger.NewLogger(
		func(message string) {
			b.Logger().Info(message)
		},
		func(message string) {
			b.Logger().Debug(message)
		}))

	return nil
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
	b.locationMutex.RLock()
	defer b.locationMutex.RUnlock()

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

/*
You can search weather forecast for 5 days with data every 3 hours by geographic coordinates

Пока еще доступно на бесплатном тарифе
*/
func (b *Bind) Forecast(ctx context.Context) (current *models.Forecast, err error) {
	b.locationMutex.RLock()
	defer b.locationMutex.RUnlock()

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

/*
Уже не доступно на бесплатном тарифе
*/
func (b *Bind) OneCall(ctx context.Context, include []string) (*models.OneCall, error) {
	b.locationMutex.RLock()
	defer b.locationMutex.RUnlock()

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

func (b *Bind) OneCallMigrate(ctx context.Context) (*models.OneCall, error) {
	responseCurrent, err := b.Current(ctx)
	if err != nil {
		return nil, err
	}

	responseForecast, err := b.Forecast(ctx)
	if err != nil {
		return nil, err
	}

	var (
		hourlyItem   *models.OneCallHourlyItems0
		daylyItem    *models.OneCallDailyItems0
		forecastItem *models.ForecastListItem

		dayTime, dayItem time.Time
		dayIndex         int
	)

	now := time.Now()
	dayNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	result := &models.OneCall{
		Current: &models.OneCallCurrent{
			// Clouds: responseCurrent.Clouds.All,
			// DewPoint:
			Dt: responseCurrent.Dt,
			//FeelsLike: responseCurrent.Main.FeelsLike,
			//Humidity: responseCurrent.Main.Humidity,
			// Pop:
			//Pressure: responseCurrent.Main.Pressure,
			Rain: responseCurrent.Rain,
			Snow: responseCurrent.Snow,
			// Sunrise: responseCurrent.Sys.Sunrise,
			// Sunset:  responseCurrent.Sys.Sunset,
			// Temp:    responseCurrent.Main.Temp,
			// Uvi:
			Visibility: responseCurrent.Visibility,
			Weather:    responseCurrent.Weather,
		},
		Daily:  make([]*models.OneCallDailyItems0, 0, 5),
		Hourly: make([]*models.OneCallHourlyItems0, 0, 40),
		// Lat: responseCurrent.Coord.Lat,
		// Lon: responseCurrent.Coord.Lon,
		// Minutely:
		// Timezone:
		// TimezoneOffset: responseForecast.City.Timezone,
	}

	if responseCurrent.Clouds != nil {
		result.Current.Clouds = responseCurrent.Clouds.All
	}

	if responseCurrent.Main != nil {
		result.Current.FeelsLike = responseCurrent.Main.FeelsLike
		result.Current.Humidity = responseCurrent.Main.Humidity
		result.Current.Pressure = responseCurrent.Main.Pressure
		result.Current.Temp = responseCurrent.Main.Temp
	}

	if responseCurrent.Sys != nil {
		result.Current.Sunrise = responseCurrent.Sys.Sunrise
		result.Current.Sunset = responseCurrent.Sys.Sunset
	}

	if responseCurrent.Wind != nil {
		result.Current.WindGust = responseCurrent.Wind.Gust
		result.Current.WindSpeed = responseCurrent.Wind.Speed
	}

	if responseCurrent.Coord != nil {
		result.Lat = responseCurrent.Coord.Lat
		result.Lon = responseCurrent.Coord.Lon
	}

	if responseForecast.City != nil {
		result.TimezoneOffset = responseForecast.City.Timezone
	}

	for _, forecastItem = range responseForecast.List {
		if forecastItem.Main == nil {
			// TODO: error
			continue
		}

		dayTime = forecastItem.Dt.Time()
		dayItem = time.Date(dayTime.Year(), dayTime.Month(), dayTime.Day(), 0, 0, 0, 0, dayTime.Location())
		dayIndex = int(dayItem.Sub(dayNow) / (time.Hour * 24))

		hourlyItem = &models.OneCallHourlyItems0{
			Dt:         forecastItem.Dt,
			Rain:       forecastItem.Rain,
			Snow:       forecastItem.Snow,
			Visibility: forecastItem.Visibility,
			Weather:    forecastItem.Weather,
		}

		if len(result.Daily) == dayIndex+1 {
			daylyItem = result.Daily[dayIndex]
		} else {
			daylyItem = &models.OneCallDailyItems0{
				FeelsLike: &models.OneCallDailyItems0FeelsLike{},
				Temp: &models.OneCallDailyItems0Temp{
					Min: forecastItem.Main.TempMin,
					Max: forecastItem.Main.TempMax,
				},
			}
			result.Daily = append(result.Daily, daylyItem)
		}

		if forecastItem.Clouds != nil {
			hourlyItem.Clouds = forecastItem.Clouds.All
		}

		if forecastItem.Clouds != nil {
			hourlyItem.FeelsLike = forecastItem.Main.FeelsLike
			hourlyItem.Humidity = forecastItem.Main.Humidity
			hourlyItem.Pressure = forecastItem.Main.Pressure
			hourlyItem.Temp = forecastItem.Main.Temp
		}

		if forecastItem.Wind != nil {
			hourlyItem.WindGust = forecastItem.Wind.Gust
			hourlyItem.WindSpeed = forecastItem.Wind.Speed
		}

		/*
			day - Температура в 12:00 по местному времени
			night - Температура в 00:00 по местному времени
			eve - Температура в 18:00 по местному времени
			morn - Температура в 06:00 по местному времени
		*/

		switch h := forecastItem.Dt.Time().Hour(); h {
		case 0:
			daylyItem.FeelsLike.Night = forecastItem.Main.FeelsLike
			daylyItem.Temp.Night = forecastItem.Main.Temp

		case 6:
			daylyItem.FeelsLike.Morn = forecastItem.Main.FeelsLike
			daylyItem.Temp.Morn = forecastItem.Main.Temp

		case 12:
			daylyItem.Dt = forecastItem.Dt
			daylyItem.FeelsLike.Day = forecastItem.Main.FeelsLike
			daylyItem.Humidity = forecastItem.Main.Humidity
			daylyItem.Pop = forecastItem.Pop
			daylyItem.Pressure = forecastItem.Main.Pressure
			daylyItem.Sunrise = forecastItem.Sunrise
			daylyItem.Sunset = forecastItem.Sunset
			daylyItem.Temp.Day = forecastItem.Main.Temp
			daylyItem.Visibility = forecastItem.Visibility
			daylyItem.Weather = forecastItem.Weather

			if forecastItem.Clouds != nil {
				daylyItem.Clouds = forecastItem.Clouds.All
			}

			if forecastItem.Rain != nil {
				daylyItem.Rain = forecastItem.Rain.Nr1h
			}

			if forecastItem.Snow != nil {
				daylyItem.Snow = forecastItem.Snow.Nr1h
			}

			if forecastItem.Wind != nil {
				daylyItem.WindDeg = forecastItem.Wind.Deg
				daylyItem.WindGust = forecastItem.Wind.Gust
				daylyItem.WindSpeed = forecastItem.Wind.Speed
			}

		case 18:
			daylyItem.FeelsLike.Eve = forecastItem.Main.FeelsLike
			daylyItem.Temp.Eve = forecastItem.Main.Temp
		}

		if forecastItem.Main.TempMax > daylyItem.Temp.Max {
			daylyItem.Temp.Max = forecastItem.Main.TempMax
		}

		if forecastItem.Main.TempMin < daylyItem.Temp.Min {
			daylyItem.Temp.Min = forecastItem.Main.TempMin
		}

		result.Hourly = append(result.Hourly, hourlyItem)

		// убираем последний день если он не полный (не описан до полудня включительно)
		if len(result.Daily) > 0 && forecastItem != nil && forecastItem.Dt.Time().Hour() < 12 {
			result.Daily = result.Daily[:len(result.Daily)-1]
		}
	}

	return result, nil
}
