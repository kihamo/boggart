package di

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/snitch"
)

const (
	ProbesConfigReadinessDefaultPeriod                  = time.Minute
	ProbesConfigReadinessDefaultTimeout                 = time.Second * 30
	ProbesConfigReadinessDefaultThresholdSuccess uint64 = 1
	ProbesConfigReadinessDefaultThresholdFailure uint64 = 1
	ProbesConfigReadinessDefaultName                    = "readiness-probe"
	ProbesConfigLivenessDefaultPeriod                   = time.Minute
	ProbesConfigLivenessDefaultTimeout                  = time.Second * 30
	ProbesConfigLivenessDefaultThresholdSuccess  uint64 = 1
	ProbesConfigLivenessDefaultThresholdFailure  uint64 = 1
	ProbesConfigLivenessDefaultName                     = "liveness-probe"
)

type BindHasReadinessProbeHandler interface {
	ReadinessProbe(context.Context) error
}

type BindHasReadinessProbeSchedule interface {
	ReadinessProbeSchedule(tasks.Meta) time.Time
}

type BindHasLivenessProbeHandler interface {
	LivenessProbe(context.Context) error
}

type BindHasLivenessProbeSchedule interface {
	LivenessProbeSchedule(tasks.Meta) time.Time
}

type ProbesContainerSupport interface {
	SetProbes(*ProbesContainer)
	Probes() *ProbesContainer
}

func ProbesContainerBind(bind boggart.Bind) (*ProbesContainer, bool) {
	if support, ok := bind.(ProbesContainerSupport); ok {
		container := support.Probes()
		return container, container != nil
	}

	return nil, false
}

type ProbesBind struct {
	mutex     sync.RWMutex
	container *ProbesContainer
}

func (b *ProbesBind) SetProbes(container *ProbesContainer) {
	b.mutex.Lock()
	b.container = container
	b.mutex.Unlock()
}

func (b *ProbesBind) Probes() *ProbesContainer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.container
}

type ProbesContainer struct {
	bindItem boggart.BindItem

	statusManager func(boggart.BindStatus)
	register      func() error
	unregister    func() error
	tasksManager  *tasks.Manager

	probeIDMutex     sync.RWMutex
	probeReadinessID string
	probeLivenessID  string

	metricProbes snitch.Counter
}

func NewProbesContainer(bindItem boggart.BindItem, statusManager func(boggart.BindStatus), register, unregister func() error, manager *tasks.Manager, metricProbes snitch.Counter) *ProbesContainer {
	return &ProbesContainer{
		bindItem:      bindItem,
		statusManager: statusManager,
		register:      register,
		unregister:    unregister,
		tasksManager:  manager,
		metricProbes:  metricProbes,
	}
}

func (c *ProbesContainer) HookRegister() (err error) {
	taskReadiness := c.taskReadiness()
	taskLiveness := c.taskLiveness()

	if taskReadiness == nil && taskLiveness == nil {
		return nil
	}

	var id string

	bindWorkersSupport, ok := WorkersContainerBind(c.bindItem.Bind())

	if taskReadiness != nil {
		// если воркеры есть у привязки, то регистрируем пробы как таски, чтобы отображались в общем списке тасок
		if ok {
			id, err = bindWorkersSupport.RegisterTask(taskReadiness)
		} else {
			taskReadiness.WithName("bind/" + c.bindItem.Type() + "/" + c.bindItem.ID() + "/" + taskReadiness.Name())
			id, err = c.tasksManager.Register(taskReadiness)
		}

		if err != nil {
			return fmt.Errorf("register readiness probe failed: %w", err)
		}

		c.probeIDMutex.Lock()
		c.probeReadinessID = id
		c.probeIDMutex.Unlock()
	} else {
		c.statusManager(boggart.BindStatusOnline)
	}

	if taskLiveness != nil {
		if ok {
			id, err = bindWorkersSupport.RegisterTask(taskLiveness)
		} else {
			taskLiveness.WithName("bind/" + c.bindItem.Type() + "/" + c.bindItem.ID() + "/" + taskLiveness.Name())
			id, err = c.tasksManager.Register(taskLiveness)
		}

		if err != nil {
			return fmt.Errorf("register liveness probe failed: %w", err)
		}

		c.probeIDMutex.Lock()
		c.probeLivenessID = id
		c.probeIDMutex.Unlock()
	}

	return nil
}

