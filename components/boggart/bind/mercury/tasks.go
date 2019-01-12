package mercury

import (
	"context"
	"math"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (d *Bind) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(time.Minute)
	taskStateUpdater.SetName("bind-mercury-200-state-updater-" + d.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (d *Bind) taskStateUpdater(ctx context.Context) (interface{}, error) {
	currentT1, currentT2, currentT3, currentT4, err := d.provider.PowerCounters()
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		// TODO: log
		return nil, err
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)
	serialNumber := d.SerialNumber()

	if prevT1 := atomic.LoadUint64(&d.tariff1); currentT1 != prevT1 {
		atomic.StoreUint64(&d.tariff1, currentT1)
		d.MQTTPublishAsync(ctx, MQTTTopicTariff.Format(serialNumber, 1), 0, true, currentT1)
	}

	if prevT2 := atomic.LoadUint64(&d.tariff2); currentT2 != prevT2 {
		atomic.StoreUint64(&d.tariff2, currentT2)
		d.MQTTPublishAsync(ctx, MQTTTopicTariff.Format(serialNumber, 2), 0, true, currentT2)
	}

	if prevT3 := atomic.LoadUint64(&d.tariff3); currentT3 != prevT3 {
		atomic.StoreUint64(&d.tariff3, currentT3)
		d.MQTTPublishAsync(ctx, MQTTTopicTariff.Format(serialNumber, 3), 0, true, currentT3)
	}

	if prevT4 := atomic.LoadUint64(&d.tariff4); currentT4 != prevT4 {
		atomic.StoreUint64(&d.tariff4, currentT4)
		d.MQTTPublishAsync(ctx, MQTTTopicTariff.Format(serialNumber, 4), 0, true, currentT4)
	}

	// optimization
	if currentVoltage, currentAmperage, currentPower, err := d.provider.ParamsCurrent(); err == nil {
		if prevVoltage := atomic.LoadUint64(&d.voltage); currentVoltage != prevVoltage {
			atomic.StoreUint64(&d.voltage, currentVoltage)
			d.MQTTPublishAsync(ctx, MQTTTopicVoltage.Format(serialNumber), 0, true, currentVoltage)
		}

		if prevAmperage := math.Float64frombits(atomic.LoadUint64(&d.amperage)); currentAmperage != prevAmperage {
			atomic.StoreUint64(&d.amperage, math.Float64bits(currentAmperage))
			d.MQTTPublishAsync(ctx, MQTTTopicAmperage.Format(serialNumber), 0, true, currentAmperage)
		}

		if prevPower := atomic.LoadUint64(&d.power); currentPower != prevPower {
			atomic.StoreUint64(&d.power, currentPower)
			d.MQTTPublishAsync(ctx, MQTTTopicPower.Format(serialNumber), 0, true, currentPower)
		}
	} else {
		// TODO: log
	}

	if current, err := d.provider.BatteryVoltage(); err == nil {
		if prev := math.Float64frombits(atomic.LoadUint64(&d.batteryVoltage)); current != prev {
			atomic.StoreUint64(&d.batteryVoltage, math.Float64bits(current))
			d.MQTTPublishAsync(ctx, MQTTTopicBatteryVoltage.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if val, err := d.provider.LastPowerOffDatetime(); err == nil {
		current := val.Unix()
		if prev := atomic.LoadInt64(&d.lastPowerOff); current != prev {
			atomic.StoreInt64(&d.lastPowerOff, current)
			d.MQTTPublishAsync(ctx, MQTTTopicLastPowerOff.Format(serialNumber), 0, true, val)
		}
	} else {
		// TODO: log
	}

	if val, err := d.provider.LastPowerOnDatetime(); err == nil {
		current := val.Unix()
		if prev := atomic.LoadInt64(&d.lastPowerOn); current != prev {
			atomic.StoreInt64(&d.lastPowerOn, current)
			d.MQTTPublishAsync(ctx, MQTTTopicLastPowerOn.Format(serialNumber), 0, true, val)
		}
	} else {
		// TODO: log
	}

	if val, err := d.provider.MakeDate(); err == nil {
		current := val.Unix()
		if prev := atomic.LoadInt64(&d.makeDate); current != prev {
			atomic.StoreInt64(&d.makeDate, current)
			d.MQTTPublishAsync(ctx, MQTTTopicMakeDate.Format(serialNumber), 0, true, val)
		}
	} else {
		// TODO: log
	}

	if version, date, err := d.provider.Version(); err == nil {
		current := date.Unix()
		if prev := atomic.LoadInt64(&d.firmwareDate); current != prev {
			atomic.StoreInt64(&d.firmwareDate, current)
			d.MQTTPublishAsync(ctx, MQTTTopicFirmwareDate.Format(serialNumber), 0, true, date)
		}

		d.mutex.Lock()
		if version != d.firmwareVersion {
			d.firmwareVersion = version
			d.MQTTPublishAsync(ctx, MQTTTopicFirmwareVersion.Format(serialNumber), 0, true, version)
		}
		d.mutex.Unlock()
	} else {
		// TODO: log
	}

	return nil, nil
}
