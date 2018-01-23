package internal

import (
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
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
)

type Component struct {
	mutex sync.RWMutex

	application shadow.Application
	annotations annotations.Component
	config      config.Component
	logger      logger.Logger
	workers     workers.Component
	routes      []dashboard.Route
	collector   *MetricsCollector

	connectionRS485 *rs485.Connection
	doorEntrance    *doors.Door
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
	c.annotations = a.GetComponent(annotations.ComponentName).(annotations.Component)
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.workers = a.GetComponent(workers.ComponentName).(workers.Component)
	c.collector = NewMetricsCollector(c)

	return nil
}

func (c *Component) Run() (err error) {
	c.logger = logger.NewOrNop(c.GetName(), c.application)

	c.initConnectionRS485()

	taskMercury := task.NewFunctionTask(c.taskMercury)
	taskMercury.SetRepeats(-1)
	taskMercury.SetRepeatInterval(c.config.GetDuration(boggart.ConfigMercuryRepeatInterval))
	taskMercury.SetName(c.GetName() + "-mercury-updater")
	c.workers.AddTask(taskMercury)

	taskPulsar := task.NewFunctionTask(c.taskPulsar)
	taskPulsar.SetRepeats(-1)
	taskPulsar.SetRepeatInterval(c.config.GetDuration(boggart.ConfigPulsarRepeatInterval))
	taskPulsar.SetName(c.GetName() + "-pulsar-updater")
	c.workers.AddTask(taskPulsar)

	taskSoftVideo := task.NewFunctionTask(c.taskSoftVideo)
	taskSoftVideo.SetRepeats(-1)
	taskSoftVideo.SetRepeatInterval(c.config.GetDuration(boggart.ConfigSoftVideoRepeatInterval))
	taskSoftVideo.SetName(c.GetName() + "-softvideo-updater")
	c.workers.AddTask(taskSoftVideo)

	c.doorEntrance, err = doors.NewDoor(c.config.GetInt(boggart.ConfigDoorsEntrancePin))
	if err != nil {
		return err
	}

	if c.application.HasComponent(annotations.ComponentName) {
		c.doorEntrance.SetCallbackChange(c.doorCallback)
	}

	return nil
}

func (c *Component) taskMercury(context.Context) (interface{}, error) {
	if c.config.GetBool(boggart.ConfigMercuryEnabled) {
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

func (c *Component) initConnectionRS485() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.connectionRS485 = rs485.NewConnection(
		c.config.GetString(boggart.ConfigRS485Address),
		c.config.GetDuration(boggart.ConfigRS485Timeout))
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
	// door is open
	if !status {
		c.logger.Info("Entrance door opened")
		return
	}

	if changed == nil {
		changed = c.application.GetStartDate()
	}

	timeEnd := time.Now()
	diff := timeEnd.Sub(*changed)

	annotation := annotations.NewAnnotation(
		"Door is open",
		fmt.Sprintf("%.2f seconds", diff.Seconds()),
		[]string{"door", "entrance"},
		changed,
		&timeEnd)

	if err := c.annotations.Create(annotation); err != nil {
		c.logger.Error("Create annotation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
