package internal

import (
	"fmt"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
	"github.com/kihamo/snitch"
)

const (
	MetricMercuryT1             = boggart.ComponentName + "_mercury_t1_watt_hour_total"
	MetricMercuryT2             = boggart.ComponentName + "_mercury_t2_watt_hour_total"
	MetricMercuryT3             = boggart.ComponentName + "_mercury_t3_watt_hour_total"
	MetricMercuryT4             = boggart.ComponentName + "_mercury_t4_watt_hour_total"
	MetricMercuryVoltage        = boggart.ComponentName + "_mercury_voltage_volt"
	MetricMercuryAmperage       = boggart.ComponentName + "_mercury_amperage_ampere"
	MetricMercuryPower          = boggart.ComponentName + "_mercury_power_watt"
	MetricMercuryBatteryVoltage = boggart.ComponentName + "_mercury_battery_voltage_volt"
)

var (
	metricMercuryT1             = snitch.NewGauge(MetricMercuryT1, "Mercury electricity meter T1 value")
	metricMercuryT2             = snitch.NewGauge(MetricMercuryT2, "Mercury electricity meter T2 value")
	metricMercuryT3             = snitch.NewGauge(MetricMercuryT3, "Mercury electricity meter T3 value")
	metricMercuryT4             = snitch.NewGauge(MetricMercuryT4, "Mercury electricity meter T4 value")
	metricMercuryVoltage        = snitch.NewGauge(MetricMercuryVoltage, "Mercury electricity meter current voltage")
	metricMercuryAmperage       = snitch.NewGauge(MetricMercuryAmperage, "Mercury electricity meter current amperage")
	metricMercuryPower          = snitch.NewGauge(MetricMercuryPower, "Mercury electricity meter current power")
	metricMercuryBatteryVoltage = snitch.NewGauge(MetricMercuryBatteryVoltage, "Mercury electricity meter current battery voltage")
)

func (c *MetricsCollector) UpdaterMercury() error {
	device := mercury.NewElectricityMeter200(
		mercury.ConvertSerialNumber(c.component.config.String(boggart.ConfigMercuryDeviceAddress)),
		c.component.ConnectionRS485())

	t1, t2, t3, t4, err := device.PowerCounters()
	if err != nil {
		return fmt.Errorf("PowerCounters error: %s", err.Error())
	}
	metricMercuryT1.Set(float64(t1))
	metricMercuryT2.Set(float64(t2))
	metricMercuryT3.Set(float64(t3))
	metricMercuryT4.Set(float64(t4))

	voltage, amperage, power, err := device.ParamsCurrent()
	if err != nil {
		return fmt.Errorf("ParamsCurrent error: %s", err.Error())
	}
	metricMercuryVoltage.Set(voltage)
	metricMercuryAmperage.Set(amperage)
	metricMercuryPower.Set(float64(power))

	voltage, err = device.BatteryVoltage()
	if err != nil {
		return fmt.Errorf("BatteryVoltage error: %s", err.Error())
	}
	metricMercuryBatteryVoltage.Set(voltage)

	return nil
}

func (c *MetricsCollector) DescribeMercury(ch chan<- *snitch.Description) {
	metricMercuryT1.Describe(ch)
	metricMercuryT2.Describe(ch)
	metricMercuryT3.Describe(ch)
	metricMercuryT4.Describe(ch)
	metricMercuryVoltage.Describe(ch)
	metricMercuryAmperage.Describe(ch)
	metricMercuryPower.Describe(ch)
	metricMercuryBatteryVoltage.Describe(ch)
}

func (c *MetricsCollector) CollectMercury(ch chan<- snitch.Metric) {
	metricMercuryT1.Collect(ch)
	metricMercuryT2.Collect(ch)
	metricMercuryT3.Collect(ch)
	metricMercuryT4.Collect(ch)
	metricMercuryVoltage.Collect(ch)
	metricMercuryAmperage.Collect(ch)
	metricMercuryPower.Collect(ch)
	metricMercuryBatteryVoltage.Collect(ch)
}
