package tasks

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	background = ScheduleFunc(func(Meta) (t time.Time) {
		return t
	})
)

type Schedule interface {
	Next(Meta) time.Time
}

type ScheduleFunc func(Meta) time.Time

func (f ScheduleFunc) Next(meta Meta) time.Time {
	return f(meta)
}

func ScheduleBackground() Schedule {
	return background
}

func ScheduleNow() Schedule {
	var once sync.Once

	return ScheduleFunc(func(Meta) (t time.Time) {
		once.Do(func() {
			t = time.Now()
		})

		return t
	})
}

func ScheduleStartAt(parent Schedule, startAt time.Time) Schedule {
	return ScheduleFunc(func(meta Meta) (t time.Time) {
		if startAt.After(time.Now()) {
			return startAt
		}

		return parent.Next(meta)
	})
}

func ScheduleWithStopFunc(parent Schedule, fn func(meta Meta) bool) Schedule {
	return ScheduleFunc(func(meta Meta) (t time.Time) {
		if fn(meta) {
			return t
		}

		if parent != nil {
			t = parent.Next(meta)
		}

		return t
	})
}

func ScheduleWithDuration(parent Schedule, duration time.Duration) Schedule {
	if duration <= 0 {
		return parent
	}

	return ScheduleFunc(func(meta Meta) (t time.Time) {
		if parent != nil {
			t = parent.Next(meta)
		}

		next := time.Now().Add(duration)
		if t.IsZero() || next.Before(t) {
			t = next
		}

		return t
	})
}

func ScheduleWithDailyTime(parent Schedule, hour, min, sec int, loc *time.Location) Schedule {
	return ScheduleFunc(func(meta Meta) (t time.Time) {
		if parent != nil {
			t = parent.Next(meta)
		}

		now := time.Now()
		if loc == nil {
			loc = now.Location()
		}

		next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, loc)
		if next.Before(now) {
			next = next.Add(time.Hour * 24)
		}

		if t.IsZero() || next.After(t) {
			t = next
		}

		return t
	})
}

func ScheduleWithAttemptsLimit(parent Schedule, limit uint64) Schedule {
	return ScheduleWithStopFunc(parent, func(meta Meta) bool {
		return limit > 0 && meta.Attempts() >= limit
	})
}

func ScheduleWithSuccessLimit(parent Schedule, limit uint64) Schedule {
	return ScheduleWithStopFunc(parent, func(meta Meta) bool {
		return limit > 0 && meta.Success() >= limit
	})
}

func ScheduleWithFailsLimit(parent Schedule, limit uint64) Schedule {
	return ScheduleWithStopFunc(parent, func(meta Meta) bool {
		return limit > 0 && meta.Fails() >= limit
	})
}

type ScheduleControl struct {
	Schedule

	flag uint32
}

func (s *ScheduleControl) Next(meta Meta) (t time.Time) {
	if atomic.LoadUint32(&s.flag) == 1 || s.Schedule == nil {
		return t
	}

	return s.Schedule.Next(meta)
}

func (s *ScheduleControl) Enable() {
	atomic.StoreUint32(&s.flag, 0)
}

func (s *ScheduleControl) Disable() {
	atomic.StoreUint32(&s.flag, 1)
}

func ScheduleWithControl(parent Schedule) *ScheduleControl {
	return &ScheduleControl{
		Schedule: parent,
	}
}
