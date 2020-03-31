// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// ItemChannelLinkDTO item channel link d t o
// swagger:model ItemChannelLinkDTO
type ItemChannelLinkDTO struct {

	// channel UID
	ChannelUID string `json:"channelUID,omitempty"`

	// configuration
	Configuration map[string]interface{} `json:"configuration,omitempty"`

	// item name
	ItemName string `json:"itemName,omitempty"`
}

// Validate validates this item channel link d t o
func (m *ItemChannelLinkDTO) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ItemChannelLinkDTO) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ItemChannelLinkDTO) UnmarshalBinary(b []byte) error {
	var res ItemChannelLinkDTO
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}