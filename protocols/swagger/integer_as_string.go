package swagger

import (
	"context"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/performance"
)

type IntegerAsString int64

func (m *IntegerAsString) Validate(strfmt.Registry) error {
	return nil
}

func (m *IntegerAsString) UnmarshalJSON(b []byte) error {
	val, err := strconv.ParseInt(replacerAsStringCleanValue.Replace(performance.UnsafeBytes2String(b)), 0, 64)
	if err != nil {
		return err
	}

	*m = IntegerAsString(val)
	return nil
}

func (m *IntegerAsString) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}

func (m *IntegerAsString) Value() int64 {
	return int64(*m)
}
