package owntracks

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Point struct {
	Lat    float64
	Lon    float64
	Radius float64
}

type Config struct {
	User        string  `valid:"required"`
	Device      string  `valid:"required"`
	PointRadius float64 `mapstructure:"point_radius" yaml:"point_radius"`
	// максимальное значение погрешности с которым допустимо производить операции
	MaxAccuracy int64 `mapstructure:"max_accuracy" yaml:"max_accuracy"`
	// точки для определения регионов
	Regions map[string]Point
	// разрешить обрабатывать события по не зарегистрированным в бинде точкам (но они могут быть в списке на устройстве)
	UnregisterPointsAllowed bool `mapstructure:"unregister_points_enabled" yaml:"unregister_points_enabled"`
	// включить принудятельную регистрацию точек из бинда на устройстве, в случае если таких нет
	RegionsSyncEnabled bool `mapstructure:"regions_sync_enabled" yaml:"regions_sync_enabled"`
	// обрабатывать поле in_regions в _type=location как событие присутсвия
	CheckInRegionEnabled bool `mapstructure:"check_in_region_enabled" yaml:"check_in_region_enabled"`
	// включить расчет дистанции до точек, чтобы сработало событие
	CheckDistanceEnabled       bool       `mapstructure:"check_distance_enabled" yaml:"check_distance_enabled"`
	TopicOwnTracksUserLocation mqtt.Topic `mapstructure:"topic_owntracks_user_location" yaml:"topic_owntracks_user_location"`
	TopicOwnTracksTransition   mqtt.Topic `mapstructure:"topic_owntracks_transition" yaml:"topic_owntracks_transition"`
	TopicOwnTracksStep         mqtt.Topic `mapstructure:"topic_owntracks_step" yaml:"topic_owntracks_step"`
	TopicOwnTracksBeacon       mqtt.Topic `mapstructure:"topic_owntracks_beacon" yaml:"topic_owntracks_beacon"`
	TopicOwnTracksDump         mqtt.Topic `mapstructure:"topic_owntracks_dump" yaml:"topic_owntracks_dump"`
	TopicOwnTracksWayPoints    mqtt.Topic `mapstructure:"topic_owntracks_waypoints" yaml:"topic_owntracks_waypoints"`
	TopicOwnTracksCommand      mqtt.Topic `mapstructure:"topic_owntracks_command" yaml:"topic_owntracks_command"`
	TopicOwnTracksCard         mqtt.Topic `mapstructure:"topic_owntracks_card" yaml:"topic_owntracks_card"`
	TopicCommandReportLocation mqtt.Topic `mapstructure:"topic_command_report_location" yaml:"topic_command_report_location"`
	TopicCommandRestart        mqtt.Topic `mapstructure:"topic_command_restart" yaml:"topic_command_restart"`
	TopicCommandReconnect      mqtt.Topic `mapstructure:"topic_command_reconnect" yaml:"topic_command_reconnect"`
	TopicCommandWayPoints      mqtt.Topic `mapstructure:"topic_command_waypoints" yaml:"topic_command_waypoints"`
	TopicRegion                mqtt.Topic `mapstructure:"topic_region" yaml:"topic_region"`
	TopicStateLat              mqtt.Topic `mapstructure:"topic_state_lat" yaml:"topic_state_lat"`
	TopicStateLon              mqtt.Topic `mapstructure:"topic_state_lon" yaml:"topic_state_lon"`
	TopicStateGeoHash          mqtt.Topic `mapstructure:"topic_state_geohash" yaml:"topic_state_geohash"`
	TopicStateAccuracy         mqtt.Topic `mapstructure:"topic_state_accuracy" yaml:"topic_state_accuracy"`
	TopicStateAltitude         mqtt.Topic `mapstructure:"topic_state_altitude" yaml:"topic_state_altitude"`
	TopicStateBatteryLevel     mqtt.Topic `mapstructure:"topic_state_battery_level" yaml:"topic_state_battery_level"`
	TopicStateVelocity         mqtt.Topic `mapstructure:"topic_state_velocity" yaml:"topic_state_velocity"`
	TopicStateConnection       mqtt.Topic `mapstructure:"topic_state_connection" yaml:"topic_state_connection"`
	TopicStateLocation         mqtt.Topic `mapstructure:"topic_state_location" yaml:"topic_state_location"`
}

func (Type) Config() interface{} {
	var (
		prefixOT mqtt.Topic = "owntracks/+/+"
		prefix   mqtt.Topic = boggart.ComponentName + "/owntracks/+/+/"
	)

	return &Config{
		PointRadius:                100,
		MaxAccuracy:                100,
		UnregisterPointsAllowed:    true,
		RegionsSyncEnabled:         true,
		CheckInRegionEnabled:       true,
		CheckDistanceEnabled:       true,
		TopicOwnTracksUserLocation: prefixOT,
		TopicOwnTracksTransition:   prefixOT + "/event",
		TopicOwnTracksStep:         prefixOT + "/step",
		TopicOwnTracksBeacon:       prefixOT + "/beacon",
		TopicOwnTracksDump:         prefixOT + "/dump",
		TopicOwnTracksWayPoints:    prefixOT + "/waypoint",
		TopicOwnTracksCommand:      prefixOT + "/cmd",
		TopicOwnTracksCard:         prefixOT + "/info",
		TopicCommandReportLocation: prefix + "cmd/report-location",
		TopicCommandRestart:        prefix + "cmd/restart",
		TopicCommandReconnect:      prefix + "cmd/reconnect",
		TopicCommandWayPoints:      prefix + "cmd/waypoints",
		TopicRegion:                prefix + "event/+",
		TopicStateLat:              prefix + "state/lat",
		TopicStateLon:              prefix + "state/lon",
		TopicStateGeoHash:          prefix + "state/geohash",
		TopicStateAccuracy:         prefix + "state/accuracy",
		TopicStateAltitude:         prefix + "state/altitude",
		TopicStateBatteryLevel:     prefix + "state/battery-level",
		TopicStateVelocity:         prefix + "state/velocity",
		TopicStateConnection:       prefix + "state/connection",
		TopicStateLocation:         prefix + "state/location",
	}
}
