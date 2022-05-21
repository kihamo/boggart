package swagger

import (
	"context"
	"time"

	"github.com/go-openapi/strfmt"
)

type Date struct {
	strfmt.Date
}

func (m *Date) Time() time.Time {
	return time.Time(m.Date)
}

func (m *Date) Validate(strfmt.Registry) error {
	return nil
}

func (m *Date) UnmarshalJSON(b []byte) error {
	return m.Date.UnmarshalJSON(b)
}

func (m *Date) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}
