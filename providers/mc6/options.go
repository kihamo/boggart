package mc6

import (
	"io"
	"time"
)

type options struct {
	slaveID     uint8
	timeout     time.Duration
	idleTimeout time.Duration
	maxTries    uint8
	logger      io.Writer
}

type Option interface {
	apply(*options)
}

type funcOption struct {
	f func(*options)
}

func (fdo *funcOption) apply(do *options) {
	fdo.f(do)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func defaultOptions() options {
	return options{
		slaveID:     0x1,
		timeout:     time.Second,
		idleTimeout: time.Minute,
		maxTries:    3,
	}
}

func WithSlaveID(id uint8) Option {
	return newFuncOption(func(o *options) {
		o.slaveID = id
	})
}

func WithTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *options) {
		o.timeout = timeout
	})
}

func WithIdleTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *options) {
		o.idleTimeout = timeout
	})
}

func WithMaxTries(tries uint8) Option {
	return newFuncOption(func(o *options) {
		if tries == 0 {
			tries = 1
		}

		o.maxTries = tries
	})
}

func WithLogger(logger io.Writer) Option {
	return newFuncOption(func(o *options) {
		o.logger = logger
	})
}
