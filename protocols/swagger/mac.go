package swagger

import (
	"context"
	"net"

	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/performance"
)

type MAC net.HardwareAddr

func (m *MAC) Validate(strfmt.Registry) error {
	return nil
}

func (m *MAC) UnmarshalJSON(b []byte) error {
	decoded, err := net.ParseMAC(replacerAsStringCleanValue.Replace(performance.UnsafeBytes2String(b)))
	if err != nil {
		return err
	}

	*m = MAC(decoded)
	return nil
}

func (m *MAC) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}

func (m *MAC) Value() net.HardwareAddr {
	return net.HardwareAddr(*m)
}
