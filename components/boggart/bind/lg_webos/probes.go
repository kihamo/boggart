package lg_webos

import (
	"context"
)

func (b *Bind) LivenessProbe(ctx context.Context) error {
	if client := b.Client(); client != nil {
		_, _, _, _, _, _, err := client.GetCurrentTime()
		return err
	}

	return nil
}

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	client := b.Client()

	if client == nil {
		err = b.initClient()
	} else {
		_, _, _, _, _, _, err = client.GetCurrentTime()
	}

	return err
}
