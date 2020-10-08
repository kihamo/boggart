package swagger

import (
	"bytes"
	"time"

	"github.com/go-openapi/strfmt"
)

type DateTime struct {
	strfmt.DateTime
}

func (m *DateTime) Time() time.Time {
	return time.Time(m.DateTime)
}

func (m *DateTime) Validate(formats strfmt.Registry) error {
	return nil
}

func (m *DateTime) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02 15:04:05", string(bytes.Trim(b, "\"")))
	if err != nil {
		return err
	}

	m.DateTime = strfmt.DateTime(t)

	return nil
}
