package owntracks

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/mmcloughlin/geohash"
	"go.uber.org/multierr"
)

/*
https://owntracks.org/booklet/tech/json/

Android
- card
- cmd
- configuration
- encrypted
- location
- lwt
- transition
- waypoint
- waypoints

In MQTT mode the apps publish to:
- owntracks/user/device with _type=location for location updates, and with _type=lwt
- owntracks/user/device/cmd with _type=cmd for remote commands
- owntracks/user/device/event with _type=transition for enter/leave events
- owntracks/user/device/step to report step counter
- owntracks/user/device/beacon for beacon ranging
- owntracks/user/device/dump for config dumps

In MQTT mode apps subscribe to:
- owntracks/user/device/cmd if remote commands are enabled
- owntracks/+/+ for seeing other user's locations, depending on broker ACL
- owntracks/+/+/event for transition messages (enter/leave)
- owntracks/+/+/info for obtaining cards.
*/

const (
	// owntracks
	MQTTOwnTracksPrefix                     mqtt.Topic = "owntracks/+/+"
	MQTTOwnTracksSubscribeTopicUserLocation            = MQTTOwnTracksPrefix
	MQTTOwnTracksSubscribeTopicTransition              = MQTTOwnTracksPrefix + "/event"
	MQTTOwnTracksSubscribeTopicStep                    = MQTTOwnTracksPrefix + "/step"
	MQTTOwnTracksSubscribeTopicBeacon                  = MQTTOwnTracksPrefix + "/beacon"
	MQTTOwnTracksSubscribeTopicDump                    = MQTTOwnTracksPrefix + "/dump"
	MQTTOwnTracksSubscribeTopicWayPoints               = MQTTOwnTracksPrefix + "/waypoint"
	MQTTOwnTracksPublishTopicCommand                   = MQTTOwnTracksPrefix + "/cmd"
	MQTTOwnTracksPublishTopicUserLocation              = MQTTOwnTracksPrefix
	MQTTOwnTracksPublishTopicTransition                = MQTTOwnTracksPrefix + "/event"
	MQTTOwnTracksPublishTopicCard                      = MQTTOwnTracksPrefix + "/info"

	// custom
	MQTTPrefix                            mqtt.Topic = boggart.ComponentName + "/owntracks/+/+/"
	MQTTSubscribeTopicCommand                        = MQTTPrefix + "cmd/+"
	MQTTPublishTopicRegion                           = MQTTPrefix + "event/+"
	MQTTPublishTopicUserStateLat                     = MQTTPrefix + "state/lat"
	MQTTPublishTopicUserStateLon                     = MQTTPrefix + "state/lon"
	MQTTPublishTopicUserStateGeoHash                 = MQTTPrefix + "state/geohash"
	MQTTPublishTopicUserStateAccuracy                = MQTTPrefix + "state/accuracy"
	MQTTPublishTopicUserStateAltitude                = MQTTPrefix + "state/altitude"
	MQTTPublishTopicUserStateBatteryLevel            = MQTTPrefix + "state/battery-level"
	MQTTPublishTopicUserStateVelocity                = MQTTPrefix + "state/velocity"
	MQTTPublishTopicUserStateConnection              = MQTTPrefix + "state/connection"
	MQTTPublishTopicUserStateLocation                = MQTTPrefix + "state/location"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	topics := []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicUserStateLat.Format(b.config.User, b.config.Device)),
		mqtt.Topic(MQTTPublishTopicUserStateLon.Format(b.config.User, b.config.Device)),
		mqtt.Topic(MQTTPublishTopicUserStateGeoHash.Format(b.config.User, b.config.Device)),
		mqtt.Topic(MQTTPublishTopicUserStateAccuracy.Format(b.config.User, b.config.Device)),
		mqtt.Topic(MQTTPublishTopicUserStateAltitude.Format(b.config.User, b.config.Device)),
		mqtt.Topic(MQTTPublishTopicUserStateBatteryLevel.Format(b.config.User, b.config.Device)),
		mqtt.Topic(MQTTPublishTopicUserStateVelocity.Format(b.config.User, b.config.Device)),
		mqtt.Topic(MQTTPublishTopicUserStateConnection.Format(b.config.User, b.config.Device)),
		mqtt.Topic(MQTTPublishTopicUserStateLocation.Format(b.config.User, b.config.Device)),
	}

	if len(b.config.WayPoints) > 0 {
		for name := range b.config.WayPoints {
			topics = append(
				topics,
				mqtt.Topic(MQTTPublishTopicRegion.Format(b.config.User, b.config.Device, name)),
			)
		}
	}

	return topics
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	subscribers := []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicCommand.Format(b.config.User, b.config.Device, "report-location"), 0, b.subscribeCommand(b.CommandReportLocation)),
		mqtt.NewSubscriber(MQTTSubscribeTopicCommand.Format(b.config.User, b.config.Device, "restart"), 0, b.subscribeCommand(b.CommandRestart)),
		mqtt.NewSubscriber(MQTTSubscribeTopicCommand.Format(b.config.User, b.config.Device, "reconnect"), 0, b.subscribeCommand(b.CommandReconnect)),
		mqtt.NewSubscriber(MQTTSubscribeTopicCommand.Format(b.config.User, b.config.Device, "waypoints"), 0, b.subscribeCommand(b.CommandWayPoints)),
		mqtt.NewSubscriber(MQTTOwnTracksSubscribeTopicUserLocation.Format(b.config.User, b.config.Device), 0, b.subscribeUserLocation),
	}

	if b.config.WayPointsSyncEnabled {
		subscribers = append(
			subscribers,
			mqtt.NewSubscriber(MQTTOwnTracksSubscribeTopicWayPoints.Format(b.config.User, b.config.Device), 0, b.subscribeSyncWayPoints),
		)
	}

	return subscribers
}

