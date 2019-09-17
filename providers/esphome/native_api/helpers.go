package native_api

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
)

type MessageEntity interface {
	GetObjectId() string
	GetKey() uint32
	GetName() string
	GetUniqueId() string
}

type MessageState interface {
	GetKey() uint32
}

const (
	EntityTypeUnknown      = "unknown"
	EntityTypeBinarySensor = "binary_sensor"
	EntityTypeCover        = "cover"
	EntityTypeFan          = "fan"
	EntityTypeLight        = "light"
	EntityTypeSensor       = "sensor"
	EntityTypeSwitch       = "switch"
	EntityTypeTextSensor   = "text_sensor"
	EntityTypeCamera       = "camera"
	EntityTypeClimate      = "climate"
)

func EntityType(message proto.Message) string {
	switch message.(type) {
	case *ListEntitiesBinarySensorResponse:
		return EntityTypeBinarySensor
	case *ListEntitiesCoverResponse:
		return EntityTypeCover
	case *ListEntitiesFanResponse:
		return EntityTypeFan
	case *ListEntitiesLightResponse:
		return EntityTypeLight
	case *ListEntitiesSensorResponse:
		return EntityTypeSensor
	case *ListEntitiesSwitchResponse:
		return EntityTypeSwitch
	case *ListEntitiesTextSensorResponse:
		return EntityTypeTextSensor
	case *ListEntitiesCameraResponse:
		return EntityTypeCamera
	case *ListEntitiesClimateResponse:
		return EntityTypeClimate
	}

	return EntityTypeUnknown
}

func State(entityMessage, stateMessage proto.Message) (state string, err error) {
	switch v := stateMessage.(type) {
	case *BinarySensorStateResponse:
		if v.GetState() {
			state = "on"
		} else {
			state = "off"
		}
	case *CoverStateResponse:
		// TODO:
	case *FanStateResponse:
		if v.GetState() {
			state = "on"
		} else {
			state = "off"
		}
	case *LightStateResponse:
		if v.GetState() {
			state = "on"
		} else {
			state = "off"
		}
	case *SensorStateResponse:
		state = strconv.FormatFloat(float64(v.GetState()), 'f', -1, 64)

		if e, ok := entityMessage.(*ListEntitiesSensorResponse); ok {
			state = fmt.Sprintf("%s "+e.GetUnitOfMeasurement(), state)
		}
	case *SwitchStateResponse:
		if v.GetState() {
			state = "on"
		} else {
			state = "off"
		}
	case *TextSensorStateResponse:
		state = v.GetState()
	case *ClimateStateResponse:
		state = v.GetMode().String()
	default:
		err = errors.New("unknown state type " + proto.MessageName(stateMessage))
	}

	return
}
