package devices

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/apcupsd"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

var (
	metricUPSApcupsdLineVoltage             = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_line_voltage_volt", "The current line voltage as returned by the UPS")
	metricUPSApcupsdLoadPercent             = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_load_percent", "The percentage of load capacity as estimated by the UPS")
	metricUPSApcupsdBatteryChargePercent    = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_battery_charge_percent", "The percentage charge on the batteries")
	metricUPSApcupsdTimeLeft                = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_time_left_seconds", "The remaining runtime left on batteries as estimated by the UPS")
	metricUPSApcupsdOutputVoltage           = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_output_voltage_volt", "The voltage the UPS is supplying to your equipment")
	metricUPSApcupsdInternalTemp            = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_internal_temp_celsius", "Internal UPS temperature as supplied by the UPS")
	metricUPSApcupsdBatteryVoltage          = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_battery_voltage_volt", "Battery voltage as supplied by the UPS")
	metricUPSApcupsdLineFrequency           = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_line_frequency_hertz", "Line frequency in hertz as given by the UPS")
	metricUPSApcupsdTransfers               = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_transfers_total", "The number of transfers to batteries since apcupsd startup")
	metricUPSApcupsdTimeOnBattery           = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_time_on_battery_seconds", "Time in seconds currently on batteries")
	metricUPSApcupsdCumulativeTimeOnBattery = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_cumulative_time_on_battery_seconds", "Total (cumulative) time on batteries in seconds since apcupsd startup")
	metricUPSApcupsdNominalBatteryVoltage   = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_nominal_battery_voltage_volt", "The nominal battery voltage")
	metricUPSApcupsdHumidity                = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_humidity_percent", "The humidity as measured by the UPS")
	metricUPSApcupsdAmbientTemperature      = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_ambient_temperature_celsius", "The ambient temperature as measured by the UPS")
	metricUPSApcupsdBadBatteryPacks         = snitch.NewGauge(boggart.ComponentName+"_device_ups_apcupsd_bad_battery_packs", "The number of bad battery packs")
)

type ApcupsdUPS struct {
	boggart.DeviceWithSerialNumber

	client   *apcupsd.Client
	interval time.Duration
}

func NewApcupsdUPS(client *apcupsd.Client, interval time.Duration) *ApcupsdUPS {
	device := &ApcupsdUPS{
		client:   client,
		interval: interval,
	}
	device.Init()
	device.SetDescription("UPS")

	return device
}

func (d *ApcupsdUPS) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeUPS,
	}
}

func (d *ApcupsdUPS) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricUPSApcupsdLineVoltage.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdLoadPercent.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdBatteryChargePercent.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdTimeLeft.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdOutputVoltage.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdInternalTemp.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdBatteryVoltage.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdLineFrequency.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdTransfers.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdTimeOnBattery.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdCumulativeTimeOnBattery.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdNominalBatteryVoltage.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdHumidity.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdAmbientTemperature.With("serial_number", serialNumber).Describe(ch)
	metricUPSApcupsdBadBatteryPacks.With("serial_number", serialNumber).Describe(ch)
}

func (d *ApcupsdUPS) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}

	metricUPSApcupsdLineVoltage.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdLoadPercent.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdBatteryChargePercent.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdTimeLeft.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdOutputVoltage.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdInternalTemp.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdBatteryVoltage.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdLineFrequency.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdTransfers.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdTimeOnBattery.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdCumulativeTimeOnBattery.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdNominalBatteryVoltage.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdHumidity.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdAmbientTemperature.With("serial_number", serialNumber).Collect(ch)
	metricUPSApcupsdBadBatteryPacks.With("serial_number", serialNumber).Collect(ch)
}

func (d *ApcupsdUPS) Ping(ctx context.Context) bool {
	status, err := d.client.Status(ctx)
	if err == nil && status.Status != nil {
		return *status.Status == "ONLINE"
	}

	return false
}

func (d *ApcupsdUPS) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillStopTask(d.taskSerialNumber)
	taskSerialNumber.SetTimeout(time.Second * 5)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Minute)
	taskSerialNumber.SetName("device-ups-apcupsd-serial-number")

	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-ups-apcupsd-updater-" + d.Id())

	return []workers.Task{
		taskSerialNumber,
		taskUpdater,
	}
}

func (d *ApcupsdUPS) taskSerialNumber(ctx context.Context) (interface{}, error, bool) {
	if !d.IsEnabled() {
		return nil, nil, false
	}

	status, err := d.client.Status(ctx)
	if err != nil || status.SerialNumber == nil {
		return nil, err, false
	}

	d.SetSerialNumber(*status.SerialNumber)
	d.SetDescription("UPS with serial number " + *status.SerialNumber)

	return nil, nil, true
}

func (d *ApcupsdUPS) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return nil, nil
	}

	status, err := d.client.Status(ctx)
	if err != nil {
		return nil, err
	}

	if status.LineVoltage != nil {
		metricUPSApcupsdLineVoltage.With("serial_number", serialNumber).Set(*status.LineVoltage)
	}

	if status.LoadPercent != nil {
		metricUPSApcupsdLoadPercent.With("serial_number", serialNumber).Set(*status.LoadPercent)
	}

	if status.BatteryChargePercent != nil {
		metricUPSApcupsdBatteryChargePercent.With("serial_number", serialNumber).Set(*status.BatteryChargePercent)
	}

	if status.TimeLeft != nil {
		metricUPSApcupsdTimeLeft.With("serial_number", serialNumber).Set(status.TimeLeft.Seconds())
	}

	if status.OutputVoltage != nil {
		metricUPSApcupsdOutputVoltage.With("serial_number", serialNumber).Set(*status.OutputVoltage)
	}

	if status.InternalTemp != nil {
		metricUPSApcupsdInternalTemp.With("serial_number", serialNumber).Set(*status.InternalTemp)
	}

	if status.BatteryVoltage != nil {
		metricUPSApcupsdBatteryVoltage.With("serial_number", serialNumber).Set(*status.BatteryVoltage)
	}

	if status.LineFrequency != nil {
		metricUPSApcupsdLineFrequency.With("serial_number", serialNumber).Set(*status.LineFrequency)
	}

	if status.Transfers != nil {
		metricUPSApcupsdTransfers.With("serial_number", serialNumber).Set(float64(*status.Transfers))
	}

	if status.TimeOnBattery != nil {
		metricUPSApcupsdTimeOnBattery.With("serial_number", serialNumber).Set(status.TimeOnBattery.Seconds())
	}

	if status.CumulativeTimeOnBattery != nil {
		metricUPSApcupsdCumulativeTimeOnBattery.With("serial_number", serialNumber).Set(status.CumulativeTimeOnBattery.Seconds())
	}

	if status.NominalBatteryVoltage != nil {
		metricUPSApcupsdNominalBatteryVoltage.With("serial_number", serialNumber).Set(*status.NominalBatteryVoltage)
	}

	if status.Humidity != nil {
		metricUPSApcupsdHumidity.With("serial_number", serialNumber).Set(*status.Humidity)
	}

	if status.AmbientTemperature != nil {
		metricUPSApcupsdAmbientTemperature.With("serial_number", serialNumber).Set(*status.AmbientTemperature)
	}

	if status.BadBatteryPacks != nil {
		metricUPSApcupsdBadBatteryPacks.With("serial_number", serialNumber).Set(float64(*status.BadBatteryPacks))
	}

	// TODO: автоматически подстроится под интервал обновления на apcupsd, что бы лишний раз не гонять тикет

	return nil, nil
}
