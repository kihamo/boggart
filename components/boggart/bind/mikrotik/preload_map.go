package mikrotik

import (
	"sync"

	"github.com/kihamo/boggart/atomic"
)

type PreloadMap struct {
	sync.Map

	readyNotify chan struct{}
	ready       *atomic.Bool
}

func NewPreloadMap() *PreloadMap {
	return &PreloadMap{
		readyNotify: make(chan struct{}),
		ready:       atomic.NewBoolDefault(false),
	}
}

func (s *PreloadMap) Ready() {
	if s.ready.IsFalse() {
		s.ready.True()
		close(s.readyNotify)
	}
}

func (s *PreloadMap) IsReady() bool {
	return s.ready.IsTrue()
}

func (s *PreloadMap) LoadWait(key interface{}) (value interface{}, ok bool) {
	<-s.readyNotify
	return s.Load(key)
}
