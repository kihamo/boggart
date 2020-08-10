package net

import (
	"time"
)

type options struct {
	network      string
	address      string
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (o *options) Map() map[string]interface{} {
	return map[string]interface{}{
		"network":      o.network,
		"address":      o.address,
		"readTimeout":  o.readTimeout,
		"writeTimeout": o.writeTimeout,
	}
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

func DefaultOptions() options {
	return options{
		network: "tcp",
	}
}

func WithNetwork(network string) Option {
	return newFuncOption(func(o *options) {
		o.network = network
	})
}

func WithAddress(address string) Option {
	return newFuncOption(func(o *options) {
		o.address = address
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

func WithTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *options) {
		o.readTimeout = timeout
		o.writeTimeout = timeout
	})
}
