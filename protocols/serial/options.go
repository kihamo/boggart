package serial

import (
	"time"

	s "github.com/goburrow/serial"
)

type options struct {
	s.Config
	allowMultiRequest bool
	once              bool
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
		allowMultiRequest: false,
		once:              true,
		Config: s.Config{
			Address:  "/dev/ttyUSB0",
			BaudRate: 9600,
			DataBits: 8,
			StopBits: 1,
			Parity:   "N",
			Timeout:  time.Second,
		},
	}
}

func WithAddress(address string) Option {
	return newFuncOption(func(o *options) {
		o.Config.Address = address
	})
}

func WithBaudRate(baudRate int) Option {
	return newFuncOption(func(o *options) {
		o.Config.BaudRate = baudRate
	})
}

func WithDataBits(dataBits int) Option {
	return newFuncOption(func(o *options) {
		o.Config.DataBits = dataBits
	})
}

func WithStopBits(stopBits int) Option {
	return newFuncOption(func(o *options) {
		o.Config.StopBits = stopBits
	})
}

func WithParity(parity string) Option {
	return newFuncOption(func(o *options) {
		o.Config.Parity = parity
	})
}

func WithTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *options) {
		o.Config.Timeout = timeout
	})
}
func WithOnce(once bool) Option {
	return newFuncOption(func(o *options) {
		o.once = once
	})
}
