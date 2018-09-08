package players

import (
	"sync/atomic"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
)

type streamWrapper struct {
	closed int64

	source effects.Volume
	format beep.Format
}

func NewStreamWrapper(source beep.StreamSeekCloser, format beep.Format) *streamWrapper {
	return &streamWrapper{
		source: effects.Volume{
			Streamer: source,
			Base:     2,
			Volume:   -float64(100-50) / 100.0 * 5,
			Silent:   false,
		},
		format: format,
	}
}

func (w *streamWrapper) Err() error {
	return nil
}

func (w *streamWrapper) SampleRate() int {
	if w == nil {
		return -1
	}

	return int(w.format.SampleRate)
}

func (w *streamWrapper) ChannelNum() int {
	if w == nil {
		return -1
	}

	return w.format.NumChannels
}

func (w *streamWrapper) BytesPerSample() int {
	if w == nil {
		return -1
	}

	return w.format.Precision
}

func (w *streamWrapper) Volume(v int64) {
	if w == nil {
		return
	}

	w.source.Volume = -float64(100-v) / 100.0 * 5
}

func (w *streamWrapper) Stream(samples [][2]float64) (n int, ok bool) {
	if w == nil {
		return -1, false
	}

	return w.source.Stream(samples)
}

func (w *streamWrapper) Close() (err error) {
	if w == nil {
		return nil
	}

	closer, ok := w.source.Streamer.(beep.StreamCloser)
	if !ok {
		return nil
	}

	err = closer.Close()
	if err == nil {
		atomic.StoreInt64(&w.closed, 1)
	}

	return err
}

func (w *streamWrapper) IsClosed() bool {
	if w == nil {
		return true
	}

	return atomic.LoadInt64(&w.closed) == 1
}
