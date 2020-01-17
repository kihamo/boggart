package rkcm

import (
	"context"
	"net"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/probes"
	"github.com/kihamo/boggart/providers/rkcm"
)

var address = "http://" + net.JoinHostPort(rkcm.DefaultHost, strconv.FormatInt(rkcm.DefaultPort, 10)) + "/"

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	return probes.HTTPProbe(ctx, address, nil)
}
