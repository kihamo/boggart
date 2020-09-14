package models

import (
	"strconv"
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
	sec, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	m.DateTime = strfmt.DateTime(time.Unix(int64(sec), 0))

	return nil
}
