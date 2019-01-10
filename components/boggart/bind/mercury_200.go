package bind

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

const (
	Mercury200MQTTTopicTariff          mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/tariff/+"
	Mercury200MQTTTopicVoltage         mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/voltage"
	Mercury200MQTTTopicAmperage        mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/amperage"
	Mercury200MQTTTopicPower           mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/power"
	Mercury200MQTTTopicBatteryVoltage  mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/battery_voltage"
	Mercury200MQTTTopicLastPowerOff    mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/last-power-off"
	Mercury200MQTTTopicLastPowerOn     mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/last-power-on"
	Mercury200MQTTTopicMakeDate        mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/make-date"
	Mercury200MQTTTopicFirmwareDate    mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/firmware/date"
	Mercury200MQTTTopicFirmwareVersion mqtt.Topic = boggart.ComponentName + "/meter/mercury200/+/firmware/version"
)

type Mercury200 struct {
	tariff1         uint64
	tariff2         uint64
	tariff3         uint64
	tariff4         uint64
	voltage         uint64
	amperage        uint64
	power           uint64
	batteryVoltage  uint64
	lastPowerOff    int64
	lastPowerOn     int64
	makeDate        int64
	firmwareDate    int64
	firmwareVersion string

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	mutex    sync.Mutex
	provider *mercury.ElectricityMeter200
}

type Mercury200Config struct {
	RS485 struct {
		Address string `valid:"required"`
		Timeout string
	} `valid:"required"`
	Address string `valid:"required"`
}

func (d Mercury200) Config() interface{} {
	return &Mercury200Config{}
}

func (d Mercury200) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Mercury200Config)

	var err error
	timeout := time.Second

	if config.RS485.Timeout != "" {
		timeout, err = time.ParseDuration(config.RS485.Timeout)
		if err != nil {
			return nil, err
		}
	}

	provider := mercury.NewElectricityMeter200(
		mercury.ConvertSerialNumber(config.Address),
		rs485.GetConnection(config.RS485.Address, timeout))

	device := &Mercury200{
		provider: provider,

		tariff1:        math.MaxUint64,
		tariff2:        math.MaxUint64,
		tariff3:        math.MaxUint64,
		tariff4:        math.MaxUint64,
		voltage:        math.MaxUint64,
		amperage:       math.MaxUint64,
		power:          math.MaxInt64,
		batteryVoltage: math.MaxUint64,
	}
	device.Init()

	// TODO: read real serial number
	device.SetSerialNumber(config.Address)

	// TODO: MQTT publish version

	return device, nil
}

