// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TemperatureData temperature data
//
// swagger:model TemperatureData
type TemperatureData struct {

	// actual
	Actual float64 `json:"actual,omitempty"`

	// offset
	Offset int64 `json:"offset,omitempty"`

	// target
	Target float64 `json:"target,omitempty"`
}

// Validate validates this temperature data
func (m *TemperatureData) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this temperature data based on context it is used
func (m *TemperatureData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TemperatureData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TemperatureData) UnmarshalBinary(b []byte) error {
	var res TemperatureData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
