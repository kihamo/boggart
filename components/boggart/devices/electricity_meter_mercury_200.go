package devices

import (
	"context"
	"math"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

const (
	Mercury200ElectricityMeterTariff1 = "t1"
	Mercury200ElectricityMeterTariff2 = "t2"
	Mercury200ElectricityMeterTariff3 = "t3"
	Mercury200ElectricityMeterTariff4 = "t4"

	Mercury200ElectricityMeterMQTTTopicTariff1        mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/tariff_1"
	Mercury200ElectricityMeterMQTTTopicTariff2        mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/tariff_2"
	Mercury200ElectricityMeterMQTTTopicTariff3        mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/tariff_3"
	Mercury200ElectricityMeterMQTTTopicTariff4        mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/tariff_4"
	Mercury200ElectricityMeterMQTTTopicVoltage        mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/voltage"
	Mercury200ElectricityMeterMQTTTopicAmperage       mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/amperage"
	Mercury200ElectricityMeterMQTTTopicPower          mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/power"
	Mercury200ElectricityMeterMQTTTopicBatteryVoltage mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/battery_voltage"
)

var (
	metricElectricityMeterMercuryTariff         = snitch.NewGauge(boggart.ComponentName+"_device_electricity_meter_mercury_200_tariff_watt_hour_total", "Mercury 200 electricity meter tariff value")
	metricElectricityMeterMercuryVoltage        = snitch.NewGauge(boggart.ComponentName+"_device_electricity_meter_mercury_200_voltage_volt", "Mercury 200 electricity meter current voltage")
	metricElectricityMeterMercuryAmperage       = snitch.NewGauge(boggart.ComponentName+"_device_electricity_meter_mercury_200_amperage_ampere", "Mercury 200 electricity meter current amperage")
	metricElectricityMeterMercuryPower          = snitch.NewGauge(boggart.ComponentName+"_device_electricity_meter_mercury_200_power_watt", "Mercury 200 electricity meter current power")
	metricElectricityMeterMercuryBatteryVoltage = snitch.NewGauge(boggart.ComponentName+"_device_electricity_meter_mercury_200_battery_voltage_volt", "Mercury 200 electricity meter current battery voltage")
)

type Mercury200ElectricityMeter struct {
	tariff1        uint64
	tariff2        uint64
	tariff3        uint64
	tariff4        uint64
	voltage        uint64
	amperage       uint64
	power          int64
	batteryVoltage uint64

	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	provider *mercury.ElectricityMeter200
	interval time.Duration
}

func NewMercury200ElectricityMeter(serialNumber string, provider *mercury.ElectricityMeter200, interval time.Duration) *Mercury200ElectricityMeter {
	device := &Mercury200ElectricityMeter{
		provider: provider,
		interval: interval,
	}
	device.Init()
	device.SetSerialNumber(serialNumber)
	device.SetDescription("Mercury 200 electricity meter with serial number " + serialNumber)

	return device
}

func (d *Mercury200ElectricityMeter) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeElectricityMeter,
	}
}

func (d *Mercury200ElectricityMeter) Tariffs(_ context.Context) (map[string]float64, error) {
	t1, t2, t3, t4, err := d.provider.PowerCounters()
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		Mercury200ElectricityMeterTariff1: float64(t1),
		Mercury200ElectricityMeterTariff2: float64(t2),
		Mercury200ElectricityMeterTariff3: float64(t3),
		Mercury200ElectricityMeterTariff4: float64(t4),
	}, nil
}

func (d *Mercury200ElectricityMeter) Voltage(_ context.Context) (float64, error) {
	voltage, _, _, err := d.provider.ParamsCurrent()
	if err != nil {
		return -1, err
	}

	return voltage, nil
}

func (d *Mercury200ElectricityMeter) Amperage(_ context.Context) (float64, error) {
	_, amperage, _, err := d.provider.ParamsCurrent()
	if err != nil {
		return -1, err
	}

	return amperage, nil
}

func (d *Mercury200ElectricityMeter) Power(_ context.Context) (float64, error) {
	_, _, power, err := d.provider.ParamsCurrent()
	if err != nil {
		return -1, err
	}

	return float64(power), nil
}

func (d *Mercury200ElectricityMeter) BatteryVoltage(_ context.Context) (float64, error) {
	return d.provider.BatteryVoltage()
}

func (d *Mercury200ElectricityMeter) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()

	metricElectricityMeterMercuryTariff.With("serial_number", serialNumber).Describe(ch)
	metricElectricityMeterMercuryVoltage.With("serial_number", serialNumber).Describe(ch)
	metricElectricityMeterMercuryAmperage.With("serial_number", serialNumber).Describe(ch)
	metricElectricityMeterMercuryPower.With("serial_number", serialNumber).Describe(ch)
	metricElectricityMeterMercuryBatteryVoltage.With("serial_number", serialNumber).Describe(ch)
}

func (d *Mercury200ElectricityMeter) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()

	metricElectricityMeterMercuryTariff.With("serial_number", serialNumber).Collect(ch)
	metricElectricityMeterMercuryVoltage.With("serial_number", serialNumber).Collect(ch)
	metricElectricityMeterMercuryAmperage.With("serial_number", serialNumber).Collect(ch)
	metricElectricityMeterMercuryPower.With("serial_number", serialNumber).Collect(ch)
	metricElectricityMeterMercuryBatteryVoltage.With("serial_number", serialNumber).Collect(ch)
}

func (d *Mercury200ElectricityMeter) Ping(_ context.Context) bool {
	_, _, err := d.provider.Version()
	return err == nil
}

