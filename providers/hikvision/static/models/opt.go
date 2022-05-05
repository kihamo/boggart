package models

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

type Opt struct {
}

func (m *Opt) Validate(formats strfmt.Registry) error {
	return nil
}

func (m *Opt) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	return swag.WriteJSON(m)
}

func (m *Opt) UnmarshalBinary(b []byte) error {
	var res Opt

	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}

	*m = res

	return nil
}

func (m *Opt) ContextValidate(context.Context, strfmt.Registry) error {
	return nil
}