func (c *ProbesContainer) HookUnregister() {
	// если есть DI воркеров, то там удалиться все через стандартный механизм,
	// дополнительно ничего не нужно делать
	if _, ok := WorkersContainerBind(c.bindItem.Bind()); ok {
		return
	}

	if id := c.ReadinessTaskID(); id != "" {
		c.tasksManager.Unregister(id)

		c.probeIDMutex.Lock()
		c.probeReadinessID = ""
		c.probeIDMutex.Unlock()
	}

	if id := c.LivenessTaskID(); id != "" {
		c.tasksManager.Unregister(id)

		c.probeIDMutex.Lock()
		c.probeLivenessID = ""
		c.probeIDMutex.Unlock()
	}
}

func (c *ProbesContainer) taskReadiness() *tasks.TaskBase {
	bind := c.bindItem.Bind()

	hasHandler, ok := bind.(BindHasReadinessProbeHandler)
	if !ok {
		return nil
	}

	// options
	logger, _ := LoggerContainerBind(bind)

	probePeriod := ProbesConfigReadinessDefaultPeriod
	probeTimeout := ProbesConfigReadinessDefaultTimeout
	probeThresholdSuccess := ProbesConfigReadinessDefaultThresholdSuccess
	probeThresholdFailure := ProbesConfigReadinessDefaultThresholdFailure

	var probeSuccess, probeFailure uint64

	if config, ok := ConfigForBind(bind); ok {
		if probeConfig, ok := config.(ProbesConfigReadiness); ok {
			probePeriod = probeConfig.ReadinessProbePeriod()
			probeTimeout = probeConfig.ReadinessProbeTimeout()
			probeThresholdSuccess = probeConfig.ReadinessProbeThresholdSuccess()
			probeThresholdFailure = probeConfig.ReadinessProbeThresholdFailure()
		}
	}

	// task
	var schedule tasks.Schedule

	if hasSchedule, ok := bind.(BindHasReadinessProbeSchedule); ok {
		schedule = tasks.ScheduleFunc(hasSchedule.ReadinessProbeSchedule)
	} else {
		schedule = tasks.ScheduleWithDuration(tasks.ScheduleNow(), probePeriod)
	}

	handler := func(ctx context.Context) (err error) {
		if s := c.bindItem.Status(); !s.IsStatusInitializing() && !s.IsStatusOnline() && !s.IsStatusOffline() {
			return nil
		}

		ch := make(chan error, 1)
		go func() {
			ch <- hasHandler.ReadinessProbe(ctx)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case err = <-ch:
		}

		if err == nil {
			threshold := atomic.AddUint64(&probeSuccess, 1)
			atomic.StoreUint64(&probeFailure, 0)

			c.metricProbes.With(
				"probe", "readiness",
				"status", "success",
				"id", c.bindItem.ID(),
				"type", c.bindItem.Type(),
			).Inc()

			if threshold != probeThresholdSuccess {
				return nil
			}

			atomic.StoreUint64(&probeSuccess, 0)
		} else {
			threshold := atomic.AddUint64(&probeFailure, 1)
			atomic.StoreUint64(&probeSuccess, 0)

			c.metricProbes.With(
				"probe", "readiness",
				"status", "failed",
				"id", c.bindItem.ID(),
				"type", c.bindItem.Type(),
			).Inc()

			logger.Error("Readiness probe failure",
				"type", c.bindItem.Type(),
				"id", c.bindItem.ID(),
				"error", err.Error(),
				"threshold", threshold,
			)

			if threshold != probeThresholdFailure {
				return err
			}

			atomic.StoreUint64(&probeFailure, 0)
		}

		if err != nil {
			c.statusManager(boggart.BindStatusOffline)
		} else {
			c.statusManager(boggart.BindStatusOnline)
		}

		return err
	}

	return tasks.NewTask().
		WithName(ProbesConfigReadinessDefaultName).
		WithSchedule(schedule).
		WithHandler(tasks.HandlerWithTimeout(tasks.HandlerFuncFromShortToLong(handler), probeTimeout))
}

