package owntracks

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

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
	MQTTSubscribeTopicUserLocation mqtt.Topic = "owntracks/+/+"
	MQTTSubscribeTopicCommand      mqtt.Topic = "owntracks/+/+/cmd"
	MQTTSubscribeTopicTransition   mqtt.Topic = "owntracks/+/+/event"
	MQTTSubscribeTopicStep         mqtt.Topic = "owntracks/+/+/step"
	MQTTSubscribeTopicBeacon       mqtt.Topic = "owntracks/+/+/beacon"
	MQTTSubscribeTopicDump         mqtt.Topic = "owntracks/+/+/dump"
	MQTTSubscribeTopicWayPoints    mqtt.Topic = "owntracks/+/+/waypoint"

	MQTTPublishTopicCommand      mqtt.Topic = "owntracks/+/+/cmd"
	MQTTPublishTopicUserLocation mqtt.Topic = "owntracks/+/+"
	MQTTPublishTopicTransition   mqtt.Topic = "owntracks/+/+/event"
	MQTTPublishTopicCard         mqtt.Topic = "owntracks/+/+/info"

	MQTTPublishTopicUserStateGeoHash      mqtt.Topic = "owntracks/+/+/state/geohash"
	MQTTPublishTopicUserStateBatteryLevel mqtt.Topic = "owntracks/+/+/state/battery-level"
	MQTTPublishTopicUserStateVelocity     mqtt.Topic = "owntracks/+/+/state/velocity"
	MQTTPublishTopicUserStateConnection   mqtt.Topic = "owntracks/+/+/state/connection"
	MQTTPublishTopicUserStateLocation     mqtt.Topic = "owntracks/+/+/state/location"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicUserStateGeoHash,
		MQTTPublishTopicUserStateBatteryLevel,
		MQTTPublishTopicUserStateVelocity,
		MQTTPublishTopicUserStateConnection,
		MQTTPublishTopicUserStateLocation,
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
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
		mqtt.NewSubscriber(MQTTSubscribeTopicUserLocation.String(), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 2 {
				return errors.New("bad topic name")
			}

			user := route[len(route)-2]
			device := route[len(route)-1]
			q := message.Qos()
			r := message.Retained()

			var payload map[string]interface{}
			if err := json.Unmarshal(message.Payload(), &payload); err != nil {
				return err
			}

			t, ok := payload["_type"]
			if ok && t == "lwt" {
				// skip last will and testament
				return nil
			}

			if !ok || t != "location" {
				return errors.New("location not found in payload")
			}

			lat, ok := payload["lat"]
			if !ok {
				return errors.New("lat not found in payload")
			}

			lon, ok := payload["lon"]
			if !ok {
				return errors.New("lon not found in payload")
			}

			location := strconv.FormatFloat(lat.(float64), 'f', -1, 64) + "," + strconv.FormatFloat(lon.(float64), 'f', -1, 64)
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateLocation.Format(user, device), q, r, location); e != nil {
				err = multierr.Append(err, e)
			}

			if batteryLevel, ok := payload["batt"]; ok {
				metricBatteryLevel.With("user", user, "device", device).Set(batteryLevel.(float64))

				if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateBatteryLevel.Format(user, device), q, r, batteryLevel); e != nil {
					err = multierr.Append(err, e)
				}
			}

			if velocity, ok := payload["vel"]; ok {
				metricVelocity.With("user", user, "device", device).Set(velocity.(float64))

				if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateVelocity.Format(user, device), q, r, velocity); e != nil {
					err = multierr.Append(err, e)
				}
			}

			if connection, ok := payload["conn"]; ok {
				var v string
				switch connection {
				case "w":
					v = "wifi"
				case "o":
					v = "offline"
				case "m":
					v = "mobile"
				default:
					v = "unknown"
				}

				if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateConnection.Format(user, device), q, r, v); e != nil {
					err = multierr.Append(err, e)
				}
			}

			hash := geohash.Encode(lat.(float64), lon.(float64))
			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicUserStateGeoHash.Format(user, device), q, r, hash); e != nil {
				err = multierr.Append(err, e)
			}

			return err
		}),
		mqtt.NewSubscriber(MQTTSubscribeTopicWayPoints.String(), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 2 {
				return errors.New("bad topic name")
			}

			/*
				user := route[len(route)-2]
				device := route[len(route)-1]
				q := message.Qos()
				r := message.Retained()
			*/

			var payload map[string]interface{}
			if err := json.Unmarshal(message.Payload(), &payload); err != nil {
				return err
			}

			// fmt.Println(payload)

			return nil
		}),
	}
}
