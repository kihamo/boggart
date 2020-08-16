package z_stack

type options struct {
	Channel              uint32
	panID                uint16
	extendedPanID        []byte
	networkKeyDistribute bool
	networkKey           []byte
	LEDEnabled           bool
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
		Channel:              11,
		panID:                0x1A62,
		extendedPanID:        []byte{0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD},
		networkKeyDistribute: false,
		networkKey:           []byte{0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F, 0x00, 0x02, 0x04, 0x06, 0x08, 0x0A, 0x0C, 0x0D},
		LEDEnabled:           true,
	}
}

func WithChannel(channel uint32) Option {
	return newFuncOption(func(o *options) {
		o.Channel = channel
	})
}

func WithPanID(panID uint16) Option {
	return newFuncOption(func(o *options) {
		o.panID = panID
	})
}

func WithExtendedPanID(extendedPanID []byte) Option {
	return newFuncOption(func(o *options) {
		o.extendedPanID = extendedPanID
	})
}

func WithNetworkKeyDistribute(flag bool) Option {
	return newFuncOption(func(o *options) {
		o.networkKeyDistribute = flag
	})
}

func WithNetworkKey(networkKey []byte) Option {
	return newFuncOption(func(o *options) {
		o.networkKey = networkKey
	})
}

func WithLEDEnabled(flag bool) Option {
	return newFuncOption(func(o *options) {
		o.LEDEnabled = flag
	})
}
