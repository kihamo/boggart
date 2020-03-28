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
}

func NewLogger(printf func(string), debug func(string)) o.Logger {
	return &logger{
		f1: printf,
		f2: debug,
	}
}

func (l logger) Printf(format string, args ...interface{}) {
	record := fmt.Sprintf(format, args...)

	cut := len(record)
	if cut > cutSize {
		cut = cutSize
	}

	l.f1(record[:cut])
}

func (l logger) Debugf(format string, args ...interface{}) {
	record := fmt.Sprintf(format, args...)

	cut := len(record)
	if cut > cutSize {
		cut = cutSize
	}

	l.f2(record[:cut])
}
