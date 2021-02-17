package owntracks

import (
	"context"
	"strconv"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/mmcloughlin/geohash"
	"go.uber.org/multierr"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()

	subscribers := []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicCommandReportLocation.Format(cfg.User, cfg.Device), 0, b.subscribeCommand(b.CommandReportLocation)),
		mqtt.NewSubscriber(cfg.TopicCommandRestart.Format(cfg.User, cfg.Device), 0, b.subscribeCommand(b.CommandRestart)),
		mqtt.NewSubscriber(cfg.TopicCommandReconnect.Format(cfg.User, cfg.Device), 0, b.subscribeCommand(b.CommandReconnect)),
		mqtt.NewSubscriber(cfg.TopicCommandWayPoints.Format(cfg.User, cfg.Device), 0, b.subscribeCommand(b.CommandWayPoints)),
		mqtt.NewSubscriber(cfg.TopicOwnTracksUserLocation.Format(cfg.User, cfg.Device), 0, b.subscribeUserLocation),
		mqtt.NewSubscriber(cfg.TopicOwnTracksTransition.Format(cfg.User, cfg.Device), 0, b.subscribeTransition),
	}

	if cfg.RegionsSyncEnabled {
		subscribers = append(
			subscribers,
			mqtt.NewSubscriber(cfg.TopicOwnTracksWayPoints.Format(cfg.User, cfg.Device), 0, b.subscribeSyncRegions),
		)
	}

	return subscribers
}

func (b *Bind) notifyRegionEvent(ctx context.Context, payload *LocationPayload, q byte, r bool) (err error) {
	alreadyNotify := make(map[string]struct{})
	cfg := b.config()

	if cfg.CheckInRegionEnabled && payload.InRegions != nil {
		for _, name := range *payload.InRegions {
			alreadyNotify[name] = struct{}{}

			var checker *atomic.BoolNull

			if cfg.UnregisterPointsAllowed {
				checker = b.getOrSetRegionChecker(name)
			} else {
				checker, _ = b.getRegionChecker(name)
			}

			if checker != nil {
				if ok := checker.Set(true); ok {
					if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicRegion.Format(cfg.User, cfg.Device, name), q, r, true); e != nil {
						err = multierr.Append(err, e)
					}
				}
			}
		}
	}

	if cfg.CheckDistanceEnabled && b.validAccuracy(payload.Acc, cfg.MaxAccuracy) {
		for name, point := range b.getAllRegions() {
			// если точно больше радиуса, то вероятность ошибки большая, пропускаем
			if !b.validAccuracy(payload.Acc, int64(point.Radius)) {
				continue
			}

			if _, ok := alreadyNotify[name]; ok {
				continue
			}

			checker, ok := b.getRegionChecker(name)
			if !ok {
				continue
			}

			alreadyNotify[name] = struct{}{}

			checkResult := calculateDistance(*payload.Lat, *payload.Lon, point.Lat, point.Lon) < point.Radius

			if ok := checker.Set(checkResult); ok {
				if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicRegion.Format(cfg.User, cfg.Device, name), q, r, checkResult); e != nil {
					err = multierr.Append(err, e)
				}
			}
		}
	}

	for name, checker := range b.getAllRegionCheckers() {
		if _, ok := alreadyNotify[name]; ok {
			continue
		}

		if ok := checker.Set(false); ok {
			if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicRegion.Format(cfg.User, cfg.Device, name), q, r, false); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	return err
}

