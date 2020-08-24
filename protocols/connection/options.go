package connection

type Options struct {
	onceInit   bool
	lockLocal  bool
	lockGlobal bool
	dumpRead   func([]byte)
	dumpWrite  func([]byte)
	readCheck  func([]byte) bool
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
	return Options{}
}

func WithOnceInit(flag bool) Option {
	return newFuncOption(func(o *Options) {
		o.onceInit = flag
	})
}

func WithLock(flag bool) Option {
	return WithGlobalLock(flag)
}

func WithGlobalLock(flag bool) Option {
	return newFuncOption(func(o *Options) {
		o.lockGlobal = flag
	})
}

func WithLocalLock(flag bool) Option {
	return newFuncOption(func(o *Options) {
		o.lockLocal = flag
	})
}

func WithDump(dump func([]byte)) Option {
	return newFuncOption(func(o *Options) {
		o.dumpRead = dump
		o.dumpWrite = dump
	})
}

func WithDumpRead(dump func([]byte)) Option {
	return newFuncOption(func(o *Options) {
		o.dumpRead = dump
	})
}

func WithDumpWrite(dump func([]byte)) Option {
	return newFuncOption(func(o *Options) {
		o.dumpWrite = dump
	})
}

func WithReadCheck(check func([]byte) bool) Option {
	return newFuncOption(func(o *Options) {
		o.readCheck = check
	})
}
