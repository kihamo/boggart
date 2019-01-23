package owntracks

const (
	DefaultPointRadius             = 100
	DefaultMaxAccuracy             = 100
	DefaultWayPointsSyncEnabled    = true
	DefaultUnregisterPointsAllowed = true
	DefaultCheckInRegionEnabled    = true
	DefaultCheckDistanceEnabled    = true
)

type Point struct {
	Lat    float64
	Lon    float64
	Radius float64
}

type Config struct {
	User   string `valid:"required"`
	Device string `valid:"required"`
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
	CheckDistanceEnabled bool `mapstructure:"check_distance_enabled" yaml:"check_distance_enabled"`
}

func (Type) Config() interface{} {
	return &Config{
		MaxAccuracy:             DefaultMaxAccuracy,
		UnregisterPointsAllowed: DefaultUnregisterPointsAllowed,
		RegionsSyncEnabled:      DefaultWayPointsSyncEnabled,
		CheckInRegionEnabled:    DefaultCheckInRegionEnabled,
		CheckDistanceEnabled:    DefaultCheckDistanceEnabled,
	}
}