func (c *ProbesContainer) ReadinessTaskID() string {
	c.probeIDMutex.RLock()
	defer c.probeIDMutex.RUnlock()

	return c.probeReadinessID
}

func (c *ProbesContainer) taskLiveness() *tasks.TaskBase {
	bind := c.bindItem.Bind()

	hasHandler, ok := bind.(BindHasLivenessProbeHandler)
	if !ok {
		return nil
	}

	// options
	logger, _ := LoggerContainerBind(bind)

	probePeriod := ProbesConfigLivenessDefaultPeriod
	probeTimeout := ProbesConfigLivenessDefaultTimeout
	probeThresholdSuccess := ProbesConfigLivenessDefaultThresholdSuccess
	probeThresholdFailure := ProbesConfigLivenessDefaultThresholdFailure

	var probeSuccess, probeFailure uint64

	if config, ok := ConfigForBind(bind); ok {
		if probeConfig, ok := config.(ProbesConfigLiveness); ok {
			probePeriod = probeConfig.LivenessProbePeriod()
			probeTimeout = probeConfig.LivenessProbeTimeout()
			probeThresholdSuccess = probeConfig.LivenessProbeThresholdSuccess()
			probeThresholdFailure = probeConfig.LivenessProbeThresholdFailure()
		}
	}

	// task
	var schedule tasks.Schedule

	if hasSchedule, ok := bind.(BindHasLivenessProbeSchedule); ok {
		schedule = tasks.ScheduleFunc(hasSchedule.LivenessProbeSchedule)
	} else {
		schedule = tasks.ScheduleWithDuration(tasks.ScheduleNow(), probePeriod)
	}

	handler := func(ctx context.Context) (err error) {
		if s := c.bindItem.Status(); !s.IsStatusOnline() && !s.IsStatusOffline() {
			return nil
		}

		ch := make(chan error, 1)
		go func() {
			ch <- hasHandler.LivenessProbe(ctx)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case err = <-ch:
		}

		if err == nil {
			threshold := atomic.AddUint64(&probeSuccess, 1)
			atomic.StoreUint64(&probeFailure, 0)

			c.metricProbes.With(
				"probe", "liveness",
				"status", "success",
				"id", c.bindItem.ID(),
				"type", c.bindItem.Type(),
			).Inc()

			if threshold == probeThresholdSuccess {
				atomic.StoreUint64(&probeSuccess, 0)
			}

			// success return
			return nil
		}

		threshold := atomic.AddUint64(&probeFailure, 1)
		atomic.StoreUint64(&probeSuccess, 0)

		c.metricProbes.With(
			"probe", "liveness",
			"status", "failed",
			"id", c.bindItem.ID(),
			"type", c.bindItem.Type(),
		).Inc()

		logger.Error("Liveness probe failure",
			"type", c.bindItem.Type(),
			"id", c.bindItem.ID(),
			"error", err.Error(),
			"threshold", threshold,
		)

		if threshold != probeThresholdFailure {
			return err
		}

		atomic.StoreUint64(&probeFailure, 0)

		if err := c.unregister(); err != nil {
			logger.Error("Unregister after liveness probe failure",
				"type", c.bindItem.Type(),
				"id", c.bindItem.ID(),
				"error", err.Error(),
			)
			return err
		}

		if err := c.register(); err != nil {
			logger.Error("Register after liveness probe failure",
				"type", c.bindItem.Type(),
				"id", c.bindItem.ID(),
				"error", err.Error(),
			)

			return err
		}

		return nil
	}

	return tasks.NewTask().
		WithName(ProbesConfigLivenessDefaultName).
		WithSchedule(schedule).
		WithHandler(tasks.HandlerWithTimeout(tasks.HandlerFuncFromShortToLong(handler), probeTimeout))
}

