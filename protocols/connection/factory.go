package connection

import (
	"errors"
	"fmt"

	"github.com/kihamo/boggart/protocols/connection/transport"
	"github.com/kihamo/boggart/protocols/connection/transport/net"
	"github.com/kihamo/boggart/protocols/connection/transport/serial"
)

func NewByDSNString(dsn string) (conn Connection, err error) {
	d, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	return NewByDSN(d)
}

func NewByDSN(dsn *DSN) (Connection, error) {
	var t transport.Transport

	switch dsn.Scheme {
	case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6", "unixgram":
		options := []net.Option{
			net.WithNetwork(dsn.Scheme),
			net.WithAddress(dsn.Host),
		}

		if dsn.Timeout != nil {
			options = append(options, net.WithTimeout(*dsn.Timeout))
		}

		if dsn.ReadTimeout != nil {
			options = append(options, net.WithReadTimeout(*dsn.ReadTimeout))
		}

		if dsn.WriteTimeout != nil {
			options = append(options, net.WithWriteTimeout(*dsn.WriteTimeout))
		}

		t = net.New(options...)

	case "serial":
		options := []serial.Option{
			serial.WithAddress(dsn.EscapedPath()),
		}

		if dsn.BaudRate != nil {
			options = append(options, serial.WithBaudRate(int(*dsn.BaudRate)))
		}

		if dsn.DataBits != nil {
			options = append(options, serial.WithDataBits(int(*dsn.DataBits)))
		}

		if dsn.StopBits != nil {
			options = append(options, serial.WithStopBits(int(*dsn.StopBits)))
		}

		if dsn.Parity != nil {
			options = append(options, serial.WithParity(*dsn.Parity))
		}

		if dsn.Timeout != nil {
			options = append(options, serial.WithTimeout(*dsn.Timeout))
		}

		t = serial.New(options...)

	default:
		return nil, errors.New("unknown connection type for DSN " + dsn.String())
	}

	// Wrappers
	options := make([]Option, 0)

	if dsn.OnceInit != nil {
		options = append(options, WithOnceInit(*dsn.OnceInit))
	}

	if dsn.LockLocal != nil {
		options = append(options, WithLocalLock(*dsn.LockLocal))
	}

	if dsn.LockGlobal != nil {
		options = append(options, WithGlobalLock(*dsn.LockGlobal))
	}

	if dsn.Dump != nil && *dsn.Dump {
		options = append(options, WithDump(func(bytes []byte) {
			fmt.Println(bytes)
		}))
	}

	return New(t, options...), nil
}
