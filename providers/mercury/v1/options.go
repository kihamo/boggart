package v1

import (
	"time"
)

type options struct {
	address  []byte
	location *time.Location
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
		location: time.Now().Location(),
	}
}

func WithAddress(address []byte) Option {
	return newFuncOption(func(o *options) {
		o.address = address
	})
}

func WithAddress200AsString(address string) Option {
	return newFuncOption(func(o *options) {
		o.address = ConvertSerialNumber200(address)
	})
}

func WithAddressAsString(address string) Option {
	return newFuncOption(func(o *options) {
		o.address = ConvertSerialNumber(address)
	})
}

func WithLocation(location *time.Location) Option {
	return newFuncOption(func(o *options) {
		o.location = location
	})
}