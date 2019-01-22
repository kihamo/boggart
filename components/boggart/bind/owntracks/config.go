package owntracks

const (
	DefaultPointRadius                   = 100
	DefaultMaxAccuracyDefaultMaxAccuracy = 100
	DefaultWayPointsSyncEnabled          = true
	DefaultWayPointsCheckInRegionEnabled = true
	DefaultWayPointsCheckDistanceEnabled = true
)

type Point struct {
	Lat    float64
	Lon    float64
	Radius float64
}

type Config struct {
	User                          string `valid:"required"`
	Device                        string `valid:"required"`
	MaxAccuracy                   int64  `mapstructure:"max_accuracy" yaml:"max_accuracy"`
	WayPoints                     map[string]Point
	WayPointsSyncEnabled          bool `mapstructure:"waypoints_sync_enabled" yaml:"waypoints_sync_enabled"`
	WayPointsCheckInRegionEnabled bool `mapstructure:"waypoints_check_in_region_enabled" yaml:"waypoints_check_in_region_enabled"`
	WayPointsCheckDistanceEnabled bool `mapstructure:"waypoints_check_distance_enabled" yaml:"waypoints_check_distance_enabled"`
}

func (Type) Config() interface{} {
	return &Config{
		MaxAccuracy:                   DefaultMaxAccuracyDefaultMaxAccuracy,
		WayPointsSyncEnabled:          DefaultWayPointsSyncEnabled,
		WayPointsCheckInRegionEnabled: DefaultWayPointsCheckInRegionEnabled,
		WayPointsCheckDistanceEnabled: DefaultWayPointsCheckDistanceEnabled,
	}
}
