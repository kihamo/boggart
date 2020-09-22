package di

import (
	"context"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	w "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow/components/workers"
	"github.com/kihamo/snitch"
	"go.uber.org/multierr"
)

const (
	ProbesConfigReadinessDefaultPeriod  = time.Minute
	ProbesConfigReadinessDefaultTimeout = time.Second * 30
	ProbesConfigLivenessDefaultPeriod   = time.Minute
	ProbesConfigLivenessDefaultTimeout  = time.Second * 30
)

type ProbesHasReadinessProbe interface {
	ReadinessProbe(context.Context) error
}

type ProbesHasLivenessProbe interface {
	LivenessProbe(context.Context) error
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
	bind boggart.BindItem

	statusManager func(boggart.BindStatus)
	register      func() error
	unregister    func() error
	client        workers.Component

	probeReadiness w.Task
	probeLiveness  w.Task

	metricProbes snitch.Counter
}

func NewProbesContainer(bind boggart.BindItem, statusManager func(boggart.BindStatus), register, unregister func() error, client workers.Component, metricProbes snitch.Counter) *ProbesContainer {
	return &ProbesContainer{
		bind:          bind,
		statusManager: statusManager,
		register:      register,
		unregister:    unregister,
		client:        client,
		metricProbes:  metricProbes,
	}
}

func (c *ProbesContainer) HookRegister() {
	probeReadiness := c.Readiness()
	probeLiveness := c.Liveness()

	if probeReadiness == nil && probeLiveness == nil {
		return
	}

	bindWorkersSupport, ok := WorkersContainerBind(c.bind.Bind())

	if probeReadiness != nil {
		// если воркеры есть у привязки, то регистрируем пробы как таски, чтобы отображались в общем списке тасок
		if ok {
			bindWorkersSupport.RegisterTask(probeReadiness)
		} else {
			c.client.AddTask(probeReadiness)
		}
	} else {
		c.statusManager(boggart.BindStatusOnline)
	}

	if probeLiveness != nil {
		if ok {
			bindWorkersSupport.RegisterTask(probeLiveness)
		} else {
			c.client.AddTask(probeLiveness)
		}
	}
}

func (c *ProbesContainer) HookUnregister() {
	// если есть DI воркеров, то там удалиться все через стандартный механизм,
	// дополнительно ничего не нужно делать
	if _, ok := WorkersContainerBind(c.bind.Bind()); ok {
		return
	}

	if probe := c.Readiness(); probe != nil {
		c.client.RemoveTask(probe)
	}

	if probe := c.Liveness(); probe != nil {
		c.client.RemoveTask(probe)
	}
}

func (c *ProbesContainer) Readiness() w.Task {
	if c.probeReadiness != nil {
		return c.probeReadiness
	}

	has, ok := c.bind.Bind().(ProbesHasReadinessProbe)
	if !ok {
		return nil
	}

	logger, _ := LoggerContainerBind(c.bind.Bind())

	probeTask := task.NewFunctionTask(func(ctx context.Context) (_ interface{}, err error) {
		switch c.bind.Status() {
		case boggart.BindStatusInitializing, boggart.BindStatusOnline, boggart.BindStatusOffline:
		default:
			return nil, nil
		}

		ch := make(chan error, 1)
		go func() {
			ch <- has.ReadinessProbe(ctx)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case err = <-ch:
		}

		if err != nil {
			c.statusManager(boggart.BindStatusOffline)
			c.metricProbes.With("probe", "readiness", "status", "failed", "id", c.bind.ID(), "type", c.bind.Type()).Inc()

			if logger != nil {
				logger.Error("Readiness probe failure",
					"type", c.bind.Type(),
					"id", c.bind.ID(),
					"error", err.Error(),
				)
			}
		} else {
			c.statusManager(boggart.BindStatusOnline)
			c.metricProbes.With("probe", "readiness", "status", "success", "id", c.bind.ID(), "type", c.bind.Type()).Inc()
		}

		return nil, err
	})

	probePeriod := ProbesConfigReadinessDefaultPeriod
	probeTimeout := ProbesConfigReadinessDefaultTimeout

	if probeConfig, ok := c.bind.Config().(ProbesConfigReadiness); ok {
		probePeriod = probeConfig.ReadinessProbePeriod()
		probeTimeout = probeConfig.ReadinessProbeTimeout()
	}

	probeTask.SetRepeatInterval(probePeriod)
	probeTask.SetTimeout(probeTimeout)
	probeTask.SetRepeats(-1)

	if _, ok := c.bind.Bind().(WorkersContainerSupport); ok {
		probeTask.SetName("readiness-probe")
	} else {
		probeTask.SetName("bind-" + c.bind.ID() + "-" + c.bind.Type() + "-readiness-probe")
	}

	c.probeReadiness = probeTask

	return probeTask
}

