package swagger

import (
	"time"

	"github.com/go-openapi/strfmt"
)

const (
	Format = "02.01.2006"
)

type Date struct {
	strfmt.Date
}

func (m Date) String() string {
	return m.Time().Format(Format)
}

func (m *Date) Time() time.Time {
	return time.Time(m.Date)
}

func (m *Date) Validate(formats strfmt.Registry) error {
	return nil
}

func (m *Date) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(Format, string(b))
	if err != nil {
		return err
	}

	m.Date = strfmt.Date(t)

	return nil
}
