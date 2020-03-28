package scale

import (
	"context"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/go-ble/ble"
)

const (
	DefaultScanTimeout = time.Second * 10
)

type Client struct {
	addr                 net.HardwareAddr
	device               ble.Device
	duration             time.Duration
	scanResult           chan []byte
	scanError            chan error
	ignoreEmptyImpedance bool
}

var scaleUUID = ble.UUID16(0x181B)

func NewClient(device ble.Device, addr net.HardwareAddr, duration time.Duration, ignoreEmptyImpedance bool) *Client {
	return &Client{
		addr:                 addr,
		device:               device,
		duration:             duration,
		scanResult:           make(chan []byte),
		scanError:            make(chan error, 1),
		ignoreEmptyImpedance: ignoreEmptyImpedance,
	}
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

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, DefaultScanTimeout)
	}

	defer cancel()

	var (
		wg  sync.WaitGroup
		err error
	)

	measuresCache := make(map[time.Time]*Measure)

	wg.Add(1)
	go func() {
		defer wg.Done()

		e := s.device.Scan(ctx, false, s.advHandler(s.scanResult))
		if e != nil && e != context.DeadlineExceeded && e != context.Canceled {
			s.scanError <- e
		}
	}()

SCAN:
	for {
		select {
		case data := <-s.scanResult:
			if len(data) == 0 {
				break SCAN
			}

			// в v2 impedance равен 0 в промежуточных результах взвешивания, поэтому такое значение можно игнорировать
			impedance := uint64(data[10])*256 + uint64(data[9])
			if impedance == 0 && s.ignoreEmptyImpedance {
				continue
			}

			m := NewMeasure(
				time.Date(int(data[3])*256+int(data[2]), time.Month(int(data[4])), int(data[5]), int(data[6]), int(data[7]), int(data[8]), 0, time.Local),
				Unit(data[0]),
				((float64(data[12])*256)+float64(data[11]))*0.01,
				impedance)

			if m.unit == UnitKG || m.unit == UnitKG2 {
				m.weight = m.weight / 2
			}

			measuresCache[m.datetime] = m

		case err = <-s.scanError:
			break SCAN

		case <-ctx.Done():
			err = ctx.Err()

			break SCAN
		}
	}

	wg.Wait()

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

func (s *Client) Close() (err error) {
	err = s.device.Stop()

	if err == nil {
		close(s.scanResult)
		close(s.scanError)
	}

	return err
}
