package bind

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/robbiet480/go.nut"
)

const (
	NUTMQTTTopicVariable mqtt.Topic = boggart.ComponentName + "/ups/+/+"
)

type NUT struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	config *NUTConfig

	mutex     sync.Mutex
	variables map[string]interface{}
}

type NUTConfig struct {
	Host string `valid:"host,required"`
	UPS  string `valid:"required"`
}

func (d NUT) Config() interface{} {
	return &NUTConfig{}
}

func (d NUT) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	device := &NUT{
		config:    c.(*NUTConfig),
		variables: make(map[string]interface{}, 0),
	}
	device.Init()

	return device, nil
}

func (d *NUT) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Second * 30)
	taskLiveness.SetName("bind-nut-liveness-" + d.config.UPS)

	taskStateUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(time.Minute)
	taskStateUpdater.SetName("bind-nut-state-updater-" + d.config.UPS)

	return []workers.Task{
		taskLiveness,
		taskStateUpdater,
	}
}

func (d *NUT) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(NUTMQTTTopicVariable.Format(d.config.UPS)),
	}
}

func (d *NUT) getUPS() (ups nut.UPS, err error) {
	client, err := nut.Connect(d.config.Host)
	if err != nil {
		return ups, err
	}

	devices, err := client.GetUPSList()
	if err != nil {
		return ups, err
	}

	for _, device := range devices {
		if device.Name == d.config.UPS {
			for _, v := range device.Variables {
				if v.Name == "device.serial" {
					d.SetSerialNumber(v.Value.(string))
					return device, nil
				}
			}

			break
		}
	}

	return ups, errors.New("device not found")
}

func (d *NUT) taskLiveness(ctx context.Context) (interface{}, error) {
	_, err := d.getUPS()
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)
	return nil, nil
}

func (d *NUT) taskStateUpdater(ctx context.Context) (interface{}, error) {
	if d.Status() != boggart.DeviceStatusOnline {
		return nil, nil
	}

	ups, err := d.getUPS()
	if err != nil {
		return nil, err
	}

	d.mutex.Lock()

	for _, v := range ups.Variables {
		prev, ok := d.variables[v.Name]
		if !ok || prev != v.Value {
			d.variables[v.Name] = v.Value
			name := mqtt.NameReplace(v.Name)

			d.MQTTPublishAsync(ctx, NUTMQTTTopicVariable.Format(d.config.UPS, name), 2, true, v.Value)
		}
	}

	d.mutex.Unlock()
	return nil, nil
}
