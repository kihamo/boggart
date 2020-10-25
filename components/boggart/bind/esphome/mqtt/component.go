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
	GetID() string
	GetType() ComponentType
	GetUniqueID() string
	GetName() string
	GetState() interface{}
	GetCommandTopic() mqtt.Topic
	GetDevice() Device
	CommandToPayload(cmd interface{}) interface{}
	Subscribers() []mqtt.Subscriber
	TopicState() mqtt.Topic
}
