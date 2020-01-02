package scale

import (
	"context"
	"net"
	"sort"
	"time"

	"github.com/go-ble/ble"
)

type Client struct {
	addr     net.HardwareAddr
	device   ble.Device
	duration time.Duration
}

var scaleUUID = ble.UUID16(0x181B)

func NewClient(device ble.Device, addr net.HardwareAddr, duration time.Duration) (*Client, error) {
	return &Client{
		addr:     addr,
		device:   device,
		duration: duration,
	}, nil
}

func (s *Client) advHandler(chResult chan []byte) func(a ble.Advertisement) {
	return func(a ble.Advertisement) {
		if a.Addr().String() != s.addr.String() || !ble.Contains(a.Services(), scaleUUID) {
			return
		}

		for _, sd := range a.ServiceData() {
			if !sd.UUID.Equal(scaleUUID) || len(sd.Data) < 13 {
				continue
			}

			chResult <- sd.Data
		}
	}
}

func (s *Client) Measures(ctx context.Context) ([]*Measure, error) {
	ctx, cancel := context.WithTimeout(ctx, s.duration)
	defer cancel()

	chResult := make(chan []byte)
	chError := make(chan error, 1)

	measuresCache := make(map[time.Time]*Measure, 0)

	defer func() {
		close(chResult)
	}()

	go func() {
		err := s.device.Scan(ctx, false, s.advHandler(chResult))
		if err != nil && err != context.DeadlineExceeded && err != context.Canceled {
			chError <- err
		}
	}()

	for {
		select {
		case data := <-chResult:
			m := NewMeasure(
				time.Date(int(data[3])*256+int(data[2]), time.Month(int(data[4])), int(data[5]), int(data[6]), int(data[7]), int(data[8]), 0, time.Local),
				Unit(data[0]),
				((float64(data[12])*256)+float64(data[11]))*0.01,
				uint64(data[10])*256+uint64(data[9]))

			if m.unit == UnitKG || m.unit == UnitKG2 {
				m.weight = m.weight / 2
			}

			measuresCache[m.datetime] = m

		case err := <-chError:
			return nil, err

		case _ = <-ctx.Done():
			err := ctx.Err()

			if err == nil || err == context.DeadlineExceeded || err == context.Canceled {
				measures := make([]*Measure, 0, len(measuresCache))

				for _, metric := range measuresCache {
					measures = append(measures, metric)
				}

				sort.SliceStable(measures, func(i, j int) bool {
					return measures[i].datetime.Before(measures[j].datetime)
				})

				return measures, nil
			}

			return nil, err
		}
	}
}

func (s *Client) Close() error {
	return s.device.Stop()
}
