package connection

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type DSN struct {
	url.URL

	ReadTimeout  *time.Duration
	WriteTimeout *time.Duration
	Timeout      *time.Duration
	OnceInit     *bool
	LockLocal    *bool
	LockGlobal   *bool
	Dump         *bool
	BaudRate     *int64
	DataBits     *int64
	StopBits     *int64
	Parity       *string
}

func ParseDSN(dsn string) (*DSN, error) {
	if dsn == "" {
		return nil, errors.New("DSN is empty")
	}

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	d := &DSN{
		URL: *u,
	}

	for key, value := range u.Query() {
		switch strings.ToLower(key) {
		case "read-timeout":
			v, err := time.ParseDuration(value[0])
			if err != nil {
				return nil, err
			}

			d.ReadTimeout = &v

		case "write-timeout":
			v, err := time.ParseDuration(value[0])
			if err != nil {
				return nil, err
			}

			d.WriteTimeout = &v

		case "timeout":
			v, err := time.ParseDuration(value[0])
			if err != nil {
				return nil, err
			}

			d.Timeout = &v

		case "once":
			v, err := strconv.ParseBool(value[0])
			if err != nil {
				return nil, err
			}

			d.OnceInit = &v

		case "lock", "lock-global":
			v, err := strconv.ParseBool(value[0])
			if err != nil {
				return nil, err
			}

			d.LockGlobal = &v

		case "lock-local":
			v, err := strconv.ParseBool(value[0])
			if err != nil {
				return nil, err
			}

			d.LockLocal = &v

		case "dump", "debug":
			v, err := strconv.ParseBool(value[0])
			if err != nil {
				return nil, err
			}

			d.Dump = &v

		case "baudrate", "baud-rate":
			v, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return nil, err
			}

			d.BaudRate = &v

		case "databits", "data-bits":
			v, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return nil, err
			}

			d.DataBits = &v

		case "stopbits", "stop-bits":
			v, err := strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				return nil, err
			}

			d.StopBits = &v

		case "parity":
			v := strings.ToUpper(value[0])

			switch v {
			case "N", "E", "O":
				// skip
			default:
				return nil, errors.New("parity value " + v + " is wrong")
			}

			d.Parity = &v
		}
	}

	return d, nil
}
