// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ShowDefaultsResponse show defaults response
//
// swagger:model ShowDefaultsResponse
type ShowDefaultsResponse struct {

	// country
	Country string `json:"country,omitempty"`

	// ctrlsum
	Ctrlsum string `json:"ctrlsum,omitempty"`

	// integrity
	Integrity string `json:"integrity,omitempty"`

	// locked
	Locked bool `json:"locked,omitempty"`

	// ndmhwid
	Ndmhwid string `json:"ndmhwid,omitempty"`

	// product
	Product string `json:"product,omitempty"`

	// serial
	Serial string `json:"serial,omitempty"`

	// servicehost
	Servicehost string `json:"servicehost,omitempty"`

	// servicepass
	Servicepass string `json:"servicepass,omitempty"`

	// servicetag
	Servicetag string `json:"servicetag,omitempty"`

	// signature
	Signature string `json:"signature,omitempty"`

	// wlankey
	Wlankey string `json:"wlankey,omitempty"`

	// wlanssid
	Wlanssid string `json:"wlanssid,omitempty"`

	// wlanwps
	Wlanwps string `json:"wlanwps,omitempty"`
}

// Validate validates this show defaults response
func (m *ShowDefaultsResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this show defaults response based on context it is used
func (m *ShowDefaultsResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ShowDefaultsResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ShowDefaultsResponse) UnmarshalBinary(b []byte) error {
	var res ShowDefaultsResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
