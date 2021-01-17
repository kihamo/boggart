package mqtt

import (
	"github.com/kihamo/boggart/components/mqtt"
)

type ComponentType string

const (
	ComponentTypeUnknown      ComponentType = "unknown"
	ComponentTypeBinarySensor ComponentType = "binary_sensor"
	ComponentTypeCover        ComponentType = "cover"
	ComponentTypeFan          ComponentType = "fan"
	ComponentTypeLight        ComponentType = "light"
	ComponentTypeSensor       ComponentType = "sensor"
	ComponentTypeSwitch       ComponentType = "switch"
	ComponentTypeTextSensor   ComponentType = "text_sensor"
	ComponentTypeCamera       ComponentType = "camera"
	ComponentTypeClimate      ComponentType = "climate"
)

func (t ComponentType) String() string {
	return string(t)
}

type Component interface {
	ID() string
	Type() ComponentType
	UniqueID() string
	Name() string
	Icon() string
	State() interface{}
	StateFormat() string
	SetState(mqtt.Message) error
	StateTopic() mqtt.Topic
	CommandTopic() mqtt.Topic
	AvailabilityTopic() mqtt.Topic
	Device() Device
	CommandToPayload(cmd interface{}) interface{}
}
