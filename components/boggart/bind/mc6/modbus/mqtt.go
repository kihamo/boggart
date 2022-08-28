package modbus

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTWrapSubscriberChangeState(f func(message mqtt.Message) error) func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	return func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
		err := f(message)

		if err == nil {
			err = b.Workers().TaskRunByName(ctx, "state-updater")
			if err != nil && err == tasks.ErrAlreadyRunning {
				return nil
			}
		}

		return err
	}
}
