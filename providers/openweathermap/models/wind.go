// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Wind wind
//
// swagger:model Wind
type Wind struct {

	// deg
	Deg uint64 `json:"deg,omitempty"`

	// gust
	Gust float64 `json:"gust,omitempty"`

	// speed
	Speed float64 `json:"speed,omitempty"`
}

// Validate validates this wind
func (m *Wind) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this wind based on context it is used
func (m *Wind) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Wind) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Wind) UnmarshalBinary(b []byte) error {
	var res Wind
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