func (b *Bind) subscribeUserLocation(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
	// skip lwt
	var payloadLWT *LWTPayload
	if err = message.JSONUnmarshal(&payloadLWT); err == nil {
		if err = payloadLWT.Valid(); err == nil {
			return nil
		}
	}

	var payload *LocationPayload
	if err = message.JSONUnmarshal(&payload); err != nil {
		return err
	}

	if err = payload.Valid(); err != nil {
		return err
	}

	var changeLocation bool
	cfg := b.config()
	q := message.Qos()
	r := message.Retained()

	if ok := b.lat.Set(float32(*payload.Lat)); ok {
		changeLocation = true

		if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicStateLat.Format(cfg.User, cfg.Device), q, r, *payload.Lat); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if ok := b.lon.Set(float32(*payload.Lon)); ok {
		changeLocation = true

		if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicStateLon.Format(cfg.User, cfg.Device), q, r, *payload.Lon); e != nil {
			err = multierr.Append(err, e)
		}
	}

	hash := geohash.Encode(*payload.Lat, *payload.Lon)
	if ok := b.geoHash.Set(hash); ok {
		if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicStateGeoHash.Format(cfg.User, cfg.Device), q, r, hash); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if changeLocation {
		location := strconv.FormatFloat(*payload.Lat, 'f', -1, 64) + "," + strconv.FormatFloat(*payload.Lon, 'f', -1, 64)
		if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicStateLocation.Format(cfg.User, cfg.Device), q, r, location); e != nil {
			err = multierr.Append(err, e)
		}
	}

	// detect event
	if e := b.notifyRegionEvent(ctx, payload, message.Qos(), message.Retained()); e != nil {
		err = multierr.Append(err, e)
	}

	if payload.Acc != nil {
		if ok := b.acc.Set(int32(*payload.Acc)); ok {
			if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicStateAccuracy.Format(cfg.User, cfg.Device), q, r, *payload.Acc); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	if payload.Alt != nil {
		if ok := b.alt.Set(int32(*payload.Alt)); ok {
			if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicStateAltitude.Format(cfg.User, cfg.Device), q, r, *payload.Alt); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	if payload.Batt != nil {
		metricBatteryLevel.With("user", cfg.User, "device", cfg.Device).Set(*payload.Batt)

		if ok := b.batt.Set(float32(*payload.Batt)); ok {
			if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicStateBatteryLevel.Format(cfg.User, cfg.Device), q, r, *payload.Batt); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	if payload.Vel != nil {
		metricVelocity.With("user", cfg.User, "device", cfg.Device).Set(float64(*payload.Vel))

		if ok := b.vel.Set(int32(*payload.Vel)); ok {
			if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicStateVelocity.Format(cfg.User, cfg.Device), q, r, *payload.Vel); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	if payload.Conn != nil {
		var v string

		switch *payload.Conn {
		case "w":
			v = "wifi"
		case "o":
			v = "offline"
		case "m":
			v = "mobile"
		default:
			v = "unknown"
		}

		if ok := b.conn.Set(v); ok {
			if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicStateConnection.Format(cfg.User, cfg.Device), q, r, v); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	return err
}

// обработка единичного события попадания/уход из региона.
// !!! Это событие не означает что в других регионах ситуация поменялась
func (b *Bind) subscribeTransition(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
	var payload *TransitionPayload
	if err = message.JSONUnmarshal(&payload); err != nil {
		return err
	}

	if err = payload.Valid(); err != nil {
		return err
	}

	var check *atomic.BoolNull
	cfg := b.config()

	if cfg.UnregisterPointsAllowed {
		check = b.getOrSetRegionChecker(payload.Desc)
	} else {
		check, _ = b.getRegionChecker(payload.Desc)
	}

	if check != nil {
		checkResult := payload.IsEnter()

		if ok := check.Set(checkResult); ok {
			if e := b.MQTT().PublishAsyncRaw(ctx, cfg.TopicRegion.Format(cfg.User, cfg.Device, payload.Desc), message.Qos(), message.Retained(), checkResult); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	return err
}

func (b *Bind) subscribeSyncRegions(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
	// _type == waypoint
	// одиночное добавление, вручную внесли новый пункт в список (или результат синка уже)
	// такие добавления игнорируем, тикет все равно дернет запрос всего списка
	var payloadOne *WayPointPayload
	if err = message.JSONUnmarshal(&payloadOne); err == nil {
		if err = payloadOne.Valid(); err == nil {
			return nil
		}
	}

	// _type == waypoints
	var payload *WayPointsPayload
	if err = message.JSONUnmarshal(&payload); err != nil {
		return err
	}

	if err = payload.Valid(); err != nil {
		return err
	}

	existsDesc := make(map[string]struct{})
	existsTst := make(map[int64]struct{})

	for _, point := range payload.WayPoints {
		existsTst[point.Tst] = struct{}{}
		existsDesc[point.Desc] = struct{}{}

		// регистрируем новый регион, который не определен в конфиге
		if b.config().UnregisterPointsAllowed {
			b.registerRegion(point.Desc, Point{
				Lat:    point.Lat,
				Lon:    point.Lon,
				Radius: point.Rad,
			})
		}
	}

	existsRegions := b.getAllRegions()
	newRegions := make([]WayPointPayload, 0, len(existsRegions))
	lastTst := time.Now().Unix()

	for name, point := range existsRegions {
		// точка зарегистрирована и у нас и на девайсе, пропускаем
		if _, ok := existsDesc[name]; ok {
			continue
		}

		// генерируем tst
		for {
			lastTst++

			if _, ok := existsTst[lastTst]; !ok {
				existsTst[lastTst] = struct{}{}
				break
			}
		}

		newRegions = append(newRegions, WayPointPayload{
			Desc: name,
			Lat:  point.Lat,
			Lon:  point.Lon,
			Rad:  point.Radius,
			Tst:  lastTst,
		})
	}

	if len(newRegions) == 0 {
		return nil
	}

	return b.CommandSetWayPoints(newRegions)
}

func (b *Bind) subscribeCommand(cmd func() error) mqtt.MessageHandler {
	return func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
		return cmd()
	}
}
