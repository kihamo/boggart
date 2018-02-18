package boggart

import (
	"github.com/kihamo/go-workers"
)

type DeviceTriggerEventBase struct {
	event     workers.Event
	arguments []interface{}
}

func NewDeviceTriggerEventBase(event workers.Event, arguments []interface{}) *DeviceTriggerEventBase {
	return &DeviceTriggerEventBase{
		event:     event,
		arguments: arguments,
	}
}

func (e *DeviceTriggerEventBase) Event() workers.Event {
	return e.event
}

func (e *DeviceTriggerEventBase) Arguments() []interface{} {
	return e.arguments
}
