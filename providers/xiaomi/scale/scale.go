package scale

import (
	"context"
	"net"
	"sort"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

type Scale struct {
	addr     net.HardwareAddr
	device   ble.Device
	duration time.Duration
}

const (
	DefaultDevice = "default"
)

var scaleUUID = ble.UUID16(0x181B)

func New(addr net.HardwareAddr, device string, duration time.Duration) (*Scale, error) {
	if device == "" {
		device = DefaultDevice
	}

	d, err := dev.NewDevice(device)
	if err != nil {
		return nil, err
	}

	ble.SetDefaultDevice(d)

	return &Scale{
		addr:     addr,
		device:   d,
		duration: duration,
	}, nil
}

func (s *Scale) advFilter(a ble.Advertisement) bool {
	if a.Addr().String() != s.addr.String() {
		return false
	}

	return ble.Contains(a.Services(), scaleUUID)
}

func (s *Scale) advHandler(chResult chan []byte) func(a ble.Advertisement) {
	return func(a ble.Advertisement) {
		for _, sd := range a.ServiceData() {
			if !sd.UUID.Equal(scaleUUID) || len(sd.Data) < 13 {
				continue
			}

			chResult <- sd.Data
		}
	}
}

func (s *Scale) Metrics(ctx context.Context) ([]*Metrics, error) {
	ctx, cancel := context.WithTimeout(ctx, s.duration)
	defer cancel()

	chResult := make(chan []byte)
	chError := make(chan error, 1)

	metricsCache := make(map[time.Time]*Metrics, 0)

	defer func() {
		close(chResult)
	}()

	go func() {
		err := ble.Scan(ctx, false, s.advHandler(chResult), s.advFilter)
		if err != nil && err != context.DeadlineExceeded && err != context.Canceled {
			chError <- err
		}
	}()

	for {
		select {
		case data := <-chResult:
			m := &Metrics{
				datetime: time.Date(
					int(data[3])*256+int(data[2]),
					time.Month(int(data[4])),
					int(data[5]),
					int(data[6]),
					int(data[7]),
					int(data[8]),
					0,
					time.Local),
				unit:      Unit(data[0]),
				weight:    ((float64(data[12]) * 256) + float64(data[11])) * 0.01,
				impedance: int64(data[10])*256 + int64(data[9]),
			}

			if m.unit == UnitKG || m.unit == UnitKG2 {
				m.weight = m.weight / 2
			}

			metricsCache[m.datetime] = m

		case err := <-chError:
			return nil, err

		case _ = <-ctx.Done():
			err := ctx.Err()

			if err == nil || err == context.DeadlineExceeded || err == context.Canceled {
				metrics := make([]*Metrics, 0, len(metricsCache))

				for _, metric := range metricsCache {
					metrics = append(metrics, metric)
				}

				sort.SliceStable(metrics, func(i, j int) bool {
					return metrics[i].datetime.Before(metrics[j].datetime)
				})

				return metrics, nil
			}

			return nil, err
		}
	}
}

func (s *Scale) Close() error {
	return ble.Stop()
}
