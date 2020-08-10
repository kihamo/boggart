package serial

import (
	"time"

	s "github.com/goburrow/serial"
)

type options struct {
	s.Config
}

func (o *options) Map() map[string]interface{} {
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
