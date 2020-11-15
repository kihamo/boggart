package zigbee2mqtt

import (
	"context"
	"errors"
	"time"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	if b.config.NewAPI {
		err = b.MQTT().PublishWithoutCache(ctx, b.config.TopicHealthCheckRequest, true)
		if err != nil {
			return err
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(time.Second * 5)
		}

		time.Sleep(deadline.Sub(time.Now()) - time.Millisecond*500)
	}

	if b.status.IsTrue() {
		return nil
	}

	if b.config.NewAPI {
		return errors.New("health check failed")
	}

	return errors.New("status isn't online")
}
