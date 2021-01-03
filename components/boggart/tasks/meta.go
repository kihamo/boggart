package tasks

import (
	"time"

	"github.com/kihamo/boggart/atomic"
)

type Meta struct {
	id         string
	firstRunAt atomic.TimeNull
	lastRunAt  atomic.TimeNull

	attempts atomic.Uint64
	fails    atomic.Uint64
	status   atomic.Uint32
}

func (m Meta) ID() string {
	return m.id
}

func (m Meta) FirstRunAt() *time.Time {
	return m.firstRunAt.Load()
}

func (m Meta) LastRunAt() *time.Time {
	return m.lastRunAt.Load()
}

func (m Meta) Attempts() uint64 {
	return m.attempts.Load()
}

func (m Meta) Success() uint64 {
	return m.Attempts() - m.Fails()
}

func (m Meta) Fails() uint64 {
	return m.fails.Load()
}

func (m Meta) Status() Status {
	return Status(m.status.Load())
}
