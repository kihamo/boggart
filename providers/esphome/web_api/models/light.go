// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Light light
// swagger:model Light
type Light struct {

	// brightness
	Brightness int64 `json:"brightness,omitempty"`

	// color
	Color *LightColor `json:"color,omitempty"`

	// color temp
	ColorTemp int64 `json:"color_temp,omitempty"`

	// effect
	Effect string `json:"effect,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// state
	State string `json:"state,omitempty"`

	// white value
	WhiteValue int64 `json:"white_value,omitempty"`
}

// Validate validates this light
func (m *Light) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateColor(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Light) validateColor(formats strfmt.Registry) error {

	if swag.IsZero(m.Color) { // not required
		return nil
	}

	if m.Color != nil {
		if err := m.Color.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("color")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Light) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Light) UnmarshalBinary(b []byte) error {
	var res Light
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// LightColor light color
// swagger:model LightColor
type LightColor struct {

	// b
	B int64 `json:"b,omitempty"`

	// g
	G int64 `json:"g,omitempty"`

	// r
	R int64 `json:"r,omitempty"`
}

// Validate validates this light color
func (m *LightColor) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *LightColor) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LightColor) UnmarshalBinary(b []byte) error {
	var res LightColor
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}