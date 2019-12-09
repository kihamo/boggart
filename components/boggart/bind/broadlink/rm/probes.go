package rm

import (
	"context"
	"errors"
	"time"

	"github.com/sparrc/go-ping"
)

const (
	overhead = time.Millisecond * 100
)

func (b *Bind) ReadinessProbe(ctx context.Context) error {
	pinger, err := ping.NewPinger(b.config.Host)
	if err != nil {
		return err
	}

	pinger.SetPrivileged(true)

	deadline, ok := ctx.Deadline()
	if !ok {
		return errors.New("get deadline failed")
	}

	pinger.Timeout = time.Until(deadline)
	if pinger.Timeout > overhead {
		pinger.Timeout -= overhead
	}

	pinger.Run()
	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		return errors.New("packets is zero")
	}

	return nil
}
