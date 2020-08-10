package connection

type options struct {
	onceInit   bool
	lockLocal  bool
	lockGlobal bool
	dumpRead   func([]byte)
	dumpWrite  func([]byte)
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
	return options{}
}

func WithOnceInit(flag bool) Option {
	return newFuncOption(func(o *options) {
		o.onceInit = flag
	})
}

func WithLock(flag bool) Option {
	return WithGlobalLock(flag)
}

func WithGlobalLock(flag bool) Option {
	return newFuncOption(func(o *options) {
		o.lockGlobal = flag
	})
}

func WithLocalLock(flag bool) Option {
	return newFuncOption(func(o *options) {
		o.lockLocal = flag
	})
}

func WithDump(dump func([]byte)) Option {
	return newFuncOption(func(o *options) {
		o.dumpRead = dump
		o.dumpWrite = dump
	})
}

func WithDumpRead(dump func([]byte)) Option {
	return newFuncOption(func(o *options) {
		o.dumpRead = dump
	})
}

func WithDumpWrite(dump func([]byte)) Option {
	return newFuncOption(func(o *options) {
		o.dumpWrite = dump
	})
}