func (d *Mercury200) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(time.Minute)
	taskStateUpdater.SetName("bind-mercury-200-state-updater-" + d.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (d *Mercury200) taskStateUpdater(ctx context.Context) (interface{}, error) {
	currentT1, currentT2, currentT3, currentT4, err := d.provider.PowerCounters()
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		// TODO: log
		return nil, err
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)
	serialNumber := d.SerialNumber()

	if prevT1 := atomic.LoadUint64(&d.tariff1); currentT1 != prevT1 {
		atomic.StoreUint64(&d.tariff1, currentT1)
		d.MQTTPublishAsync(ctx, Mercury200MQTTTopicTariff.Format(serialNumber, 1), 0, true, currentT1)
	}

	if prevT2 := atomic.LoadUint64(&d.tariff2); currentT2 != prevT2 {
		atomic.StoreUint64(&d.tariff1, currentT2)
		d.MQTTPublishAsync(ctx, Mercury200MQTTTopicTariff.Format(serialNumber, 2), 0, true, currentT2)
	}

	if prevT3 := atomic.LoadUint64(&d.tariff3); currentT3 != prevT3 {
		atomic.StoreUint64(&d.tariff3, currentT3)
		d.MQTTPublishAsync(ctx, Mercury200MQTTTopicTariff.Format(serialNumber, 3), 0, true, currentT3)
	}

	if prevT4 := atomic.LoadUint64(&d.tariff4); currentT4 != prevT4 {
		atomic.StoreUint64(&d.tariff4, currentT4)
		d.MQTTPublishAsync(ctx, Mercury200MQTTTopicTariff.Format(serialNumber, 4), 0, true, currentT4)
	}

	// optimization
	if currentVoltage, currentAmperage, currentPower, err := d.provider.ParamsCurrent(); err == nil {
		if prevVoltage := atomic.LoadUint64(&d.voltage); currentVoltage != prevVoltage {
			atomic.StoreUint64(&d.voltage, currentVoltage)
			d.MQTTPublishAsync(ctx, Mercury200MQTTTopicVoltage.Format(serialNumber), 0, true, currentVoltage)
		}

		if prevAmperage := math.Float64frombits(atomic.LoadUint64(&d.amperage)); currentAmperage != prevAmperage {
			atomic.StoreUint64(&d.amperage, math.Float64bits(currentAmperage))
			d.MQTTPublishAsync(ctx, Mercury200MQTTTopicAmperage.Format(serialNumber), 0, true, currentAmperage)
		}

		if prevPower := atomic.LoadUint64(&d.power); currentPower != prevPower {
			atomic.StoreUint64(&d.power, currentPower)
			d.MQTTPublishAsync(ctx, Mercury200MQTTTopicPower.Format(serialNumber), 0, true, currentPower)
		}
	} else {
		// TODO: log
	}

	if current, err := d.provider.BatteryVoltage(); err == nil {
		if prev := math.Float64frombits(atomic.LoadUint64(&d.batteryVoltage)); current != prev {
			atomic.StoreUint64(&d.batteryVoltage, math.Float64bits(current))
			d.MQTTPublishAsync(ctx, Mercury200MQTTTopicBatteryVoltage.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if val, err := d.provider.LastPowerOffDatetime(); err == nil {
		current := val.Unix()
		if prev := atomic.LoadInt64(&d.lastPowerOff); current != prev {
			atomic.StoreInt64(&d.lastPowerOff, current)
			d.MQTTPublishAsync(ctx, Mercury200MQTTTopicLastPowerOff.Format(serialNumber), 0, true, val)
		}
	} else {
		// TODO: log
	}

	if val, err := d.provider.LastPowerOnDatetime(); err == nil {
		current := val.Unix()
		if prev := atomic.LoadInt64(&d.lastPowerOn); current != prev {
			atomic.StoreInt64(&d.lastPowerOn, current)
			d.MQTTPublishAsync(ctx, Mercury200MQTTTopicLastPowerOn.Format(serialNumber), 0, true, val)
		}
	} else {
		// TODO: log
	}

	if val, err := d.provider.MakeDate(); err == nil {
		current := val.Unix()
		if prev := atomic.LoadInt64(&d.makeDate); current != prev {
			atomic.StoreInt64(&d.makeDate, current)
			d.MQTTPublishAsync(ctx, Mercury200MQTTTopicMakeDate.Format(serialNumber), 0, true, val)
		}
	} else {
		// TODO: log
	}

	if version, date, err := d.provider.Version(); err == nil {
		current := date.Unix()
		if prev := atomic.LoadInt64(&d.firmwareDate); current != prev {
			atomic.StoreInt64(&d.firmwareDate, current)
			d.MQTTPublishAsync(ctx, Mercury200MQTTTopicFirmwareDate.Format(serialNumber), 0, true, date)
		}

		d.mutex.Lock()
		if version != d.firmwareVersion {
			d.firmwareVersion = version
			d.MQTTPublishAsync(ctx, Mercury200MQTTTopicFirmwareVersion.Format(serialNumber), 0, true, version)
		}
		d.mutex.Unlock()
	} else {
		// TODO: log
	}

	return nil, nil
}

func (d *Mercury200) MQTTTopics() []mqtt.Topic {
	sn := d.SerialNumberMQTTEscaped()

	return []mqtt.Topic{
		mqtt.Topic(Mercury200MQTTTopicTariff.Format(sn, 1)),
		mqtt.Topic(Mercury200MQTTTopicTariff.Format(sn, 2)),
		mqtt.Topic(Mercury200MQTTTopicTariff.Format(sn, 3)),
		mqtt.Topic(Mercury200MQTTTopicTariff.Format(sn, 4)),
		mqtt.Topic(Mercury200MQTTTopicVoltage.Format(sn)),
		mqtt.Topic(Mercury200MQTTTopicAmperage.Format(sn)),
		mqtt.Topic(Mercury200MQTTTopicPower.Format(sn)),
		mqtt.Topic(Mercury200MQTTTopicBatteryVoltage.Format(sn)),
		mqtt.Topic(Mercury200MQTTTopicLastPowerOff.Format(sn)),
		mqtt.Topic(Mercury200MQTTTopicLastPowerOn.Format(sn)),
		mqtt.Topic(Mercury200MQTTTopicMakeDate.Format(sn)),
		mqtt.Topic(Mercury200MQTTTopicFirmwareDate.Format(sn)),
		mqtt.Topic(Mercury200MQTTTopicFirmwareVersion.Format(sn)),
	}
}
