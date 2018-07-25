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

// ChannelDTO channel d t o
// swagger:model ChannelDTO
type ChannelDTO struct {

	// channel type UID
	ChannelTypeUID string `json:"channelTypeUID,omitempty"`

	// configuration
	Configuration map[string]interface{} `json:"configuration,omitempty"`

	// default tags
	// Unique: true
	DefaultTags []string `json:"defaultTags"`

	// description
	Description string `json:"description,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// item type
	ItemType string `json:"itemType,omitempty"`

	// kind
	Kind string `json:"kind,omitempty"`

	// label
	Label string `json:"label,omitempty"`

	// properties
	Properties map[string]string `json:"properties,omitempty"`

	// uid
	UID string `json:"uid,omitempty"`
}

// Validate validates this channel d t o
func (m *ChannelDTO) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDefaultTags(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ChannelDTO) validateDefaultTags(formats strfmt.Registry) error {

	if swag.IsZero(m.DefaultTags) { // not required
		return nil
	}

	if err := validate.UniqueItems("defaultTags", "body", m.DefaultTags); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ChannelDTO) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ChannelDTO) UnmarshalBinary(b []byte) error {
	var res ChannelDTO
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}