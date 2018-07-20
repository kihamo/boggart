package devices

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
	"gobot.io/x/gobot/drivers/i2c"
)

var (
	metricSensorBME280Temperature = snitch.NewGauge(boggart.ComponentName+"_device_sensor_bme280_temperature_celsius", "Temperature")
	metricSensorBME280Altitude    = snitch.NewGauge(boggart.ComponentName+"_device_sensor_bme280_altitude_metre", "Altitude")
	metricSensorBME280Humidity    = snitch.NewGauge(boggart.ComponentName+"_device_sensor_bme280_humidity_percent", "Humidity")
	metricSensorBME280Pressure    = snitch.NewGauge(boggart.ComponentName+"_device_sensor_bme280_pressure_psi", "Pressure")
)

type BME280SensorChange struct {
	Temperature float64
	Altitude    float64
	Humidity    float64
	Pressure    float64
}

type BME280Sensor struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber

	driver   *i2c.BME280Driver
	interval time.Duration

	mutex      sync.Mutex
	lastValues BME280SensorChange
}

func NewBME280Sensor(connector i2c.Connector, interval time.Duration, bus int, address int) *BME280Sensor {
	driver := i2c.NewBME280Driver(
		connector,
		i2c.WithBus(bus),
		i2c.WithAddress(address))

	device := &BME280Sensor{
		driver:   driver,
		interval: interval,
	}

	device.Init()
	device.SetSerialNumber(fmt.Sprintf("%d_%d", bus, address))
	device.SetDescription("Sensor BME280")

	return device
}

func (d *BME280Sensor) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeThermometer,
		boggart.DeviceTypeBarometer,
		boggart.DeviceTypeHygrometer,
	}
}

func (d *BME280Sensor) Describe(ch chan<- *snitch.Description) {
	metricSensorBME280Temperature.Describe(ch)
	metricSensorBME280Altitude.Describe(ch)
	metricSensorBME280Humidity.Describe(ch)
	metricSensorBME280Pressure.Describe(ch)
}

func (d *BME280Sensor) Collect(ch chan<- snitch.Metric) {
	metricSensorBME280Temperature.Collect(ch)
	metricSensorBME280Altitude.Collect(ch)
	metricSensorBME280Humidity.Collect(ch)
	metricSensorBME280Pressure.Collect(ch)
}

func (d *BME280Sensor) Ping(_ context.Context) bool {
	_, err := d.driver.Temperature()
	return err == nil
}

func (d *BME280Sensor) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-sensor-bme280-updater-" + d.driver.Name())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *BME280Sensor) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	if err := d.driver.Start(); err != nil {
		return nil, err
	}

	serialNumber := d.SerialNumber()
	currentValues := BME280SensorChange{}

	if value, err := d.driver.Temperature(); err == nil {
		currentValues.Temperature = float64(value)
		metricSensorBME280Temperature.With("serial_number", serialNumber).Set(currentValues.Temperature)
	} else {
		return nil, err
	}

	if value, err := d.driver.Altitude(); err == nil {
		currentValues.Altitude = float64(value)
		metricSensorBME280Altitude.With("serial_number", serialNumber).Set(currentValues.Altitude)
	} else {
		return nil, err
	}

	if value, err := d.driver.Humidity(); err == nil {
		currentValues.Humidity = float64(value)
		metricSensorBME280Humidity.With("serial_number", serialNumber).Set(currentValues.Humidity)
	} else {
		return nil, err
	}

	if value, err := d.driver.Pressure(); err == nil {
		currentValues.Pressure = float64(int(value/133.322*10)) / 10
		metricSensorBME280Pressure.With("serial_number", serialNumber).Set(currentValues.Pressure)
	} else {
		return nil, err
	}

	d.mutex.Lock()
	if !reflect.DeepEqual(d.lastValues, currentValues) {
		d.lastValues = currentValues
		d.mutex.Unlock()

		d.TriggerEvent(boggart.DeviceEventBME280Changed, currentValues, serialNumber)
	} else {
		d.mutex.Unlock()
	}

	return nil, nil
}
