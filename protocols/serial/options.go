package serial

import (
	"time"

	s "github.com/goburrow/serial"
)

type Options struct {
	s.Config
	allowMultiRequest bool
	once              bool
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
	return newFuncOption(func(o *Options) {
		o.Config.Address = address
	})
}

func WithBaudRate(baudRate int) Option {
	return newFuncOption(func(o *Options) {
		o.Config.BaudRate = baudRate
	})
}

func WithDataBits(dataBits int) Option {
	return newFuncOption(func(o *Options) {
		o.Config.DataBits = dataBits
	})
}

func WithStopBits(stopBits int) Option {
	return newFuncOption(func(o *Options) {
		o.Config.StopBits = stopBits
	})
}

func WithParity(parity string) Option {
	return newFuncOption(func(o *Options) {
		o.Config.Parity = parity
	})
}

func WithTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *Options) {
		o.Config.Timeout = timeout
	})
}
func WithOnce(once bool) Option {
	return newFuncOption(func(o *Options) {
		o.once = once
	})
}
