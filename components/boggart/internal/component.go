package internal

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/boggart/components/boggart/providers/doors"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/annotations"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	annotations annotations.Component
	config      config.Component
	logger      logger.Logger
	messenger   messengers.Messenger
	workers     workers.Component
	routes      []dashboard.Route
	collector   *MetricsCollector

	devices boggart.DeviceManager

	connectionRS485 *rs485.Connection
	doorEntrance    *doors.Door
}

func (c *Component) Name() string {
	return boggart.ComponentName
}

func (c *Component) Version() string {
	return boggart.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name: annotations.ComponentName,
		},
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: dashboard.ComponentName,
		},
		{
			Name: logger.ComponentName,
		},
		{
			Name: messengers.ComponentName,
		},
		{
			Name:     metrics.ComponentName,
			Required: true,
		},
		{
			Name:     workers.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a

	if a.HasComponent(annotations.ComponentName) {
		c.annotations = a.GetComponent(annotations.ComponentName).(annotations.Component)
	}

	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.workers = a.GetComponent(workers.ComponentName).(workers.Component)
	c.collector = NewMetricsCollector(c)
	c.devices = NewDeviceManager(c.workers)

	return nil
}

func (c *Component) Run() (err error) {
	c.logger = logger.NewOrNop(c.Name(), c.application)

	if c.application.HasComponent(messengers.ComponentName) {
		c.messenger = c.application.GetComponent(messengers.ComponentName).(messengers.Component).Messenger(messengers.MessengerTelegram)
	}

	c.initVideoRecorders()
	c.initCameras()
	c.initPhones()
	c.initRouters()

	c.initConnectionRS485()

	taskMercury := task.NewFunctionTask(c.taskMercury)
	taskMercury.SetRepeats(-1)
	taskMercury.SetRepeatInterval(c.config.Duration(boggart.ConfigMercuryRepeatInterval))
	taskMercury.SetName(c.Name() + "-mercury-updater")
	c.workers.AddTask(taskMercury)

	taskPulsar := task.NewFunctionTask(c.taskPulsar)
	taskPulsar.SetRepeats(-1)
	taskPulsar.SetRepeatInterval(c.config.Duration(boggart.ConfigPulsarRepeatInterval))
	taskPulsar.SetName(c.Name() + "-pulsar-updater")
	c.workers.AddTask(taskPulsar)

	taskSoftVideo := task.NewFunctionTask(c.taskSoftVideo)
	taskSoftVideo.SetRepeats(-1)
	taskSoftVideo.SetRepeatInterval(c.config.Duration(boggart.ConfigSoftVideoRepeatInterval))
	taskSoftVideo.SetName(c.Name() + "-softvideo-updater")
	c.workers.AddTask(taskSoftVideo)

	c.doorEntrance, err = doors.NewDoor(c.config.Int(boggart.ConfigDoorsEntrancePin))
	if err != nil {
		return err
	}

	if c.application.HasComponent(annotations.ComponentName) {
		c.doorEntrance.SetCallbackChange(c.doorCallback)
	}

	return nil
}

func (c *Component) taskMercury(context.Context) (interface{}, error) {
	if c.config.Bool(boggart.ConfigMercuryEnabled) {
		err := c.collector.UpdaterMercury()
		if err != nil {
			c.logger.Error("Mercury updater failed", map[string]interface{}{
				"error": err.Error(),
			})
		}

		return nil, err
	}

	return nil, nil
}

func (c *Component) taskPulsar(context.Context) (interface{}, error) {
	if c.config.Bool(boggart.ConfigPulsarEnabled) {
		err := c.collector.UpdaterPulsar()
		if err != nil {
			c.logger.Error("Puslar updater failed", map[string]interface{}{
				"error": err.Error(),
			})
		}

		return nil, err
	}

	return nil, nil
}

func (c *Component) taskSoftVideo(context.Context) (interface{}, error) {
	if c.config.Bool(boggart.ConfigSoftVideoEnabled) {
		err := c.collector.UpdaterSoftVideo()
		if err != nil {
			c.logger.Error("SoftVideo updater failed", map[string]interface{}{
				"error": err.Error(),
			})
		}

		return nil, err
	}

	return nil, nil
}

func (c *Component) initConnectionRS485() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.connectionRS485 = rs485.NewConnection(
		c.config.String(boggart.ConfigRS485Address),
		c.config.Duration(boggart.ConfigRS485Timeout))
}

func (c *Component) ConnectionRS485() *rs485.Connection {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.connectionRS485
}

func (c *Component) DoorEntrance() boggart.Door {
	return c.doorEntrance
}

func (c *Component) doorCallback(status bool, changed *time.Time) {
	if c.messenger != nil {
		device := c.devices.Device(boggart.DeviceIdCameraHall.String())
		if device == nil && !device.IsEnabled() {
			return
		}

		image, err := device.(boggart.Camera).Snapshot(context.Background())
		if err == nil {
			if status {
				// TODO: changeUserId
				c.messenger.SendMessage("238815343", "Entrance door is opened")
			} else {
				// TODO: changeUserId
				c.messenger.SendMessage("238815343", "Entrance door is closed")
			}

			time.AfterFunc(time.Second, func() {
				// TODO: changeUserId
				c.messenger.SendPhoto("238815343", "Hall snapshot", bytes.NewReader(image))
			})
		}
	}

	if c.annotations == nil || !status {
		return
	}

	if changed == nil {
		changed = c.application.StartDate()
	}

	timeEnd := time.Now()
	diff := timeEnd.Sub(*changed)

	if c.annotations != nil {
		annotation := annotations.NewAnnotation(
			"Door is closed",
			fmt.Sprintf("Door was open for %.2f seconds", diff.Seconds()),
			[]string{"door", "entrance", "close"},
			changed,
			&timeEnd)

		if err := c.annotations.Create(annotation); err != nil {
			c.logger.Error("Create annotation failed", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}
