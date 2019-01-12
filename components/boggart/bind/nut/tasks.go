package nut

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/robbiet480/go.nut"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Second * 30)
	taskLiveness.SetName("bind-nut-liveness-" + b.config.UPS)

	taskStateUpdater := task.NewFunctionTask(b.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(time.Minute)
	taskStateUpdater.SetName("bind-nut-state-updater-" + b.config.UPS)

	return []workers.Task{
		taskLiveness,
		taskStateUpdater,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	_, err := b.getUPS()
	if err != nil {
		b.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.DeviceStatusOnline)
	return nil, nil
}

func (b *Bind) taskStateUpdater(ctx context.Context) (interface{}, error) {
	if b.Status() != boggart.DeviceStatusOnline {
		return nil, nil
	}

	ups, err := b.getUPS()
	if err != nil {
		return nil, err
	}

	b.mutex.Lock()

	for _, v := range ups.Variables {
		prev, ok := b.variables[v.Name]
		if !ok || prev != v.Value {
			b.variables[v.Name] = v.Value
			name := mqtt.NameReplace(v.Name)

			b.MQTTPublishAsync(ctx, MQTTTopicVariable.Format(b.config.UPS, name), 2, true, v.Value)
		}
	}

	b.mutex.Unlock()
	return nil, nil
}

func (b *Bind) getUPS() (ups nut.UPS, err error) {
	client, err := nut.Connect(b.config.Host)
	if err != nil {
		return ups, err
	}

	devices, err := client.GetUPSList()
	if err != nil {
		return ups, err
	}

	for _, device := range devices {
		if device.Name == b.config.UPS {
			for _, v := range device.Variables {
				if v.Name == "device.serial" {
					b.SetSerialNumber(v.Value.(string))
					return device, nil
				}
			}

			break
		}
	}

	return ups, errors.New("device not found")
}
