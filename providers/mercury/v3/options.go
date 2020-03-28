package v3

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
		o.address = ConvertSerialNumber(address)
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
