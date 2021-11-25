package swagger

import (
	"time"

	"github.com/go-openapi/strfmt"
)

type Date struct {
	strfmt.Date
}

func (m *Date) Time() time.Time {
	return time.Time(m.Date)
}

func (m *Date) Validate(formats strfmt.Registry) error {
	return nil
}

func (m *Date) UnmarshalJSON(b []byte) error {
	return m.Date.UnmarshalJSON(b)
}
