package owntracks

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	config.TopicOwnTracksUserLocation = config.TopicOwnTracksUserLocation.Format(config.User, config.Device)
	config.TopicOwnTracksTransition = config.TopicOwnTracksTransition.Format(config.User, config.Device)
	config.TopicOwnTracksStep = config.TopicOwnTracksStep.Format(config.User, config.Device)
	config.TopicOwnTracksBeacon = config.TopicOwnTracksBeacon.Format(config.User, config.Device)
	config.TopicOwnTracksDump = config.TopicOwnTracksDump.Format(config.User, config.Device)
	config.TopicOwnTracksWayPoints = config.TopicOwnTracksWayPoints.Format(config.User, config.Device)
	config.TopicOwnTracksCommand = config.TopicOwnTracksCommand.Format(config.User, config.Device)
	config.TopicOwnTracksCard = config.TopicOwnTracksCard.Format(config.User, config.Device)
	config.TopicCommandReportLocation = config.TopicCommandReportLocation.Format(config.User, config.Device)
	config.TopicCommandRestart = config.TopicCommandRestart.Format(config.User, config.Device)
	config.TopicCommandReconnect = config.TopicCommandReconnect.Format(config.User, config.Device)
	config.TopicCommandWayPoints = config.TopicCommandWayPoints.Format(config.User, config.Device)
	config.TopicRegion = config.TopicRegion.Format(config.User, config.Device)
	config.TopicStateLat = config.TopicStateLat.Format(config.User, config.Device)
	config.TopicStateLon = config.TopicStateLon.Format(config.User, config.Device)
	config.TopicStateGeoHash = config.TopicStateGeoHash.Format(config.User, config.Device)
	config.TopicStateAccuracy = config.TopicStateAccuracy.Format(config.User, config.Device)
	config.TopicStateAltitude = config.TopicStateAltitude.Format(config.User, config.Device)
	config.TopicStateBatteryLevel = config.TopicStateBatteryLevel.Format(config.User, config.Device)
	config.TopicStateVelocity = config.TopicStateVelocity.Format(config.User, config.Device)
	config.TopicStateConnection = config.TopicStateConnection.Format(config.User, config.Device)
	config.TopicStateLocation = config.TopicStateLocation.Format(config.User, config.Device)

	bind := &Bind{
		config:   config,
		lat:      atomic.NewFloat32Null(),
		lon:      atomic.NewFloat32Null(),
		geoHash:  atomic.NewString(),
		conn:     atomic.NewString(),
		acc:      atomic.NewInt32Null(),
		alt:      atomic.NewInt32Null(),
		batt:     atomic.NewFloat32Null(),
		vel:      atomic.NewInt32Null(),
		regions:  make(map[string]Point),
		checkers: make(map[string]*atomic.BoolNull),
	}

	for name, region := range bind.config.Regions {
		if region.Radius <= 0 {
			region.Radius = config.PointRadius
		}

		bind.config.Regions[name] = region
		bind.registerRegion(name, region)
	}

	return bind, nil
}
