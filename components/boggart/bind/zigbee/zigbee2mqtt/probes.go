package zigbee2mqtt

import (
	"context"
	"errors"
	"time"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	cfg := b.config()

	if cfg.NewAPI {
		err = b.MQTT().PublishRawWithoutCache(ctx, cfg.TopicHealthCheckRequest, 1, false, true)
		if err != nil {
			return err
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(time.Second * 5)
		}

		time.Sleep(time.Until(deadline) - time.Millisecond*500)
	}

	if b.status.IsTrue() {
		return nil
	}

	if cfg.NewAPI {
		return errors.New("health check failed")
	}

	return errors.New("status isn't online")
}
