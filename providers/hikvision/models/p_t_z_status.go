// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PTZStatus p t z status
//
// swagger:model PTZStatus
type PTZStatus struct {

	// absolute high
	AbsoluteHigh *PtzAbsoluteHigh `json:"absoluteHigh,omitempty" xml:"AbsoluteHigh,omitempty"`
}

// Validate validates this p t z status
func (m *PTZStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAbsoluteHigh(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PTZStatus) validateAbsoluteHigh(formats strfmt.Registry) error {
	if swag.IsZero(m.AbsoluteHigh) { // not required
		return nil
	}

	if m.AbsoluteHigh != nil {
		if err := m.AbsoluteHigh.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("absoluteHigh")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("absoluteHigh")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this p t z status based on the context it is used
func (m *PTZStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAbsoluteHigh(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PTZStatus) contextValidateAbsoluteHigh(ctx context.Context, formats strfmt.Registry) error {

	if m.AbsoluteHigh != nil {

		if swag.IsZero(m.AbsoluteHigh) { // not required
			return nil
		}

		if err := m.AbsoluteHigh.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("absoluteHigh")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("absoluteHigh")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PTZStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PTZStatus) UnmarshalBinary(b []byte) error {
	var res PTZStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
