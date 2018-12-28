// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// OptIn opt in
// swagger:model OptIn
type OptIn struct {

	// crash
	Crash bool `json:"crash,omitempty"`

	// opencast
	Opencast bool `json:"opencast,omitempty"`

	// stats
	Stats bool `json:"stats,omitempty"`
}

// Validate validates this opt in
func (m *OptIn) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *OptIn) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OptIn) UnmarshalBinary(b []byte) error {
	var res OptIn
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
