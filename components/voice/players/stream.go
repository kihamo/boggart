package players

import (
	"io"
	"sync/atomic"
)

type Stream interface {
	io.ReadCloser

	SampleRate() int
}

type StreamWrapper struct {
	closed int64

	original Stream
}

func NewStreamWrapper(original Stream) *StreamWrapper {
	return &StreamWrapper{
		original: original,
	}
}

func (w *StreamWrapper) Read(p []byte) (int, error) {
	if w == nil {
		return -1, nil
	}

	return w.original.Read(p)
}

func (w *StreamWrapper) Close() error {
	if w == nil {
		return nil
	}

	err := w.original.Close()
	if err == nil {
		atomic.StoreInt64(&w.closed, 1)
	}

	return err
}

func (w *StreamWrapper) SampleRate() int {
	if w == nil {
		return 0
	}

	return w.original.SampleRate()
}

func (w *StreamWrapper) IsClosed() bool {
	if w == nil {
		return true
	}

	return atomic.LoadInt64(&w.closed) == 1
}
