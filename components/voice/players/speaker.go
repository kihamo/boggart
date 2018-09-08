package players

import (
	"io"
	"sync/atomic"
)

type SpeakerWrapper struct {
	closed int64

	original io.WriteCloser
}

func NewSpeakerWrapper(original io.WriteCloser) *SpeakerWrapper {
	return &SpeakerWrapper{
		original: original,
	}
}

func (w *SpeakerWrapper) Write(p []byte) (int, error) {
	if w == nil {
		return -1, nil
	}

	return w.original.Write(p)
}

func (w *SpeakerWrapper) Close() error {
	if w == nil {
		return nil
	}

	err := w.original.Close()
	if err == nil {
		atomic.StoreInt64(&w.closed, 1)
	}

	return err
}

func (w *SpeakerWrapper) IsClosed() bool {
	if w == nil {
		return true
	}

	return atomic.LoadInt64(&w.closed) == 1
}
