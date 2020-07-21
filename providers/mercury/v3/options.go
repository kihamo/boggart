package v3

import (
	"strconv"
)

type options struct {
	address     byte
	password    LevelPassword
	accessLevel accessLevel
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
		address:     0x0,
		accessLevel: AccessLevel1,
		password:    DefaultPasswordLevel1,
	}
}

func WithAddress(address byte) Option {
	return newFuncOption(func(o *options) {
		o.address = address
	})
}

func WithAddressAsString(address string) Option {
	return newFuncOption(func(o *options) {
		// Отделите три последние цифры серийного номера, это будет число N.
		number, _ := strconv.ParseInt(address[len(address)-3:], 10, 0)

		// Если N>=240 адресом являются две последние цифры серийного номера.
		if number >= 240 {
			number, _ = strconv.ParseInt(address[len(address)-2:], 10, 0)

			o.address = byte(number)
			return
		}

		// Если N<240 адресом являются три последние цифры.
		if number < 240 {
			o.address = byte(number)
			return
		}

		o.address = 1
	})
}

func WithAccessLevel(accessLevel accessLevel) Option {
	return newFuncOption(func(o *options) {
		o.accessLevel = accessLevel
	})
}

func WithPasswordLevel(password LevelPassword) Option {
	return newFuncOption(func(o *options) {
		o.password = password
	})
}
