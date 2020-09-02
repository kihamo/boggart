package net

import (
	"time"
)

type Options struct {
	network      string
	address      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	dialTimeout  time.Duration
}

func (o *Options) Map() map[string]interface{} {
	return map[string]interface{}{
		"network":      o.network,
		"address":      o.address,
		"readTimeout":  o.readTimeout,
		"writeTimeout": o.writeTimeout,
		"dialTimeout":  o.dialTimeout,
	}
}

type Option interface {
	apply(*Options)
}

type funcOption struct {
	f func(*Options)
}

func (fdo *funcOption) apply(do *Options) {
	fdo.f(do)
}

func newFuncOption(f func(*Options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func DefaultOptions() Options {
	return Options{
		network: "tcp",
	}
}

func WithNetwork(network string) Option {
	return newFuncOption(func(o *Options) {
		o.network = network
	})
}

func WithAddress(address string) Option {
	return newFuncOption(func(o *Options) {
		o.address = address
	})
}

func WithReadTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *Options) {
		o.readTimeout = timeout
	})
}

func WithWriteTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *Options) {
		o.writeTimeout = timeout
	})
}

func WithDialTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *Options) {
		o.dialTimeout = timeout
	})
}

func WithTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *Options) {
		o.readTimeout = timeout
		o.writeTimeout = timeout
	})
}
