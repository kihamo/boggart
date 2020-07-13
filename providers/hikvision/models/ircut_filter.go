// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// IrcutFilter ircut filter
//
// swagger:model IrcutFilter
type IrcutFilter struct {

	// type
	// Enum: [auto day night shedule eventTrigger]
	Type string `json:"type,omitempty" xml:"IrcutFilterType,omitempty"`
}

// Validate validates this ircut filter
func (m *IrcutFilter) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var ircutFilterTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["auto","day","night","shedule","eventTrigger"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		ircutFilterTypeTypePropEnum = append(ircutFilterTypeTypePropEnum, v)
	}
}

const (

	// IrcutFilterTypeAuto captures enum value "auto"
	IrcutFilterTypeAuto string = "auto"

	// IrcutFilterTypeDay captures enum value "day"
	IrcutFilterTypeDay string = "day"

	// IrcutFilterTypeNight captures enum value "night"
	IrcutFilterTypeNight string = "night"

	// IrcutFilterTypeShedule captures enum value "shedule"
	IrcutFilterTypeShedule string = "shedule"

	// IrcutFilterTypeEventTrigger captures enum value "eventTrigger"
	IrcutFilterTypeEventTrigger string = "eventTrigger"
)

// prop value enum
func (m *IrcutFilter) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, ircutFilterTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *IrcutFilter) validateType(formats strfmt.Registry) error {

	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *IrcutFilter) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IrcutFilter) UnmarshalBinary(b []byte) error {
	var res IrcutFilter
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
