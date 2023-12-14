package swagger

import (
	"bytes"
	"context"
	"time"

	"github.com/go-openapi/strfmt"
)

const (
	Format = "02.01.2006"
)

var null = []byte("null")

type Date struct {
	strfmt.Date
}

func (m Date) String() string {
	return m.Time().Format(Format)
}

func (m *Date) Time() time.Time {
	return time.Time(m.Date)
}

func (m *Date) Validate(strfmt.Registry) error {
	return nil
}

func (m *Date) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}

func (m *Date) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, null) {
		return nil
	}

	t, err := time.Parse(Format, string(bytes.Trim(b, "\"")))
	if err != nil {
		return err
	}

	m.Date = strfmt.Date(t)

	return nil
}
