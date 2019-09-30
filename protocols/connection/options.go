package connection

import (
	"time"
)

type options struct {
	network      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	once         bool
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

func WithNetwork(network string) Option {
	return newFuncOption(func(o *options) {
		o.network = network
	})
}

func WithReadTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *options) {
		o.readTimeout = timeout
	})
}

func WithWriteTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *options) {
		o.writeTimeout = timeout
	})
}

func WithOnce(once bool) Option {
	return newFuncOption(func(o *options) {
		o.once = once
	})
}
