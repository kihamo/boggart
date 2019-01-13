package nut

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricLoad           = snitch.NewGauge(boggart.ComponentName+"_bind_nut_load_percent", "Load on UPS (percent of full)")
	metricInputVoltage   = snitch.NewGauge(boggart.ComponentName+"_bind_nut_input_voltage_volts", "Input voltage volts")
	metricBatteryCharge  = snitch.NewGauge(boggart.ComponentName+"_bind_nut_battery_charge_percent", "Battery charge (percent of full)")
	metricBatteryRuntime = snitch.NewGauge(boggart.ComponentName+"_bind_nut_battery_runtime_seconds", "Battery runtime seconds")
	metricBatteryVoltage = snitch.NewGauge(boggart.ComponentName+"_bind_nut_battery_voltage_volts", "Battery voltage volts")
)

/*
battery.charge Battery charge (percent of full)
battery.charge.low Remaining battery level when UPS switches to LB (percent)
battery.charge.warning Battery level when UPS switches to Warning state (percent)
battery.date Battery change date
battery.mfr.date Battery manufacturing date
battery.runtime Battery runtime (seconds)
battery.runtime.low Remaining battery runtime when UPS switches to LB (seconds)
battery.type Battery chemistry
battery.voltage Battery voltage (V)
battery.voltage.nominal Nominal battery voltage (V)
device.mfr Description unavailable
device.model Description unavailable
device.serial Description unavailable
device.type Description unavailable
driver.name Driver name
driver.parameter.pollfreq Description unavailable
driver.parameter.pollinterval Description unavailable
driver.parameter.port Description unavailable
driver.parameter.productid Description unavailable
driver.parameter.synchronous Description unavailable
driver.parameter.vendorid Description unavailable
driver.version Driver version - NUT release
driver.version.data Description unavailable
driver.version.internal Internal driver version
input.sensitivity Input power sensitivity
input.transfer.high High voltage transfer point (V)
input.transfer.low Low voltage transfer point (V)
input.voltage Input voltage (V)
input.voltage.nominal Nominal input voltage (V)
ups.beeper.status UPS beeper status
ups.delay.shutdown Interval to wait after shutdown with delay command (seconds)
ups.firmware UPS firmware
ups.firmware.aux Auxiliary device firmware
ups.load Load on UPS (percent of full)
ups.mfr UPS manufacturer
ups.mfr.date UPS manufacturing date
ups.model UPS model
ups.productid Product ID for USB devices
ups.realpower.nominal UPS real power rating (W)
ups.serial UPS serial number
ups.status UPS status
ups.test.result Results of last self test
ups.timer.reboot Time before the load will be rebooted (seconds)
ups.timer.shutdown Time before the load will be shutdown (seconds)
ups.vendorid Vendor ID for USB devices
*/

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	metricLoad.With("serial_number", sn).Describe(ch)
	metricInputVoltage.With("serial_number", sn).Describe(ch)
	metricBatteryCharge.With("serial_number", sn).Describe(ch)
	metricBatteryRuntime.With("serial_number", sn).Describe(ch)
	metricBatteryVoltage.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	metricLoad.With("serial_number", sn).Collect(ch)
	metricInputVoltage.With("serial_number", sn).Collect(ch)
	metricBatteryCharge.With("serial_number", sn).Collect(ch)
	metricBatteryRuntime.With("serial_number", sn).Collect(ch)
	metricBatteryVoltage.With("serial_number", sn).Collect(ch)
}
