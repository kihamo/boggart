package z_stack

type Watcher struct {
	frames chan *Frame
	errors chan error
	done   chan struct{}
}

func NewWatcher() *Watcher {
	return &Watcher{
		frames: make(chan *Frame),
		errors: make(chan error),
		done:   make(chan struct{}),
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

func (w *Watcher) close() {
	close(w.done)
}

func (w *Watcher) NextFrame() <-chan *Frame {
	return w.frames
}

func (w *Watcher) NextError() <-chan error {
	return w.errors
}

func (w *Watcher) Done() <-chan struct{} {
	return w.done
}