func (c *ProbesContainer) ReadinessCheck(ctx context.Context) (err error) {
	if c.probeReadiness != nil {
		if timeout := c.probeReadiness.Timeout(); timeout > 0 {
			var cancel context.CancelFunc

			ctx, cancel = context.WithTimeout(ctx, timeout)

			defer cancel()
		}

		_, err = c.probeReadiness.Run(ctx)
	}

	return err
}

func (c *ProbesContainer) Liveness() w.Task {
	if c.probeLiveness != nil {
		return c.probeLiveness
	}

	has, ok := c.bind.Bind().(ProbesHasLivenessProbe)
	if !ok {
		return nil
	}

	logger, _ := LoggerContainerBind(c.bind.Bind())

	probeTask := task.NewFunctionTask(func(ctx context.Context) (_ interface{}, err error) {
		switch c.bind.Status() {
		case boggart.BindStatusOnline, boggart.BindStatusOffline:
		default:
			return nil, nil
		}

		ch := make(chan error, 1)
		go func() {
			ch <- has.LivenessProbe(ctx)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case err = <-ch:
		}

		if err == nil {
			c.metricProbes.With("probe", "liveness", "status", "success", "id", c.bind.ID(), "type", c.bind.Type()).Inc()
			return nil, nil
		}

		c.metricProbes.With("probe", "liveness", "status", "failed", "id", c.bind.ID(), "type", c.bind.Type()).Inc()

		if logger != nil {
			logger.Error("Liveness probe failure",
				"type", c.bind.Type(),
				"id", c.bind.ID(),
				"error", err.Error(),
			)
		}

		if e := c.unregister(); e != nil {
			if logger != nil {
				logger.Error("Unregister after liveness probe failure",
					"type", c.bind.Type(),
					"id", c.bind.ID(),
					"error", e.Error(),
				)
			}

			err = multierr.Append(err, e)
			return nil, err
		}

		if e := c.register(); e != nil {
			if logger != nil {
				logger.Error("Register after liveness probe failure",
					"type", c.bind.Type(),
					"id", c.bind.ID(),
					"error", e.Error(),
				)
			}

			err = multierr.Append(err, e)
			return nil, err
		}

		return nil, err
	})

	probePeriod := ProbesConfigLivenessDefaultPeriod
	probeTimeout := ProbesConfigLivenessDefaultTimeout

	if probeConfig, ok := c.bind.Config().(ProbesConfigLiveness); ok {
		probePeriod = probeConfig.LivenessProbePeriod()
		probeTimeout = probeConfig.LivenessProbeTimeout()
	}

	probeTask.SetRepeatInterval(probePeriod)
	probeTask.SetTimeout(probeTimeout)
	probeTask.SetRepeats(-1)

	if _, ok := c.bind.Bind().(WorkersContainerSupport); ok {
		probeTask.SetName("liveness-probe")
	} else {
		probeTask.SetName("bind-" + c.bind.ID() + "-" + c.bind.Type() + "-liveness-probe")
	}

	c.probeLiveness = probeTask

	return probeTask
}

func (c *ProbesContainer) LivenessCheck(ctx context.Context) (err error) {
	if c.probeLiveness != nil {
		if timeout := c.probeLiveness.Timeout(); timeout > 0 {
			var cancel context.CancelFunc

			ctx, cancel = context.WithTimeout(ctx, timeout)

			defer cancel()
		}

		_, err = c.probeLiveness.Run(ctx)
	}

	return err
}

type ProbesConfigReadiness interface {
	ReadinessProbePeriod() time.Duration
	ReadinessProbeTimeout() time.Duration
}

type ProbesConfigLiveness interface {
	LivenessProbePeriod() time.Duration
	LivenessProbeTimeout() time.Duration
}

type ProbesConfig struct {
	ReadinessPeriod  time.Duration `mapstructure:"readiness_probe_period" yaml:"readiness_probe_period"`
	ReadinessTimeout time.Duration `mapstructure:"readiness_probe_timeout" yaml:"readiness_probe_timeout"`
	LivenessPeriod   time.Duration `mapstructure:"liveness_probe_period" yaml:"liveness_probe_period"`
	LivenessTimeout  time.Duration `mapstructure:"liveness_probe_timeout" yaml:"liveness_probe_timeout"`
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

func (c ProbesConfig) LivenessProbePeriod() time.Duration {
	if c.LivenessPeriod <= 0 {
		return ProbesConfigLivenessDefaultPeriod
	}

	return c.LivenessPeriod
}

func (c ProbesConfig) LivenessProbeTimeout() time.Duration {
	return c.LivenessTimeout
}
