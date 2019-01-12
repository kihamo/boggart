package bind

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ghthor/gowol"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/samsung/tv"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

const (
	SamsungTizenMQTTTopicPower           mqtt.Topic = boggart.ComponentName + "/tv/+/power"
	SamsungTizenMQTTTopicKey             mqtt.Topic = boggart.ComponentName + "/tv/+/key"
	SamsungTizenMQTTTopicDeviceID        mqtt.Topic = boggart.ComponentName + "/tv/+/device/id"
	SamsungTizenMQTTTopicDeviceModelName mqtt.Topic = boggart.ComponentName + "/tv/+/device/model-name"
)

type SamsungTizen struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	mutex    sync.RWMutex
	initOnce sync.Once

	client *tv.ApiV2
	mac    string
}

type SamsungTizenConfig struct {
	Host string `valid:"host,required"`
}

func (d SamsungTizen) Config() interface{} {
	return &SamsungTizenConfig{}
}

func (d SamsungTizen) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*SamsungTizenConfig)

	device := &SamsungTizen{
		client: tv.NewApiV2(config.Host),
	}
	device.Init()

	return device, nil
}

func (d *SamsungTizen) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Second * 30)
	taskLiveness.SetName("bind-samsung-tizen-liveness")

	return []workers.Task{
		taskLiveness,
	}
}

func (d *SamsungTizen) taskLiveness(ctx context.Context) (interface{}, error) {
	info, err := d.client.Device(ctx)
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, nil
	}

	if d.Status() == boggart.DeviceStatusOnline {
		return nil, nil
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)

	if d.SerialNumber() == "" {
		parts := strings.Split(info.ID, ":")
		if len(parts) > 1 {
			d.SetSerialNumber(parts[1])
		} else {
			d.SetSerialNumber(info.ID)
		}

		d.mutex.Lock()
		d.mac = info.Device.WifiMac
		d.mutex.Unlock()

		sn := d.SerialNumber()
		d.MQTTPublishAsync(ctx, SamsungTizenMQTTTopicDeviceID.Format(sn), 0, false, info.Device.ID)
		d.MQTTPublishAsync(ctx, SamsungTizenMQTTTopicDeviceModelName.Format(sn), 0, false, info.Device.Name)
	}

	return nil, nil
}

func (d *SamsungTizen) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(SamsungTizenMQTTTopicPower.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 2) {
				return
			}

			if bytes.Equal(message.Payload(), []byte(`1`)) {
				d.mutex.RLock()
				mac := d.mac
				d.mutex.RUnlock()

				wol.MagicWake(mac, "255.255.255.255")
			} else if d.Status() == boggart.DeviceStatusOnline {
				d.client.SendCommand(tv.KeyPower)
			}
		}),
		mqtt.NewSubscriber(SamsungTizenMQTTTopicKey.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 2) {
				return
			}

			d.client.SendCommand(string(message.Payload()))
		})),
	}
}
