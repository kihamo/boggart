package z_stack

type Watcher struct {
	frames chan *Frame
	errors chan error
}

func NewWatcher() *Watcher {
	return &Watcher{
		frames: make(chan *Frame),
		errors: make(chan error),
	}
}

func (w *Watcher) notifyFrames(frames ...*Frame) {
	go func() {
		for _, frame := range frames {
			w.frames <- frame
		}
	}()
}

func (w *Watcher) notifyError(err error) {
	go func() {
		w.errors <- err
	}()
}

func (w *Watcher) NextFrame() <-chan *Frame {
	return w.frames
}

func (w *Watcher) NextError() <-chan error {
	return w.errors
}
