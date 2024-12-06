package swagger

import (
	"context"
	"encoding/hex"

	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/performance"
)

type HEXString string

func (m *HEXString) Validate(strfmt.Registry) error {
	return nil
}

func (m *HEXString) UnmarshalJSON(b []byte) error {
	decoded, err := hex.DecodeString(replacerAsStringCleanValue.Replace(performance.UnsafeBytes2String(b)))
	if err != nil {
		return err
	}

	*m = HEXString(decoded)
	return nil
}

func (m *HEXString) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}

func (m *HEXString) Value() string {
	return string(*m)
}
