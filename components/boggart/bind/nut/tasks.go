package nut

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater-" + b.config.UPS)

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	variables, err := b.Variables()
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.BindStatusOnline)
	b.mutex.Lock()

	sn := b.SerialNumber()
	snMQTT := mqtt.NameReplace(sn)

	for _, v := range variables {
		prev, ok := b.variables[v.Name]
		if !ok || prev != v.Value {
			switch v.Name {
			case "ups.load":
				metricLoad.With("serial_number", sn).Set(float64(v.Value.(int64)))
			case "input.voltage":
				metricInputVoltage.With("serial_number", sn).Set(v.Value.(float64))
			case "battery.charge":
				metricBatteryCharge.With("serial_number", sn).Set(float64(v.Value.(int64)))
			case "battery.runtime":
				metricBatteryRuntime.With("serial_number", sn).Set(float64(v.Value.(int64)))
			case "battery.voltage":
				metricBatteryVoltage.With("serial_number", sn).Set(v.Value.(float64))
			}

			b.variables[v.Name] = v.Value
			name := mqtt.NameReplace(v.Name)

			// TODO:
			_ = b.MQTTPublishAsync(ctx, MQTTPublishTopicVariable.Format(snMQTT, name), v.Value)
		}
	}

	b.mutex.Unlock()
	return nil, nil
}
