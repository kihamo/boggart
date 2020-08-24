package serial

import (
	"time"

	s "github.com/goburrow/serial"
)

type Options struct {
	s.Config
}

func (o *Options) Map() map[string]interface{} {
	return map[string]interface{}{
		"address":  o.Address,
		"baudRate": o.BaudRate,
		"dataBits": o.DataBits,
		"stopBits": o.StopBits,
		"parity":   o.Parity,
		"timeout":  o.Timeout,
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
