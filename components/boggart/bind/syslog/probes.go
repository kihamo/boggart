package syslog

import (
	"context"
	"errors"
	"sync/atomic"
)

func (b *Bind) LivenessProbe(context.Context) error {
	if atomic.LoadUint32(&b.status) == 2 {
		return errors.New("server isn't running")
	}

	return nil
}
