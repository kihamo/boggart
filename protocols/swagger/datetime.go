package swagger

import (
	"bytes"
	"context"
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

func (m *DateTime) Validate(strfmt.Registry) error {
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

func (m *DateTime) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}

type DateTimeByTimestamp struct {
	strfmt.DateTime
}

func (m *DateTimeByTimestamp) Time() time.Time {
	return time.Time(m.DateTime)
}

func (m *DateTimeByTimestamp) Validate(strfmt.Registry) error {
	return nil
}

func (m *DateTimeByTimestamp) UnmarshalJSON(b []byte) error {
	sec, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	m.DateTime = strfmt.DateTime(time.Unix(int64(sec), 0))

	return nil
}

func (m *DateTimeByTimestamp) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}