func (d *Mercury200ElectricityMeter) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-electricity-meter-mercury-200-updater-" + d.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *Mercury200ElectricityMeter) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	tariffs, err := d.Tariffs(ctx)
	if err != nil {
		return nil, err
	}

	serialNumber := d.SerialNumber()

	metricTariff := metricElectricityMeterMercuryTariff.With("serial_number", serialNumber)
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff1).Set(tariffs[Mercury200ElectricityMeterTariff1])
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff2).Set(tariffs[Mercury200ElectricityMeterTariff2])
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff3).Set(tariffs[Mercury200ElectricityMeterTariff3])
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff4).Set(tariffs[Mercury200ElectricityMeterTariff4])

	prevTariff1 := math.Float64frombits(atomic.LoadUint64(&d.tariff1))
	if tariffs[Mercury200ElectricityMeterTariff1] != prevTariff1 {
		atomic.StoreUint64(&d.tariff1, math.Float64bits(tariffs[Mercury200ElectricityMeterTariff1]))

		d.MQTTPublishAsync(ctx, Mercury200ElectricityMeterMQTTTopicTariff1.Format(serialNumber), 0, true, tariffs[Mercury200ElectricityMeterTariff1])
	}

	prevTariff2 := math.Float64frombits(atomic.LoadUint64(&d.tariff2))
	if tariffs[Mercury200ElectricityMeterTariff2] != prevTariff2 {
		atomic.StoreUint64(&d.tariff2, math.Float64bits(tariffs[Mercury200ElectricityMeterTariff2]))

		d.MQTTPublishAsync(ctx, Mercury200ElectricityMeterMQTTTopicTariff2.Format(serialNumber), 0, true, tariffs[Mercury200ElectricityMeterTariff2])
	}

	prevTariff3 := math.Float64frombits(atomic.LoadUint64(&d.tariff3))
	if tariffs[Mercury200ElectricityMeterTariff3] != prevTariff3 {
		atomic.StoreUint64(&d.tariff3, math.Float64bits(tariffs[Mercury200ElectricityMeterTariff3]))

		d.MQTTPublishAsync(ctx, Mercury200ElectricityMeterMQTTTopicTariff3.Format(serialNumber), 0, true, tariffs[Mercury200ElectricityMeterTariff3])
	}

	prevTariff4 := math.Float64frombits(atomic.LoadUint64(&d.tariff4))
	if tariffs[Mercury200ElectricityMeterTariff4] != prevTariff4 {
		atomic.StoreUint64(&d.tariff4, math.Float64bits(tariffs[Mercury200ElectricityMeterTariff4]))

		d.MQTTPublishAsync(ctx, Mercury200ElectricityMeterMQTTTopicTariff4.Format(serialNumber), 0, true, tariffs[Mercury200ElectricityMeterTariff4])
	}

	// optimization
	voltage, amperage, power, err := d.provider.ParamsCurrent()
	if err != nil {
		return nil, err
	}

	metricElectricityMeterMercuryVoltage.With("serial_number", serialNumber).Set(voltage)
	metricElectricityMeterMercuryAmperage.With("serial_number", serialNumber).Set(amperage)
	metricElectricityMeterMercuryPower.With("serial_number", serialNumber).Set(float64(power))

	prevVoltage := math.Float64frombits(atomic.LoadUint64(&d.voltage))
	if voltage != prevVoltage {
		atomic.StoreUint64(&d.voltage, math.Float64bits(voltage))

		d.MQTTPublishAsync(ctx, Mercury200ElectricityMeterMQTTTopicVoltage.Format(serialNumber), 0, true, voltage)
	}

	prevAmperage := math.Float64frombits(atomic.LoadUint64(&d.amperage))
	if amperage != prevAmperage {
		atomic.StoreUint64(&d.amperage, math.Float64bits(amperage))

		d.MQTTPublishAsync(ctx, Mercury200ElectricityMeterMQTTTopicAmperage.Format(serialNumber), 0, true, amperage)
	}

	prevPower := atomic.LoadInt64(&d.power)
	if power != prevPower {
		atomic.StoreInt64(&d.power, power)

		d.MQTTPublishAsync(ctx, Mercury200ElectricityMeterMQTTTopicPower.Format(serialNumber), 0, true, power)
	}

	batteryVoltage, err := d.BatteryVoltage(ctx)
	if err != nil {
		return nil, nil
	}

	metricElectricityMeterMercuryBatteryVoltage.With("serial_number", serialNumber).Set(batteryVoltage)

	prevBatteryVoltage := math.Float64frombits(atomic.LoadUint64(&d.batteryVoltage))
	if batteryVoltage != prevBatteryVoltage {
		atomic.StoreUint64(&d.batteryVoltage, math.Float64bits(batteryVoltage))

		d.MQTTPublishAsync(ctx, Mercury200ElectricityMeterMQTTTopicBatteryVoltage.Format(serialNumber), 0, true, batteryVoltage)
	}

	return nil, nil
}

func (d *Mercury200ElectricityMeter) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		Mercury200ElectricityMeterMQTTTopicTariff1,
		Mercury200ElectricityMeterMQTTTopicTariff2,
		Mercury200ElectricityMeterMQTTTopicTariff3,
		Mercury200ElectricityMeterMQTTTopicTariff4,
		Mercury200ElectricityMeterMQTTTopicVoltage,
		Mercury200ElectricityMeterMQTTTopicAmperage,
		Mercury200ElectricityMeterMQTTTopicPower,
		Mercury200ElectricityMeterMQTTTopicBatteryVoltage,
	}
}
