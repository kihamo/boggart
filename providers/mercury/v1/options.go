package v1

import (
	"encoding/hex"
	"fmt"
	"strconv"
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
		if len(address) < 6 {
			o.address = nil
			return
		}

		number, _ := strconv.ParseInt(address[len(address)-6:], 10, 0)
		h, _ := hex.DecodeString(fmt.Sprintf("%06x", number))

		sn := make([]byte, 4)
		copy(sn[4-len(h):], h)

		o.address = sn
	})
}

func WithAddressAsString(address string) Option {
	return newFuncOption(func(o *options) {
		if len(address) < 8 {
			o.address = nil
			return
		}

		number, _ := strconv.ParseInt(address[:8], 10, 0)
		h, _ := hex.DecodeString(fmt.Sprintf("%08x", number))

		sn := make([]byte, 4)
		copy(sn[4-len(h):], h)

		o.address = sn
	})
}

func WithLocation(location *time.Location) Option {
	return newFuncOption(func(o *options) {
		o.location = location
	})
}