/*
	acc Accuracy of the reported location in meters without unit (iOS,Android/integer/meters/optional)
	alt Altitude measured above sea level (iOS,Android/integer/meters/optional)
	batt Device battery level (iOS,Android/integer/percent/optional)
	cog Course over ground (iOS/integer/degree/optional)
	lat latitude (iOS,Android/float/meters/required)
	lon longitude (iOS,Android/float/meters/required)
	rad radius around the region when entering/leaving (iOS/integer/meters/optional)
	t trigger for the location report (iOS,Android/string/optional)
		p ping issued randomly by background task (iOS,Android)
		c circular region enter/leave event (iOS,Android)
		b beacon region enter/leave event (iOS)
		r response to a reportLocation cmd message (iOS,Android)
		u manual publish requested by the user (iOS,Android)
		t timer based publish in move move (iOS)
		v updated by Settings/Privacy/Locations Services/System Services/Frequent Locations monitoring (iOS)
		tid Tracker ID used to display the initials of a user (iOS,Android/string/optional) required for http mode
	tst UNIX epoch timestamp in seconds of the location fix (iOS,Android/integer/epoch/required)
	vac vertical accuracy of the alt element (iOS/integer/meters/optional)
	vel velocity (iOS,Android/integer/kmh/optional)
	p barometric pressure (iOS/float/kPa/optional/extended data)
	conn Internet connectivity status (route to host) when the message is created (iOS,Android/string/optional/extended data)
		w phone is connected to a WiFi connection (iOS,Android)
		o phone is offline (iOS,Android)
		m mobile data (iOS,Android)
	topic (only in HTTP payloads) contains the original publish topic (e.g. owntracks/jane/phone). (iOS)
	inregions contains a list of regions the device is currently in (e.g. ["Home","Garage"]). Might be empty. (iOS,Android/list of strings/optional)
*/
func (b *Bind) subscribeUserLocation(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
	var payload *LocationPayload
	if err := json.Unmarshal(message.Payload(), &payload); err != nil {
		return err
	}

	if payload.Type == "lwt" {
		// skip last will and testament
		return nil
	}

	if payload.Type != "location" {
		return errors.New("location not found in payload")
	}

	if payload.Lat == nil {
		return errors.New("lat not found in payload")
	}

	var changeLocation bool
	q := message.Qos()
	r := message.Retained()

	if ok := b.lat.Set(*payload.Lat); ok {
		changeLocation = true

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateLat.Format(b.config.User, b.config.Device), q, r, *payload.Lat); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if payload.Lon == nil {
		return errors.New("lon not found in payload")
	}

	if ok := b.lon.Set(*payload.Lon); ok {
		changeLocation = true

		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateLon.Format(b.config.User, b.config.Device), q, r, *payload.Lon); e != nil {
			err = multierr.Append(err, e)
		}
	}

	// detect event
	if len(b.config.WayPoints) > 0 {
		existsRegions := make(map[string]struct{})

		if b.config.WayPointsCheckInRegionEnabled && payload.InRegions != nil {
			for _, name := range *payload.InRegions {
				var cache *atomic.BoolNull
				existsRegions[name] = struct{}{}

				if _, ok := b.wayPointsCheck[name]; ok {
					cache = b.wayPointsCheck[name]
				} else {
					b.mutex.Lock()
					cache, ok = b.wayPointsCheckUnregister[name]
					if !ok {
						cache = atomic.NewBoolNull()
						b.wayPointsCheckUnregister[name] = cache
					}
					b.mutex.Unlock()
				}

				if ok := cache.Set(true); ok {
					if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicRegion.Format(b.config.User, b.config.Device, mqtt.NameReplace(name)), q, r, true); e != nil {
						err = multierr.Append(err, e)
					}
				}
			}

			// все остальные не зарегистрированные в бинде отсылаем как false
			for name, cache := range b.wayPointsCheckUnregister {
				if _, ok := existsRegions[name]; ok {
					continue
				}

				if ok := cache.Set(false); ok {
					if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicRegion.Format(b.config.User, b.config.Device, mqtt.NameReplace(name)), q, r, false); e != nil {
						err = multierr.Append(err, e)
					}
				}
			}
		}

		for name, point := range b.config.WayPoints {
			if _, ok := existsRegions[name]; ok {
				continue
			}

			var checkResult bool

			if b.config.WayPointsCheckDistanceEnabled && b.validAccuracy(payload.Acc) {
				checkResult = calculateDistance(*payload.Lat, *payload.Lon, point.Lat, point.Lon) < point.Radius
			}

			if ok := b.wayPointsCheck[name].Set(checkResult); ok {
				if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicRegion.Format(b.config.User, b.config.Device, mqtt.NameReplace(name)), q, r, checkResult); e != nil {
					err = multierr.Append(err, e)
				}
			}
		}
	}

	hash := geohash.Encode(*payload.Lat, *payload.Lon)
	if ok := b.geoHash.Set(hash); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateGeoHash.Format(b.config.User, b.config.Device), q, r, hash); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if payload.Acc != nil {
		if ok := b.acc.Set(*payload.Acc); ok {
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateAccuracy.Format(b.config.User, b.config.Device), q, r, *payload.Acc); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	if payload.Alt != nil {
		if ok := b.alt.Set(*payload.Alt); ok {
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateAltitude.Format(b.config.User, b.config.Device), q, r, *payload.Alt); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	if changeLocation {
		location := strconv.FormatFloat(*payload.Lat, 'f', -1, 64) + "," + strconv.FormatFloat(*payload.Lon, 'f', -1, 64)
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateLocation.Format(b.config.User, b.config.Device), q, r, location); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if payload.Batt != nil {
		metricBatteryLevel.With("user", b.config.User, "device", b.config.Device).Set(*payload.Batt)

		if ok := b.batt.Set(*payload.Batt); ok {
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateBatteryLevel.Format(b.config.User, b.config.Device), q, r, *payload.Batt); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	if payload.Vel != nil {
		metricVelocity.With("user", b.config.User, "device", b.config.Device).Set(float64(*payload.Vel))

		if ok := b.vel.Set(*payload.Vel); ok {
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateVelocity.Format(b.config.User, b.config.Device), q, r, *payload.Vel); e != nil {
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
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateConnection.Format(b.config.User, b.config.Device), q, r, v); e != nil {
				err = multierr.Append(err, e)
			}
		}
	}

	return err
}

func (b *Bind) subscribeSyncWayPoints(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
	// _type == waypoint
	// одиночное добавление, вручную внесли новый пункт в список (или результат синка уже)
	// такие добавления игнорируем, тикет все равно дернет запрос всего списка
	var payloadOne *WayPointPayload
	if err = json.Unmarshal(message.Payload(), &payloadOne); err == nil && payloadOne.Type == "waypoint" {
		return nil
	}

	// _type == waypoints
	var payloadList *WayPointsPayload
	if err = json.Unmarshal(message.Payload(), &payloadList); err != nil {
		return err
	}

	if payloadList.Type != "waypoints" {
		return errors.New("payload isn't waypoints")
	}

	existsDesc := make(map[string]struct{})
	existsTst := make(map[int64]struct{})

	for _, point := range payloadList.WayPoints {
		existsTst[point.Tst] = struct{}{}
		existsDesc[point.Desc] = struct{}{}

		// заполняем кэш не зарегистрированных в конфигурации точек, но зарегистрированных на устройстве
		// необходимо, чтобы каждый раз не слать в MQTT статус по ним
		if b.config.WayPointsCheckInRegionEnabled {
			if _, ok := b.config.WayPoints[point.Desc]; !ok {
				b.mutex.Lock()
				if _, ok := b.wayPointsCheckUnregister[point.Desc]; !ok {
					b.wayPointsCheckUnregister[point.Desc] = atomic.NewBoolNull()
				}
				b.mutex.Unlock()
			}
		}
	}

	points := make([]WayPointPayload, 0, len(b.config.WayPoints))
	lastTst := time.Now().Unix()
	for id, point := range b.config.WayPoints {
		if _, ok := existsDesc[id]; ok {
			continue
		}

		// generate tst
		for {
			lastTst++

			if _, ok := existsTst[lastTst]; !ok {
				existsTst[lastTst] = struct{}{}
				break
			}
		}

		points = append(points, WayPointPayload{
			Desc: id,
			Lat:  point.Lat,
			Lon:  point.Lon,
			Rad:  point.Radius,
			Tst:  lastTst,
		})
	}

	if len(points) == 0 {
		return nil
	}

	return b.CommandSetWayPoints(points)
}

func (b *Bind) subscribeCommand(cmd func() error) mqtt.MessageHandler {
	return func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
		return cmd()
	}
}
