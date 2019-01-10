package bind

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"gobot.io/x/gobot/drivers/i2c"
)

const (
	BME280SensorMQTTTopicTemperature mqtt.Topic = boggart.ComponentName + "/meter/bme280/+/temperature"
	BME280SensorMQTTTopicAltitude    mqtt.Topic = boggart.ComponentName + "/meter/bme280/+/altitude"
	BME280SensorMQTTTopicHumidity    mqtt.Topic = boggart.ComponentName + "/meter/bme280/+/humidity"
	BME280SensorMQTTTopicPressure    mqtt.Topic = boggart.ComponentName + "/meter/bme280/+/pressure"
)

type BME280Sensor struct {
	temperature uint64
	altitude    uint64
	humidity    uint64
	pressure    uint64

	boggart.DeviceBindBase
	boggart.DeviceBindSerialNumber
	boggart.DeviceBindMQTT

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

		temperature: math.MaxUint64,
		altitude:    math.MaxUint64,
		humidity:    math.MaxUint64,
		pressure:    math.MaxUint64,
	}

	device.Init()
	device.SetSerialNumber(fmt.Sprintf("%d_%d", bus, address))

	return device
}

func (d *BME280Sensor) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(d.interval)
	taskStateUpdater.SetName("bind-bme280-state-updater-" + d.driver.Name())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (d *BME280Sensor) taskStateUpdater(ctx context.Context) (interface{}, error) {
	if err := d.driver.Start(); err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)
	serialNumber := d.SerialNumber()

	if value, err := d.driver.Temperature(); err == nil {
		current := float64(value)

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

func (d *BME280Sensor) MQTTTopics() []mqtt.Topic {
	sn := d.SerialNumberMQTTEscaped()

	return []mqtt.Topic{
		mqtt.Topic(BME280SensorMQTTTopicTemperature.Format(sn)),
		mqtt.Topic(BME280SensorMQTTTopicAltitude.Format(sn)),
		mqtt.Topic(BME280SensorMQTTTopicHumidity.Format(sn)),
		mqtt.Topic(BME280SensorMQTTTopicPressure.Format(sn)),
	}
}
