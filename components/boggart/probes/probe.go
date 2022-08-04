package probes

import (
	"context"
	"errors"
	"runtime"
	"time"

	"github.com/go-ping/ping"
)

const (
	overhead = time.Millisecond * 100
)

func PingProbe(ctx context.Context, addr string) error {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return err
	}

	pinger.SetPrivileged(runtime.GOOS != "darwin")

	deadline, ok := ctx.Deadline()
	if !ok {
		return errors.New("get deadline failed")
	}

	pinger.Timeout = time.Until(deadline)
	if pinger.Timeout <= 0 {
		return errors.New("timeout value for pinger must be greater than zero")
	}

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
