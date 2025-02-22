// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	custom "github.com/kihamo/boggart/providers/smcenter/static/models"
)

// Bill bill
//
// swagger:model Bill
type Bill struct {

	// date
	// Format: date
	Date custom.Date `json:"Date,omitempty"`

	// file link
	FileLink string `json:"FileLink,omitempty"`

	// has file
	HasFile bool `json:"HasFile,omitempty"`

	// ID
	ID uint64 `json:"ID,omitempty"`

	// ident
	Ident string `json:"Ident,omitempty"`

	// period
	Period string `json:"Period,omitempty"`

	// total
	Total float64 `json:"Total,omitempty"`
}

// Validate validates this bill
func (m *Bill) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Bill) validateDate(formats strfmt.Registry) error {
	if swag.IsZero(m.Date) { // not required
		return nil
	}

	if err := m.Date.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("Date")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("Date")
		}
		return err
	}

	return nil
}

// ContextValidate validate this bill based on the context it is used
func (m *Bill) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Bill) contextValidateDate(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Date) { // not required
		return nil
	}

	if err := m.Date.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("Date")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("Date")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Bill) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Bill) UnmarshalBinary(b []byte) error {
	var res Bill
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
