package serial

import (
	"time"
)

type options struct {
	allowMultiRequest bool
	target            string
	baudRate          int
	dataBits          int
	stopBits          int
	parity            string
	timeout           time.Duration
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
		target:            "/dev/ttyUSB0",
		baudRate:          9600,
		dataBits:          8,
		stopBits:          1,
		parity:            "N",
		timeout:           time.Second,
	}
}

func WithTarget(address string) Option {
	return newFuncOption(func(o *options) {
		o.target = address
	})
}

func WithBaudRate(baudRate int) Option {
	return newFuncOption(func(o *options) {
		o.baudRate = baudRate
	})
}

func WithDataBits(dataBits int) Option {
	return newFuncOption(func(o *options) {
		o.dataBits = dataBits
	})
}

func WithStopBits(stopBits int) Option {
	return newFuncOption(func(o *options) {
		o.stopBits = stopBits
	})
}

func WithParity(parity string) Option {
	return newFuncOption(func(o *options) {
		o.parity = parity
	})
}

func WithTimeout(timeout time.Duration) Option {
	return newFuncOption(func(o *options) {
		o.timeout = timeout
	})
}
