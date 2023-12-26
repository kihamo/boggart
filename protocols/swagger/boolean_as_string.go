package swagger

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/performance"
)

var replacerAsStringCleanValue = strings.NewReplacer(
	`"`, "",
)

type BooleanAsString bool

func (m *BooleanAsString) Validate(strfmt.Registry) error {
	return nil
}

func (m *BooleanAsString) UnmarshalJSON(b []byte) error {
	val, err := strconv.ParseBool(replacerAsStringCleanValue.Replace(performance.UnsafeBytes2String(b)))
	if err != nil {
		return err
	}

	*m = BooleanAsString(val)
	return nil
}

func (m *BooleanAsString) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}

func (m *BooleanAsString) Value() bool {
	return bool(*m)
}
