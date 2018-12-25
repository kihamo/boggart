package devices

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
	TVSamsungMQTTTopicPower           mqtt.Topic = boggart.ComponentName + "/tv/+/power"
	TVSamsungMQTTTopicKey             mqtt.Topic = boggart.ComponentName + "/tv/+/key"
	TVSamsungMQTTTopicDeviceID        mqtt.Topic = boggart.ComponentName + "/tv/+/device/id"
	TVSamsungMQTTTopicDeviceModelName mqtt.Topic = boggart.ComponentName + "/tv/+/device/model-name"
)

type SamsungTV struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	mutex    sync.RWMutex
	initOnce sync.Once

	client *tv.ApiV2
	mac    string
}

func NewSamsungTV(host string) *SamsungTV {
	device := &SamsungTV{
		client: tv.NewApiV2(host),
	}
	device.Init()
	device.SetDescription("Samsung TV")

	return device
}

func (d *SamsungTV) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeTV,
	}
}

func (d *SamsungTV) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Second * 30)
	taskLiveness.SetName("device-tv-samsung-liveness")

	return []workers.Task{
		taskLiveness,
	}
}

func (d *SamsungTV) taskLiveness(ctx context.Context) (interface{}, error) {
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
		d.SetSerialNumber(info.ID)

		d.mutex.Lock()
		d.mac = info.Device.WifiMac
		d.mutex.Unlock()

		d.MQTTPublishAsync(ctx, TVSamsungMQTTTopicDeviceID.Format(info.ID), 0, false, info.Device.ID)
		d.MQTTPublishAsync(ctx, TVSamsungMQTTTopicDeviceModelName.Format(info.ID), 0, false, info.Device.Name)

		d.initOnce.Do(d.initMQTTSubscribers)
	}

	return nil, nil
}

func (d *SamsungTV) initMQTTSubscribers() {
	parts := strings.Split(d.SerialNumber(), ":")
	if len(parts) < 2 {
		return
	}

	sn := d.SerialNumber()

	d.MQTTSubscribe(TVSamsungMQTTTopicPower.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		if bytes.Equal(message.Payload(), []byte(`1`)) {
			d.mutex.RLock()
			mac := d.mac
			d.mutex.RUnlock()

			wol.MagicWake(mac, "255.255.255.255")
		} else if d.Status() == boggart.DeviceStatusOnline {
			d.client.SendCommand(tv.KeyPower)
		}
	})

	d.MQTTSubscribe(TVSamsungMQTTTopicKey.Format(sn[1]), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		d.client.SendCommand(string(message.Payload()))
	}))
}
