package bind

import (
	"context"
	"encoding/hex"
	"math"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

const (
	PulsarHeatMeterInputScale = 1000

	PulsarHeatMeterMQTTTopicTemperatureIn    mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_in"
	PulsarHeatMeterMQTTTopicTemperatureOut   mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_out"
	PulsarHeatMeterMQTTTopicTemperatureDelta mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_delta"
	PulsarHeatMeterMQTTTopicEnergy           mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/energy"
	PulsarHeatMeterMQTTTopicConsumption      mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/consumption"
	PulsarHeatMeterMQTTTopicCapacity         mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/capacity"
	PulsarHeatMeterMQTTTopicPower            mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/power"
	PulsarHeatMeterMQTTTopicInputPulses      mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/input/+/pulses"
	PulsarHeatMeterMQTTTopicInputVolume      mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/input/+/volume"
)

type PulsarHeatMeter struct {
	temperatureIn    uint64
	temperatureOut   uint64
	temperatureDelta uint64
	energy           uint64
	consumption      uint64
	capacity         uint64
	power            uint64
	input1           uint64
	input2           uint64
	input3           uint64
	input4           uint64

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	config   *PulsarHeatMeterConfig
	provider *pulsar.HeatMeter
}

type PulsarHeatMeterConfig struct {
	RS485 struct {
		Address string `valid:"required"`
		Timeout string
	} `valid:"required"`
	Address      string
	Input1Offset float64 `mapstructure:"input1_offset",valid:"float"`
	Input2Offset float64 `mapstructure:"input2_offset",valid:"float"`
	Input3Offset float64 `mapstructure:"input3_offset",valid:"float"`
	Input4Offset float64 `mapstructure:"input4_offset",valid:"float"`
}

func (d PulsarHeatMeter) Config() interface{} {
	return &PulsarHeatMeterConfig{}
}

func (d PulsarHeatMeter) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*PulsarHeatMeterConfig)

	var err error
	timeout := time.Second

	if config.RS485.Timeout != "" {
		timeout, err = time.ParseDuration(config.RS485.Timeout)
		if err != nil {
			return nil, err
		}
	}

	conn := rs485.GetConnection(config.RS485.Address, timeout)

	var deviceAddress []byte
	if config.Address == "" {
		deviceAddress, err = pulsar.DeviceAddress(conn)
	} else {
		deviceAddress, err = hex.DecodeString(config.Address)
	}

	if err != nil {
		return nil, err
	}

	device := &PulsarHeatMeter{
		config:   config,
		provider: pulsar.NewHeatMeter(deviceAddress, conn),

		temperatureIn:    math.MaxUint64,
		temperatureOut:   math.MaxUint64,
		temperatureDelta: math.MaxUint64,
		energy:           math.MaxUint64,
		consumption:      math.MaxUint64,
		capacity:         math.MaxUint64,
		power:            math.MaxUint64,
		input1:           math.MaxUint64,
		input2:           math.MaxUint64,
		input3:           math.MaxUint64,
		input4:           math.MaxUint64,
	}
	device.Init()
	device.SetSerialNumber(hex.EncodeToString(deviceAddress))

	return device, nil
}

func (d *PulsarHeatMeter) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(time.Minute)
	taskStateUpdater.SetName("bind-pulsar-heat-meter-state-updater-" + d.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (d *PulsarHeatMeter) taskStateUpdater(ctx context.Context) (interface{}, error) {
	if _, err := d.provider.Version(); err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)
	serialNumber := d.SerialNumber()

	if currentVal, err := d.provider.TemperatureIn(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&d.temperatureIn))
		if current != prev {
			atomic.StoreUint64(&d.temperatureIn, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicTemperatureIn.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := d.provider.TemperatureOut(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&d.temperatureOut))
		if current != prev {
			atomic.StoreUint64(&d.temperatureOut, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicTemperatureOut.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := d.provider.TemperatureDelta(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&d.temperatureDelta))
		if current != prev {
			atomic.StoreUint64(&d.temperatureDelta, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicTemperatureDelta.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := d.provider.Energy(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&d.energy))
		if current != prev {
			atomic.StoreUint64(&d.energy, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicEnergy.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := d.provider.Consumption(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&d.consumption))
		if current != prev {
			atomic.StoreUint64(&d.consumption, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicConsumption.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := d.provider.Capacity(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&d.capacity))
		if current != prev {
			atomic.StoreUint64(&d.capacity, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicCapacity.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := d.provider.Power(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&d.power))
		if current != prev {
			atomic.StoreUint64(&d.power, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicPower.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	// inputs
	if currentVal, err := d.provider.PulseInput1(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&d.input1)
		if current != prev {
			atomic.StoreUint64(&d.input1, current)

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicInputPulses.Format(serialNumber, 1), 0, true, current)
			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicInputVolume.Format(serialNumber, 1), 0, true, d.inputVolume(current, d.config.Input1Offset))
		}
	} else {
		// TODO: log
	}

	if currentVal, err := d.provider.PulseInput2(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&d.input2)
		if current != prev {
			atomic.StoreUint64(&d.input2, current)

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicInputPulses.Format(serialNumber, 2), 0, true, current)
			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicInputVolume.Format(serialNumber, 2), 0, true, d.inputVolume(current, d.config.Input2Offset))
		}
	} else {
		// TODO: log
	}

	if currentVal, err := d.provider.PulseInput3(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&d.input3)
		if current != prev {
			atomic.StoreUint64(&d.input3, current)

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicInputPulses.Format(serialNumber, 3), 0, true, current)
			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicInputVolume.Format(serialNumber, 3), 0, true, d.inputVolume(current, d.config.Input3Offset))
		}
	} else {
		// TODO: log
	}

	if currentVal, err := d.provider.PulseInput4(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&d.input4)
		if current != prev {
			atomic.StoreUint64(&d.input4, current)

			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicInputPulses.Format(serialNumber, 4), 0, true, current)
			d.MQTTPublishAsync(ctx, PulsarHeatMeterMQTTTopicInputVolume.Format(serialNumber, 4), 0, true, d.inputVolume(current, d.config.Input4Offset))
		}
	} else {
		// TODO: log
	}

	return nil, nil
}

func (d *PulsarHeatMeter) inputVolume(pulses uint64, offset float64) float64 {
	return (offset*PulsarHeatMeterInputScale + float64(pulses*10)) / PulsarHeatMeterInputScale
}

func (d *PulsarHeatMeter) MQTTTopics() []mqtt.Topic {
	sn := d.SerialNumber()

	return []mqtt.Topic{
		mqtt.Topic(PulsarHeatMeterMQTTTopicTemperatureIn.Format(sn)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicTemperatureOut.Format(sn)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicTemperatureDelta.Format(sn)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicEnergy.Format(sn)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicConsumption.Format(sn)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicCapacity.Format(sn)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicPower.Format(sn)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicInputPulses.Format(sn, 1)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicInputVolume.Format(sn, 1)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicInputPulses.Format(sn, 2)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicInputVolume.Format(sn, 2)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicInputPulses.Format(sn, 3)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicInputVolume.Format(sn, 3)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicInputPulses.Format(sn, 4)),
		mqtt.Topic(PulsarHeatMeterMQTTTopicInputVolume.Format(sn, 4)),
	}
}