func (c *ProbesContainer) LivenessTaskID() string {
	c.probeIDMutex.RLock()
	defer c.probeIDMutex.RUnlock()

	return c.probeLivenessID
}

type ProbesConfigReadiness interface {
	ReadinessProbePeriod() time.Duration
	ReadinessProbeTimeout() time.Duration
	ReadinessProbeThresholdSuccess() uint64
	ReadinessProbeThresholdFailure() uint64
}

type ProbesConfigLiveness interface {
	LivenessProbePeriod() time.Duration
	LivenessProbeTimeout() time.Duration
	LivenessProbeThresholdSuccess() uint64
	LivenessProbeThresholdFailure() uint64
}

type ProbesConfig struct {
	ReadinessPeriod           time.Duration `mapstructure:"readiness_probe_period" yaml:"readiness_probe_period"`
	ReadinessTimeout          time.Duration `mapstructure:"readiness_probe_timeout" yaml:"readiness_probe_timeout"`
	ReadinessThresholdSuccess uint64        `mapstructure:"readiness_threshold_success" yaml:"readiness_threshold_success"`
	ReadinessThresholdFailure uint64        `mapstructure:"readiness_threshold_failure" yaml:"readiness_threshold_failure"`
	LivenessPeriod            time.Duration `mapstructure:"liveness_probe_period" yaml:"liveness_probe_period"`
	LivenessTimeout           time.Duration `mapstructure:"liveness_probe_timeout" yaml:"liveness_probe_timeout"`
	LivenessThresholdSuccess  uint64        `mapstructure:"liveness_threshold_success" yaml:"liveness_threshold_success"`
	LivenessThresholdFailure  uint64        `mapstructure:"liveness_threshold_failure" yaml:"liveness_threshold_failure"`
}

func ProbesConfigDefaults() (c ProbesConfig) {
	c.ReadinessPeriod = ProbesConfigReadinessDefaultPeriod
	c.ReadinessThresholdSuccess = ProbesConfigReadinessDefaultThresholdSuccess
	c.ReadinessThresholdFailure = ProbesConfigReadinessDefaultThresholdFailure
	c.LivenessPeriod = ProbesConfigLivenessDefaultPeriod
	c.LivenessThresholdSuccess = ProbesConfigLivenessDefaultThresholdSuccess
	c.LivenessThresholdFailure = ProbesConfigLivenessDefaultThresholdFailure

	return c
}

func (c ProbesConfig) ReadinessProbePeriod() time.Duration {
	if c.ReadinessPeriod <= 0 {
		return ProbesConfigReadinessDefaultPeriod
	}

	return c.ReadinessPeriod
}

func (c ProbesConfig) ReadinessProbeTimeout() time.Duration {
	return c.ReadinessTimeout
}

func (c ProbesConfig) ReadinessProbeThresholdSuccess() uint64 {
	if c.ReadinessThresholdSuccess == 0 {
		return ProbesConfigReadinessDefaultThresholdSuccess
	}

	return c.ReadinessThresholdSuccess
}

func (c ProbesConfig) ReadinessProbeThresholdFailure() uint64 {
	if c.ReadinessThresholdFailure == 0 {
		return ProbesConfigReadinessDefaultThresholdFailure
	}

	return c.ReadinessThresholdFailure
}

func (c ProbesConfig) LivenessProbePeriod() time.Duration {
	if c.LivenessPeriod <= 0 {
		return ProbesConfigLivenessDefaultPeriod
	}

	return c.LivenessPeriod
}

func (c ProbesConfig) LivenessProbeTimeout() time.Duration {
	return c.LivenessTimeout
}

func (c ProbesConfig) LivenessProbeThresholdSuccess() uint64 {
	if c.LivenessThresholdSuccess == 0 {
		return ProbesConfigLivenessDefaultThresholdSuccess
	}

	return c.LivenessThresholdSuccess
}

func (c ProbesConfig) LivenessProbeThresholdFailure() uint64 {
	if c.LivenessThresholdFailure == 0 {
		return ProbesConfigLivenessDefaultThresholdFailure
	}

	return c.LivenessThresholdFailure
}
