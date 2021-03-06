// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Locale locale
// swagger:model Locale
type Locale struct {

	// display string
	// Required: true
	DisplayString *string `json:"display_string"`

	// locale
	Locale string `json:"locale,omitempty"`
}

// Validate validates this locale
func (m *Locale) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDisplayString(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Locale) validateDisplayString(formats strfmt.Registry) error {

	if err := validate.Required("display_string", "body", m.DisplayString); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Locale) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Locale) UnmarshalBinary(b []byte) error {
	var res Locale
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
