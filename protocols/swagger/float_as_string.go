package swagger

import (
	"context"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/performance"
)

type FloatAsString float64

func (m *FloatAsString) Validate(strfmt.Registry) error {
	return nil
}

func (m *FloatAsString) UnmarshalJSON(b []byte) error {
	val, err := strconv.ParseFloat(replacerAsStringCleanValue.Replace(performance.UnsafeBytes2String(b)), 64)
	if err != nil {
		return err
	}

	*m = FloatAsString(val)
	return nil
}

func (m *FloatAsString) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}

func (m *FloatAsString) Value() float64 {
	return float64(*m)
}
