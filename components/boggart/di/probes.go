package di

import (
	"context"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
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
		return support.Probes(), true
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

	probeReadiness workers.Task
	probeLiveness  workers.Task
}

func NewProbesContainer(bind boggart.BindItem, statusManager func(boggart.BindStatus), register, unregister func() error) *ProbesContainer {
	return &ProbesContainer{
		bind:          bind,
		statusManager: statusManager,
		register:      register,
		unregister:    unregister,
	}
}

func (c *ProbesContainer) Readiness() workers.Task {
	if c.probeReadiness != nil {
		return c.probeReadiness
	}

	has, ok := c.bind.Bind().(ProbesHasReadinessProbe)
	if !ok {
		return nil
	}

	logger, _ := LoggerContainerBind(c.bind.Bind())

	probeTask := task.NewFunctionTask(func(ctx context.Context) (_ interface{}, err error) {
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

			if logger != nil {
				logger.Error("Readiness probe failure",
					"type", c.bind.Type(),
					"id", c.bind.ID(),
					"error", err.Error(),
				)
			}
		} else {
			c.statusManager(boggart.BindStatusOnline)
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
			ctx, _ = context.WithTimeout(ctx, timeout)
		}

		_, err = c.probeReadiness.Run(ctx)
	}

	return err
}

func (c *ProbesContainer) Liveness() workers.Task {
	if c.probeLiveness != nil {
		return c.probeLiveness
	}

	has, ok := c.bind.Bind().(ProbesHasLivenessProbe)
	if !ok {
		return nil
	}

	logger, _ := LoggerContainerBind(c.bind.Bind())

	probeTask := task.NewFunctionTask(func(ctx context.Context) (_ interface{}, err error) {
		currentStatus := c.bind.Status()
		if currentStatus != boggart.BindStatusOnline && currentStatus != boggart.BindStatusOffline {
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
			return nil, nil
		}

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
			ctx, _ = context.WithTimeout(ctx, timeout)
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
