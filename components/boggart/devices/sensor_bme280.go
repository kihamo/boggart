package devices

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
	"gobot.io/x/gobot/drivers/i2c"
)

const (
	BME280SensorMQTTTopicTemperature boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/bme280/+/temperature"
	BME280SensorMQTTTopicAltitude    boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/bme280/+/altitude"
	BME280SensorMQTTTopicHumidity    boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/bme280/+/humidity"
	BME280SensorMQTTTopicPressure    boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/bme280/+/pressure"
)

var (
	metricSensorBME280Temperature = snitch.NewGauge(boggart.ComponentName+"_device_sensor_bme280_temperature_celsius", "Temperature")
	metricSensorBME280Altitude    = snitch.NewGauge(boggart.ComponentName+"_device_sensor_bme280_altitude_metre", "Altitude")
	metricSensorBME280Humidity    = snitch.NewGauge(boggart.ComponentName+"_device_sensor_bme280_humidity_percent", "Humidity")
	metricSensorBME280Pressure    = snitch.NewGauge(boggart.ComponentName+"_device_sensor_bme280_pressure_psi", "Pressure")
)

type BME280Sensor struct {
	temperature uint64
	altitude    uint64
	humidity    uint64
	pressure    uint64

	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	driver   *i2c.BME280Driver
	interval time.Duration
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

	if value, err := d.driver.Temperature(); err == nil {
		current := float64(value)

		metricSensorBME280Temperature.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.temperature))
		if current != prev {
			atomic.StoreUint64(&d.temperature, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, BME280SensorMQTTTopicTemperature.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	if value, err := d.driver.Altitude(); err == nil {
		current := float64(value)

		metricSensorBME280Altitude.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.altitude))
		if current != prev {
			atomic.StoreUint64(&d.altitude, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, BME280SensorMQTTTopicAltitude.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	if value, err := d.driver.Humidity(); err == nil {
		current := float64(value)

		metricSensorBME280Humidity.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.humidity))
		if current != prev {
			atomic.StoreUint64(&d.humidity, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, BME280SensorMQTTTopicHumidity.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	if value, err := d.driver.Pressure(); err == nil {
		current := float64(int(value/133.322*10)) / 10

		metricSensorBME280Pressure.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.pressure))
		if current != prev {
			atomic.StoreUint64(&d.pressure, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, BME280SensorMQTTTopicPressure.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	return nil, nil
}

func (d *BME280Sensor) MQTTTopics() []boggart.DeviceMQTTTopic {
	return []boggart.DeviceMQTTTopic{
		BME280SensorMQTTTopicTemperature,
		BME280SensorMQTTTopicAltitude,
		BME280SensorMQTTTopicHumidity,
		BME280SensorMQTTTopicPressure,
	}
}
