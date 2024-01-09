// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SystemNtpServers system ntp servers
//
// swagger:model SystemNtpServers
type SystemNtpServers []*SystemNtpServersItems0

// Validate validates this system ntp servers
func (m SystemNtpServers) Validate(formats strfmt.Registry) error {
	var res []error

	for i := 0; i < len(m); i++ {
		if swag.IsZero(m[i]) { // not required
			continue
		}

		if m[i] != nil {
			if err := m[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName(strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName(strconv.Itoa(i))
				}
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validate this system ntp servers based on the context it is used
func (m SystemNtpServers) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	for i := 0; i < len(m); i++ {

		if m[i] != nil {

			if swag.IsZero(m[i]) { // not required
				return nil
			}

			if err := m[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName(strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName(strconv.Itoa(i))
				}
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// SystemNtpServersItems0 system ntp servers items0
//
// swagger:model SystemNtpServersItems0
type SystemNtpServersItems0 struct {

	// addressing format type
	// Enum: [ipaddress hostname]
	AddressingFormatType string `json:"addressingFormatType,omitempty" xml:"NTPServer>addressingFormatType,omitempty"`

	// host name
	HostName *string `json:"hostName,omitempty" xml:"NTPServer>hostName,omitempty"`

	// id
	ID uint64 `json:"id,omitempty" xml:"NTPServer>id,omitempty"`

	// ip address
	IPAddress *string `json:"ipAddress,omitempty" xml:"NTPServer>ipAddress,omitempty"`

	// ipv6 address
	IPV6Address *string `json:"ipv6Address,omitempty" xml:"NTPServer>ipv6Address,omitempty"`

	// port no
	PortNo uint64 `json:"portNo,omitempty" xml:"NTPServer>portNo,omitempty"`

	// synchronize interval
	SynchronizeInterval uint64 `json:"synchronizeInterval,omitempty" xml:"NTPServer>synchronizeInterval,omitempty"`
}

// Validate validates this system ntp servers items0
func (m *SystemNtpServersItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAddressingFormatType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var systemNtpServersItems0TypeAddressingFormatTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["ipaddress","hostname"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		systemNtpServersItems0TypeAddressingFormatTypePropEnum = append(systemNtpServersItems0TypeAddressingFormatTypePropEnum, v)
	}
}

const (

	// SystemNtpServersItems0AddressingFormatTypeIpaddress captures enum value "ipaddress"
	SystemNtpServersItems0AddressingFormatTypeIpaddress string = "ipaddress"

	// SystemNtpServersItems0AddressingFormatTypeHostname captures enum value "hostname"
	SystemNtpServersItems0AddressingFormatTypeHostname string = "hostname"
)

// prop value enum
func (m *SystemNtpServersItems0) validateAddressingFormatTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, systemNtpServersItems0TypeAddressingFormatTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *SystemNtpServersItems0) validateAddressingFormatType(formats strfmt.Registry) error {
	if swag.IsZero(m.AddressingFormatType) { // not required
		return nil
	}

	// value enum
	if err := m.validateAddressingFormatTypeEnum("addressingFormatType", "body", m.AddressingFormatType); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this system ntp servers items0 based on context it is used
func (m *SystemNtpServersItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SystemNtpServersItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SystemNtpServersItems0) UnmarshalBinary(b []byte) error {
	var res SystemNtpServersItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
