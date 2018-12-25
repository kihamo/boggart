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
	UPSNUTMQTTTopicVariable mqtt.Topic = boggart.ComponentName + "/ups/+/+"
)

type UPSNUT struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	host string
	ups  string

	mutex     sync.Mutex
	variables map[string]interface{}
}

func NewUPSNUT(host, ups string) *UPSNUT {
	device := &UPSNUT{
		host:      host,
		ups:       ups,
		variables: make(map[string]interface{}, 0),
	}

	device.Init()
	device.SetDescription("UPS NUT")

	return device
}

func (d *UPSNUT) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeUPS,
	}
}

func (d *UPSNUT) Tasks() []workers.Task {
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

func (d *UPSNUT) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(UPSNUTMQTTTopicVariable.Format(d.ups)),
	}
}

func (d *UPSNUT) getUPS() (ups nut.UPS, err error) {
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

func (d *UPSNUT) taskLiveness(ctx context.Context) (interface{}, error) {
	_, err := d.getUPS()
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)
	return nil, nil
}

func (d *UPSNUT) taskStateUpdater(ctx context.Context) (interface{}, error) {
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

			d.MQTTPublishAsync(ctx, UPSNUTMQTTTopicVariable.Format(d.ups, name), 1, false, v.Value)
		}
	}

	d.mutex.Unlock()
	return nil, nil
}
