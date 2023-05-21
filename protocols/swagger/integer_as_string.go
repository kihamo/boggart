package swagger

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/performance"
)

var replacerIntegerAsStringCleanValue = strings.NewReplacer(
	`"`, "",
)

type IntegerAsString int64

func (m *IntegerAsString) Validate(strfmt.Registry) error {
	return nil
}

func (m *IntegerAsString) UnmarshalJSON(b []byte) error {
	val, err := strconv.ParseInt(replacerIntegerAsStringCleanValue.Replace(performance.UnsafeBytes2String(b)), 10, 64)
	if err != nil {
		return err
	}

	*m = IntegerAsString(val)
	return nil
}

func (m *IntegerAsString) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}
