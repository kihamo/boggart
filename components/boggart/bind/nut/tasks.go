package nut

import (
	"context"

	"github.com/kihamo/go-workers"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := b.WrapTaskIsOnline(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater")

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	variables, err := b.Variables()
	if err != nil {
		return err
	}

	b.mutex.Lock()

	sn := b.SerialNumber()

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

			// TODO:
			_ = b.MQTTContainer().PublishAsync(ctx, b.config.TopicVariable.Format(sn, v.Name), v.Value)
		}
	}

	b.mutex.Unlock()
	return nil
}
