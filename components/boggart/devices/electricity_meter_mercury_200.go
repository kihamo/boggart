package devices

import (
	"context"
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

type Mercury200ElectricityMeter struct {
	boggart.DeviceBase

	provider     *mercury.ElectricityMeter200
	serialNumber string
	interval     time.Duration
}

func NewMercury200ElectricityMeter(serialNumber string, provider *mercury.ElectricityMeter200, interval time.Duration) *Mercury200ElectricityMeter {
	device := &Mercury200ElectricityMeter{
		provider:     provider,
		serialNumber: serialNumber,
		interval:     interval,
	}
	device.Init()
	device.SetDescription("Mercury 200 electricity meter with serial number " + device.serialNumber)

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
	metricElectricityMeterMercuryTariff.With("serial_number", d.serialNumber).Describe(ch)
	metricElectricityMeterMercuryVoltage.With("serial_number", d.serialNumber).Describe(ch)
	metricElectricityMeterMercuryAmperage.With("serial_number", d.serialNumber).Describe(ch)
	metricElectricityMeterMercuryPower.With("serial_number", d.serialNumber).Describe(ch)
	metricElectricityMeterMercuryBatteryVoltage.With("serial_number", d.serialNumber).Describe(ch)
}

func (d *Mercury200ElectricityMeter) Collect(ch chan<- snitch.Metric) {
	metricElectricityMeterMercuryTariff.With("serial_number", d.serialNumber).Collect(ch)
	metricElectricityMeterMercuryVoltage.With("serial_number", d.serialNumber).Collect(ch)
	metricElectricityMeterMercuryAmperage.With("serial_number", d.serialNumber).Collect(ch)
	metricElectricityMeterMercuryPower.With("serial_number", d.serialNumber).Collect(ch)
	metricElectricityMeterMercuryBatteryVoltage.With("serial_number", d.serialNumber).Collect(ch)
}

func (d *Mercury200ElectricityMeter) Ping(_ context.Context) bool {
	_, _, err := d.provider.Version()
	return err == nil
}

func (d *Mercury200ElectricityMeter) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-electricity-meter-mercury-200-updater-" + d.serialNumber)

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

	metricTariff := metricElectricityMeterMercuryTariff.With("serial_number", d.serialNumber)
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff1).Set(tariffs[Mercury200ElectricityMeterTariff1])
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff2).Set(tariffs[Mercury200ElectricityMeterTariff2])
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff3).Set(tariffs[Mercury200ElectricityMeterTariff3])
	metricTariff.With("tariff", Mercury200ElectricityMeterTariff4).Set(tariffs[Mercury200ElectricityMeterTariff4])

	// optimization
	voltage, amperage, power, err := d.provider.ParamsCurrent()
	if err != nil {
		return nil, err
	}
	metricElectricityMeterMercuryVoltage.With("serial_number", d.serialNumber).Set(voltage)
	metricElectricityMeterMercuryAmperage.With("serial_number", d.serialNumber).Set(amperage)
	metricElectricityMeterMercuryPower.With("serial_number", d.serialNumber).Set(float64(power))

	voltage, err = d.BatteryVoltage(ctx)
	if err != nil {
		return nil, nil
	}
	metricElectricityMeterMercuryBatteryVoltage.With("serial_number", d.serialNumber).Set(voltage)

	return nil, nil
}