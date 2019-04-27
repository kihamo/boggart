package serial

import (
	"time"
)

const (
	DefaultSerialAddress = "/dev/ttyUSB0"
	DefaultTimeout       = time.Second
)

type dialOptions struct {
	allowMultiRequest bool

	timeout time.Duration
}

type DialOption interface {
	apply(*dialOptions)
}

type funcDialOption struct {
	f func(*dialOptions)
}

func (fdo *funcDialOption) apply(do *dialOptions) {
	fdo.f(do)
}

func newFuncDialOption(f func(*dialOptions)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}

func defaultDialOptions() dialOptions {
	return dialOptions{
		allowMultiRequest: false,
		timeout:           DefaultTimeout,
	}
}

func WithTimeout(timeout time.Duration) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.timeout = timeout
	})
}
