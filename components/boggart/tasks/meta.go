package tasks

import (
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/pborman/uuid"
)

type Meta struct {
	id              string
	lastRunDuration *atomic.DurationNull
	firstRunAt      *atomic.TimeNull
	lastRunAt       *atomic.TimeNull
	nextRunAt       *atomic.TimeNull
	attempts        *atomic.Uint64
	fails           *atomic.Uint64
	status          *atomic.Uint32
}

func newMeta() *Meta {
	return &Meta{
		id:              uuid.New(),
		lastRunDuration: atomic.NewDurationNull(),
		firstRunAt:      atomic.NewTimeNull(),
		lastRunAt:       atomic.NewTimeNull(),
		nextRunAt:       atomic.NewTimeNull(),
		attempts:        atomic.NewUint64(),
		fails:           atomic.NewUint64(),
		status:          atomic.NewUint32(),
	}
}

func (m Meta) ID() string {
	return m.id
}

func (m Meta) LastRunDuration() *time.Duration {
	return m.lastRunDuration.Load()
}

func (m Meta) FirstRunAt() *time.Time {
	return m.firstRunAt.Load()
}

func (m Meta) LastRunAt() *time.Time {
	return m.lastRunAt.Load()
}

func (m Meta) NextRunAt() *time.Time {
	return m.nextRunAt.Load()
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
