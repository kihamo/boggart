package devices

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

const (
	Mercury200ElectricityMeterTariff1 = "t1"
	Mercury200ElectricityMeterTariff2 = "t2"
	Mercury200ElectricityMeterTariff3 = "t3"
	Mercury200ElectricityMeterTariff4 = "t4"
)

var (
	metricElectricityMeterMercuryTariff         = snitch.NewGauge(boggart.ComponentName+"_device_electricity_meter_mercury_200_tariff_watt_hour_total", "Mercury 200 electricity meter tariff value")
	metricElectricityMeterMercuryVoltage        = snitch.NewGauge(boggart.ComponentName+"_device_electricity_meter_mercury_200_voltage_volt", "Mercury 200 electricity meter current voltage")
	metricElectricityMeterMercuryAmperage       = snitch.NewGauge(boggart.ComponentName+"_device_electricity_meter_mercury_200_amperage_ampere", "Mercury 200 electricity meter current amperage")
	metricElectricityMeterMercuryPower          = snitch.NewGauge(boggart.ComponentName+"_device_electricity_meter_mercury_200_power_watt", "Mercury 200 electricity meter current power")
	metricElectricityMeterMercuryBatteryVoltage = snitch.NewGauge(boggart.ComponentName+"_device_electricity_meter_mercury_200_battery_voltage_volt", "Mercury 200 electricity meter current battery voltage")
)

type Mercury200ElectricityMeterChange struct {
	Tariff1        float64
	Tariff2        float64
	Tariff3        float64
	Tariff4        float64
	Voltage        float64
	Amperage       float64
	Power          int64
	BatteryVoltage float64
}

type Mercury200ElectricityMeter struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber

	provider *mercury.ElectricityMeter200
	interval time.Duration

	mutex      sync.Mutex
	lastValues Mercury200ElectricityMeterChange
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
	currentValues := Mercury200ElectricityMeterChange{
		Tariff1: tariffs[Mercury200ElectricityMeterTariff1],
		Tariff2: tariffs[Mercury200ElectricityMeterTariff2],
		Tariff3: tariffs[Mercury200ElectricityMeterTariff3],
		Tariff4: tariffs[Mercury200ElectricityMeterTariff4],
	}

	metricTariff := metricElectricityMeterMercuryTariff.With("serial_number", serialNumber)
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff1).Set(currentValues.Tariff1)
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff2).Set(currentValues.Tariff2)
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff3).Set(currentValues.Tariff3)
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff4).Set(currentValues.Tariff4)

	// optimization
	voltage, amperage, power, err := d.provider.ParamsCurrent()
	if err != nil {
		return nil, err
	}
	currentValues.Voltage = voltage
	currentValues.Amperage = amperage
	currentValues.Power = power

	metricElectricityMeterMercuryVoltage.With("serial_number", serialNumber).Set(voltage)
	metricElectricityMeterMercuryAmperage.With("serial_number", serialNumber).Set(amperage)
	metricElectricityMeterMercuryPower.With("serial_number", serialNumber).Set(float64(power))

	voltage, err = d.BatteryVoltage(ctx)
	if err != nil {
		return nil, nil
	}
	currentValues.BatteryVoltage = voltage
	metricElectricityMeterMercuryBatteryVoltage.With("serial_number", serialNumber).Set(voltage)

	d.mutex.Lock()
	if !reflect.DeepEqual(d.lastValues, currentValues) {
		d.lastValues = currentValues
		d.mutex.Unlock()

		d.TriggerEvent(ctx, boggart.DeviceEventMercury200Changed, currentValues, serialNumber)
	} else {
		d.mutex.Unlock()
	}

	return nil, nil
}
