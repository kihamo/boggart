// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"

	custom "github.com/kihamo/boggart/providers/dom24/static/models"
)

// Payment payment
// swagger:model Payment
type Payment struct {

	// date
	// Format: date
	Date custom.Date `json:"Date,omitempty"`

	// ident
	Ident string `json:"Ident,omitempty"`

	// period
	Period string `json:"Period,omitempty"`

	// sum
	Sum float64 `json:"Sum,omitempty"`
}

// Validate validates this payment
func (m *Payment) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Payment) validateDate(formats strfmt.Registry) error {

	if swag.IsZero(m.Date) { // not required
		return nil
	}

	if err := m.Date.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("Date")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Payment) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Payment) UnmarshalBinary(b []byte) error {
	var res Payment
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
