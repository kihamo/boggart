package v3

import (
	"strconv"
	"strings"
)

type options struct {
	address     uint8
	password    LevelPassword
	accessLevel uint8
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
		address:     AddressUniveral,
		accessLevel: AccessLevel1,
		password:    DefaultPasswordLevel1,
	}
}

func WithAddress(address uint8) Option {
	return newFuncOption(func(o *options) {
		o.address = address
	})
}

func WithAddressAsString(address string) Option {
	return newFuncOption(func(o *options) {
		// Отделите три последние цифры серийного номера, это будет число N.
		clean := strings.Split(address, "-")
		address = clean[0]

		number, _ := strconv.ParseInt(address[len(address)-3:], 10, 0)

		// Если N>=240 адресом являются две последние цифры серийного номера.
		if number >= 240 {
			number, _ = strconv.ParseInt(address[len(address)-2:], 10, 0)

			o.address = uint8(number)
			return
		}

		// Если N<240 адресом являются три последние цифры.
		if number < 240 {
			o.address = uint8(number)
			return
		}

		o.address = 1
	})
}

func WithAccessLevel(accessLevel uint8) Option {
	return newFuncOption(func(o *options) {
		o.accessLevel = accessLevel
	})
}

func WithPasswordLevel(password LevelPassword) Option {
	return newFuncOption(func(o *options) {
		o.password = password
	})
}
