// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// ProfileTypeDTO profile type d t o
// swagger:model ProfileTypeDTO
type ProfileTypeDTO struct {

	// kind
	Kind string `json:"kind,omitempty"`

	// label
	Label string `json:"label,omitempty"`

	// supported item types
	SupportedItemTypes []string `json:"supportedItemTypes"`

	// uid
	UID string `json:"uid,omitempty"`
}

// Validate validates this profile type d t o
func (m *ProfileTypeDTO) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ProfileTypeDTO) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProfileTypeDTO) UnmarshalBinary(b []byte) error {
	var res ProfileTypeDTO
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
