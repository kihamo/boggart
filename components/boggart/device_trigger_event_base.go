package boggart

import (
	"context"

	"github.com/kihamo/go-workers"
)

type DeviceTriggerEventBase struct {
	ctx       context.Context
	event     workers.Event
	arguments []interface{}
}

func NewDeviceTriggerEventBase(ctx context.Context, event workers.Event, arguments []interface{}) *DeviceTriggerEventBase {
	return &DeviceTriggerEventBase{
		ctx:       ctx,
		event:     event,
		arguments: arguments,
	}
}

func (e *DeviceTriggerEventBase) Context() context.Context {
	return e.ctx
}

func (e *DeviceTriggerEventBase) Event() workers.Event {
	return e.event
}

func (e *DeviceTriggerEventBase) Arguments() []interface{} {
	return e.arguments
}
