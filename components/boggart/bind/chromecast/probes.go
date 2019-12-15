package chromecast

import (
	"context"
	"errors"
)

func (b *Bind) LivenessProbe(_ context.Context) error {
	if b.disconnected.IsTrue() {
		return errors.New("disconnected")
	}

	return nil
}

func (b *Bind) ReadinessProbe(_ context.Context) (err error) {
	if b.disconnected.IsNil() {
		err = b.initConnect()
	} else if b.disconnected.IsTrue() {
		err = errors.New("disconnected")
	}

	return err
}
