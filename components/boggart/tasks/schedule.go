package tasks

import (
	"time"
)

var (
	background = ScheduleFunc(func(Meta) time.Time {
		return time.Now()
	})
)

type Schedule interface {
	Next(Meta) time.Time
}

func ScheduleBackground() Schedule {
	return background
}

type ScheduleFunc func(Meta) time.Time

func (f ScheduleFunc) Next(meta Meta) time.Time {
	return f(meta)
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
	return ScheduleFunc(func(meta Meta) (t time.Time) {
		if parent != nil {
			t = parent.Next(meta)
		}

		next := time.Now().Add(duration)
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
