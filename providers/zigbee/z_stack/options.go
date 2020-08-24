package zstack

type Options struct {
	extendedPanID        []byte
	networkKey           []byte
	Channel              uint32
	panID                uint16
	networkKeyDistribute bool
	LEDEnabled           bool
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
		Channel:              11,
		panID:                0x1A62,
		extendedPanID:        []byte{0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD, 0xDD},
		networkKeyDistribute: false,
		networkKey:           []byte{0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F, 0x00, 0x02, 0x04, 0x06, 0x08, 0x0A, 0x0C, 0x0D},
		LEDEnabled:           true,
	}
}

func WithChannel(channel uint32) Option {
	return newFuncOption(func(o *Options) {
		o.Channel = channel
	})
}

func WithPanID(panID uint16) Option {
	return newFuncOption(func(o *Options) {
		o.panID = panID
	})
}

func WithExtendedPanID(extendedPanID []byte) Option {
	return newFuncOption(func(o *Options) {
		o.extendedPanID = extendedPanID
	})
}

func WithNetworkKeyDistribute(flag bool) Option {
	return newFuncOption(func(o *Options) {
		o.networkKeyDistribute = flag
	})
}

func WithNetworkKey(networkKey []byte) Option {
	return newFuncOption(func(o *Options) {
		o.networkKey = networkKey
	})
}

func WithLEDEnabled(flag bool) Option {
	return newFuncOption(func(o *Options) {
		o.LEDEnabled = flag
	})
}
