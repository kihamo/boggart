package internal

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/doors"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	workers     workers.Component
	routes      []dashboard.Route
	collector   *MetricsCollector

	doorEntrance *doors.Door
}

func (c *Component) GetName() string {
	return boggart.ComponentName
}

func (c *Component) GetVersion() string {
	return boggart.ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{
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
			Name: metrics.ComponentName,
		},
		{
			Name:     workers.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.workers = a.GetComponent(workers.ComponentName).(workers.Component)
	c.collector = NewMetricsCollector(c)

	return nil
}

func (c *Component) Run() (err error) {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	taskSoftVideo := task.NewFunctionTask(c.taskSoftVideo)
	taskSoftVideo.SetRepeats(-1)
	taskSoftVideo.SetRepeatInterval(c.config.GetDuration(boggart.ConfigSoftVideoRepeatInterval))
	taskSoftVideo.SetName(c.GetName() + "-softvideo-updater")
	c.workers.AddTask(taskSoftVideo)

	taskPulsar := task.NewFunctionTask(c.taskPulsar)
	taskPulsar.SetRepeats(-1)
	taskPulsar.SetRepeatInterval(c.config.GetDuration(boggart.ConfigPulsarRepeatInterval))
	taskPulsar.SetName(c.GetName() + "-pulsar-updater")
	c.workers.AddTask(taskPulsar)

	c.doorEntrance, err = doors.NewDoor(c.config.GetInt(boggart.ConfigDoorsEntrancePin))
	if err != nil {
		return err
	}

	return nil
}

func (c *Component) taskPulsar(context.Context) (interface{}, error) {
	if c.config.GetBool(boggart.ConfigPulsarEnabled) {
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
	if c.config.GetBool(boggart.ConfigSoftVideoEnabled) {
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

func (c *Component) DoorEntrance() boggart.Door {
	return c.doorEntrance
}
