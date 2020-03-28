package alsa

import (
	"sync/atomic"

	"github.com/hajimehoshi/oto"
)

type speakerWrapper struct {
	closed     int64
	bufferSize int
	numBytes   int
	player     *oto.Player
}

func newSpeakerWrapper(player *oto.Player, bufferSize, numBytes int) *speakerWrapper {
	return &speakerWrapper{
		bufferSize: bufferSize,
		numBytes:   numBytes,
		player:     player,
	}
}

func (w *speakerWrapper) BufferSize() int {
	return w.bufferSize
}

func (w *speakerWrapper) NumBytes() int {
	return w.numBytes
}

func (w *speakerWrapper) Write(data []byte) (int, error) {
	return w.player.Write(data)
}

func (w *speakerWrapper) Close() error {
	if w == nil {
		return nil
	}

	err := w.player.Close()
	if err == nil {
		atomic.StoreInt64(&w.closed, 1)
	}

	return err
}

func (w *speakerWrapper) IsClosed() bool {
	if w == nil {
		return true
	}

	return atomic.LoadInt64(&w.closed) == 1
}
