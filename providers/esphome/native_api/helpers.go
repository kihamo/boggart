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

func State(entityMessage, stateMessage proto.Message) (state string, raw interface{}, err error) {
	switch v := stateMessage.(type) {
	case *BinarySensorStateResponse:
		raw = v.GetKey()
		if v.GetState() {
			state = "on"
		} else {
			state = "off"
		}
	case *CoverStateResponse:
		// TODO:
	case *FanStateResponse:
		raw = v.GetState()
		if v.GetState() {
			state = "on"
		} else {
			state = "off"
		}
	case *LightStateResponse:
		raw = v.GetState()
		if v.GetState() {
			state = "on"
		} else {
			state = "off"
		}
	case *SensorStateResponse:
		state = strconv.FormatFloat(float64(v.GetState()), 'f', -1, 64)
		raw = state

		if e, ok := entityMessage.(*ListEntitiesSensorResponse); ok {
			state = fmt.Sprintf("%s "+e.GetUnitOfMeasurement(), state)
		}
	case *SwitchStateResponse:
		raw = v.GetState()
		if v.GetState() {
			state = "on"
		} else {
			state = "off"
		}
	case *TextSensorStateResponse:
		state = v.GetState()
		raw = state
	case *ClimateStateResponse:
		state = v.GetMode().String()
		raw = state
	default:
		err = errors.New("unknown state type " + proto.MessageName(stateMessage))
	}

	return
}
