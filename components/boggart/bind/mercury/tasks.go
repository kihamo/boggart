package mercury

import (
	"context"
	"math"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.updaterInterval)
	taskStateUpdater.SetName("bind-mercury:200-updater-" + b.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	currentT1, currentT2, currentT3, currentT4, err := b.provider.PowerCounters()
	if err != nil {
		b.UpdateStatus(boggart.DeviceStatusOffline)
		// TODO: log
		return nil, err
	}

	b.UpdateStatus(boggart.DeviceStatusOnline)

	sn := b.SerialNumber()
	snMQTT := mqtt.NameReplace(sn)
	mTariff := metricTariff.With("serial_number", sn)

	if prevT1 := atomic.LoadUint64(&b.tariff1); currentT1 != prevT1 {
		atomic.StoreUint64(&b.tariff1, currentT1)
		b.MQTTPublishAsync(ctx, MQTTTopicTariff.Format(snMQTT, 1), 0, true, currentT1)
		mTariff.With("tariff", "1").Set(float64(currentT1))
	}

	if prevT2 := atomic.LoadUint64(&b.tariff2); currentT2 != prevT2 {
		atomic.StoreUint64(&b.tariff2, currentT2)
		b.MQTTPublishAsync(ctx, MQTTTopicTariff.Format(snMQTT, 2), 0, true, currentT2)
		mTariff.With("tariff", "2").Set(float64(currentT2))
	}

	if prevT3 := atomic.LoadUint64(&b.tariff3); currentT3 != prevT3 {
		atomic.StoreUint64(&b.tariff3, currentT3)
		b.MQTTPublishAsync(ctx, MQTTTopicTariff.Format(snMQTT, 3), 0, true, currentT3)
		mTariff.With("tariff", "3").Set(float64(currentT3))
	}

	if prevT4 := atomic.LoadUint64(&b.tariff4); currentT4 != prevT4 {
		atomic.StoreUint64(&b.tariff4, currentT4)
		b.MQTTPublishAsync(ctx, MQTTTopicTariff.Format(snMQTT, 4), 0, true, currentT4)
		mTariff.With("tariff", "4").Set(float64(currentT4))
	}

	// optimization
	if currentVoltage, currentAmperage, currentPower, err := b.provider.ParamsCurrent(); err == nil {
		if prevVoltage := atomic.LoadUint64(&b.voltage); currentVoltage != prevVoltage {
			atomic.StoreUint64(&b.voltage, currentVoltage)
			b.MQTTPublishAsync(ctx, MQTTTopicVoltage.Format(snMQTT), 0, true, currentVoltage)
			metricVoltage.With("serial_number", sn).Set(float64(currentVoltage))
		}

		if prevAmperage := math.Float64frombits(atomic.LoadUint64(&b.amperage)); currentAmperage != prevAmperage {
			atomic.StoreUint64(&b.amperage, math.Float64bits(currentAmperage))
			b.MQTTPublishAsync(ctx, MQTTTopicAmperage.Format(snMQTT), 0, true, currentAmperage)
			metricAmperage.With("serial_number", sn).Set(currentAmperage)
		}

		if prevPower := atomic.LoadUint64(&b.power); currentPower != prevPower {
			atomic.StoreUint64(&b.power, currentPower)
			b.MQTTPublishAsync(ctx, MQTTTopicPower.Format(snMQTT), 0, true, currentPower)
			metricPower.With("serial_number", sn).Set(float64(currentPower))
		}
	} else {
		// TODO: log
	}

	if current, err := b.provider.BatteryVoltage(); err == nil {
		if prev := math.Float64frombits(atomic.LoadUint64(&b.batteryVoltage)); current != prev {
			atomic.StoreUint64(&b.batteryVoltage, math.Float64bits(current))
			b.MQTTPublishAsync(ctx, MQTTTopicBatteryVoltage.Format(snMQTT), 0, true, current)
			metricBatteryVoltage.With("serial_number", sn).Set(float64(current))
		}
	} else {
		// TODO: log
	}

	if val, err := b.provider.LastPowerOffDatetime(); err == nil {
		current := val.Unix()
		if prev := atomic.LoadInt64(&b.lastPowerOff); current != prev {
			atomic.StoreInt64(&b.lastPowerOff, current)
			b.MQTTPublishAsync(ctx, MQTTTopicLastPowerOff.Format(snMQTT), 0, true, val)
		}
	} else {
		// TODO: log
	}

	if val, err := b.provider.LastPowerOnDatetime(); err == nil {
		current := val.Unix()
		if prev := atomic.LoadInt64(&b.lastPowerOn); current != prev {
			atomic.StoreInt64(&b.lastPowerOn, current)
			b.MQTTPublishAsync(ctx, MQTTTopicLastPowerOn.Format(snMQTT), 0, true, val)
		}
	} else {
		// TODO: log
	}

	if val, err := b.provider.MakeDate(); err == nil {
		current := val.Unix()
		if prev := atomic.LoadInt64(&b.makeDate); current != prev {
			atomic.StoreInt64(&b.makeDate, current)
			b.MQTTPublishAsync(ctx, MQTTTopicMakeDate.Format(snMQTT), 0, true, val)
		}
	} else {
		// TODO: log
	}

	if version, date, err := b.provider.Version(); err == nil {
		current := date.Unix()
		if prev := atomic.LoadInt64(&b.firmwareDate); current != prev {
			atomic.StoreInt64(&b.firmwareDate, current)
			b.MQTTPublishAsync(ctx, MQTTTopicFirmwareDate.Format(snMQTT), 0, true, date)
		}

		b.mutex.Lock()
		if version != b.firmwareVersion {
			b.firmwareVersion = version
			b.MQTTPublishAsync(ctx, MQTTTopicFirmwareVersion.Format(snMQTT), 0, true, version)
		}
		b.mutex.Unlock()
	} else {
		// TODO: log
	}

	return nil, nil
}
