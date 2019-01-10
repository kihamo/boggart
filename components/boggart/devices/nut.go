package devices

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
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	host string
	ups  string

	mutex     sync.Mutex
	variables map[string]interface{}
}

func (d NUT) CreateBind(config map[string]interface{}) (boggart.DeviceBind, error) {
	host, ok := config["host"]
	if !ok {
		return nil, errors.New("config option host isn't set")
	}

	if host == "" {
		return nil, errors.New("config option host is empty")
	}

	ups, ok := config["ups"]
	if !ok {
		return nil, errors.New("config option ups isn't set")
	}

	if ups == "" {
		return nil, errors.New("config option ups is empty")
	}

	device := &NUT{
		host:      host.(string),
		ups:       ups.(string),
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
	taskLiveness.SetName("ups-nut-liveness-" + d.ups)

	taskStateUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(time.Minute)
	taskStateUpdater.SetName("ups-nut-state-updater-" + d.ups)

	return []workers.Task{
		taskLiveness,
		taskStateUpdater,
	}
}

func (d *NUT) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(NUTMQTTTopicVariable.Format(d.ups)),
	}
}

func (d *NUT) getUPS() (ups nut.UPS, err error) {
	client, err := nut.Connect(d.host)
	if err != nil {
		return ups, err
	}

	devices, err := client.GetUPSList()
	if err != nil {
		return ups, err
	}

	for _, device := range devices {
		if device.Name == d.ups {
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

			d.MQTTPublishAsync(ctx, NUTMQTTTopicVariable.Format(d.ups, name), 2, true, v.Value)
		}
	}

	d.mutex.Unlock()
	return nil, nil
}
