package mikrotik

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
)

var (
	wifiClientRegexp = regexp.MustCompile(`^([^@]+)@([^:\s]+):\s+([^\s,]+)`)
	vpnClientRegexp  = regexp.MustCompile(`^(\S+) logged (in|out), (.+?)$`)
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	provider     *mikrotik.Client
	host         string
	syslogClient string

	livenessInterval time.Duration
	livenessTimeout  time.Duration
	updaterInterval  time.Duration
}

type Mac struct {
	Address string
	ARP     struct {
		IP      string
		Comment string
	}
	DHCP struct {
		Hostname string
	}
}

func (b *Bind) Mac(ctx context.Context, mac string) (*Mac, error) {
	if b.SerialNumber() == "" {
		return nil, errors.New("serial number is empty")
	}

	info := &Mac{
		Address: mac,
	}

	if table, err := b.provider.IPARP(ctx); err == nil {
		for _, row := range table {
			if row.MacAddress == mac {
				info.ARP.IP = row.Address
				info.ARP.Comment = row.Comment
				break
			}
		}
	} else {
		return nil, err
	}

	if leases, err := b.provider.IPDHCPServerLease(ctx); err == nil {
		for _, lease := range leases {
			if lease.MacAddress == mac {
				info.DHCP.Hostname = lease.MacAddress
				break
			}
		}
	} else {
		return nil, err
	}

	return info, nil
}
