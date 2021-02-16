package serial

import (
	"errors"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/serial"
	"github.com/kihamo/boggart/protocols/serial_network"
)

type Bind struct {
	status uint32 // 0 - default, 1 - started, 2 - failed

	di.ConfigBind
	di.LoggerBind
	di.ProbesBind

	server serialnetwork.Server
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	opts := []serial.Option{
		serial.WithAddress(cfg.Target),
		serial.WithBaudRate(cfg.BaudRate),
		serial.WithDataBits(cfg.DataBits),
		serial.WithStopBits(cfg.StopBits),
		serial.WithParity(cfg.Parity),
		serial.WithTimeout(cfg.Timeout),
		serial.WithOnce(cfg.Once),
	}

	dial := serial.Dial(opts...)
	address := net.JoinHostPort(cfg.Host, strconv.FormatInt(cfg.Port, 10))

	switch cfg.Network {
	case "tcp", "tcp4", "tcp6":
		b.server = serialnetwork.NewTCPServer(cfg.Network, address, dial)

	case "udp", "udp4", "udp6":
		b.server = serialnetwork.NewUDPServer(cfg.Network, address, dial)

	default:
		return errors.New("unsupported network " + cfg.Network)
	}

	go func() {
		atomic.StoreUint32(&b.status, 1)
		defer atomic.StoreUint32(&b.status, 2)

		if err := b.server.ListenAndServe(); err != nil {
			b.Logger().Error("Failed serve with error " + err.Error())
		}
	}()

	return nil
}

func (b *Bind) Close() (err error) {
	return b.server.Close()
}
