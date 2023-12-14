// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Device device
//
// swagger:model Device
type Device struct {

	// id
	ID int64 `json:"i,omitempty"`

	// name
	Name string `json:"n,omitempty"`

	// Settings
	Settings map[string]float64 `json:"st,omitempty"`

	// device severity level
	// Enum: [1 32 64]
	SeverityLevel int64 `json:"sev,omitempty"`

	// state
	State map[string]string `json:"s,omitempty"`

	// object type
	Type int64 `json:"t,omitempty"`

	// f
	F int64 `json:"f,omitempty"`
}

// Validate validates this device
func (m *Device) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSeverityLevel(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var deviceTypeSeverityLevelPropEnum []interface{}

func init() {
	var res []int64
	if err := json.Unmarshal([]byte(`[1,32,64]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		deviceTypeSeverityLevelPropEnum = append(deviceTypeSeverityLevelPropEnum, v)
	}
}

// prop value enum
func (m *Device) validateSeverityLevelEnum(path, location string, value int64) error {
	if err := validate.EnumCase(path, location, value, deviceTypeSeverityLevelPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Device) validateSeverityLevel(formats strfmt.Registry) error {
	if swag.IsZero(m.SeverityLevel) { // not required
		return nil
	}

	// value enum
	if err := m.validateSeverityLevelEnum("sev", "body", m.SeverityLevel); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this device based on context it is used
func (m *Device) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Device) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Device) UnmarshalBinary(b []byte) error {
	var res Device
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}