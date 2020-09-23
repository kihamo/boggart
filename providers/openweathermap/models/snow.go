// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Snow snow
//
// swagger:model Snow
type Snow struct {

	// 1h
	Nr1h float64 `json:"1h,omitempty"`

	// 3h
	Nr3h float64 `json:"3h,omitempty"`
}

// Validate validates this snow
func (m *Snow) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Snow) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Snow) UnmarshalBinary(b []byte) error {
	var res Snow
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}