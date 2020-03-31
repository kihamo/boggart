// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// StateDescription state description
// swagger:model StateDescription
type StateDescription struct {

	// maximum
	Maximum float64 `json:"maximum,omitempty"`

	// minimum
	Minimum float64 `json:"minimum,omitempty"`

	// options
	Options []*StateOption `json:"options"`

	// pattern
	Pattern string `json:"pattern,omitempty"`

	// read only
	ReadOnly *bool `json:"readOnly,omitempty"`

	// step
	Step float64 `json:"step,omitempty"`
}

// Validate validates this state description
func (m *StateDescription) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateOptions(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StateDescription) validateOptions(formats strfmt.Registry) error {

	if swag.IsZero(m.Options) { // not required
		return nil
	}

	for i := 0; i < len(m.Options); i++ {
		if swag.IsZero(m.Options[i]) { // not required
			continue
		}

		if m.Options[i] != nil {
			if err := m.Options[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("options" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *StateDescription) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StateDescription) UnmarshalBinary(b []byte) error {
	var res StateDescription
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}