package swagger

import (
	"fmt"

	o "github.com/go-openapi/runtime/logger"
)

const (
	cutSize = 1024
)

type logger struct {
	o.Logger

	f1 func(string)
	f2 func(string)

	cutSize int
}

func NewLogger(printf func(string), debug func(string)) o.Logger {
	return NewLoggerWithCutSize(printf, debug, cutSize)
}

func NewLoggerWithCutSize(printf func(string), debug func(string), cutSize int) o.Logger {
	return &logger{
		f1:      printf,
		f2:      debug,
		cutSize: cutSize,
	}
}

func (l logger) Printf(format string, args ...interface{}) {
	record := fmt.Sprintf(format, args...)

	if l.cutSize > 0 {
		cut := len(record)
		if cut > l.cutSize {
			cut = l.cutSize
		}

		l.f1(record[:cut])
	} else {
		l.f1(record)
	}
}

func (l logger) Debugf(format string, args ...interface{}) {
	record := fmt.Sprintf(format, args...)

	if l.cutSize > 0 {
		cut := len(record)
		if cut > l.cutSize {
			cut = l.cutSize
		}

		l.f2(record[:cut])
	} else {
		l.f2(record)
	}
}
