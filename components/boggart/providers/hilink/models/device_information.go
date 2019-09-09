// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// DeviceInformation device information
// swagger:model DeviceInformation
type DeviceInformation struct {

	// classify
	Classify string `json:"Classify,omitempty" xml:"Classify"`

	// device name
	DeviceName string `json:"DeviceName,omitempty" xml:"DeviceName"`

	// hardware version
	HardwareVersion string `json:"HardwareVersion,omitempty" xml:"HardwareVersion"`

	// i c c ID
	ICCID string `json:"ICCID,omitempty" xml:"Iccid"`

	// i m e i
	IMEI string `json:"IMEI,omitempty" xml:"Imei"`

	// i m s i
	IMSI string `json:"IMSI,omitempty" xml:"Imsi"`

	// m s i s d n
	MSISDN string `json:"MSISDN,omitempty" xml:"Msisdn"`

	// mac address1
	MacAddress1 string `json:"MacAddress1,omitempty" xml:"MacAddress1"`

	// mac address2
	MacAddress2 string `json:"MacAddress2,omitempty" xml:"MacAddress2"`

	// product family
	ProductFamily string `json:"ProductFamily,omitempty" xml:"ProductFamily"`

	// serial number
	SerialNumber string `json:"SerialNumber,omitempty" xml:"SerialNumber"`

	// software version
	SoftwareVersion string `json:"SoftwareVersion,omitempty" xml:"SoftwareVersion"`

	// support mode
	SupportMode string `json:"SupportMode,omitempty" xml:"supportmode"`

	// wan IP address
	WanIPAddress string `json:"WanIPAddress,omitempty" xml:"WanIPAddress"`

	// wan IPv6 address
	WanIPV6Address string `json:"WanIPv6Address,omitempty" xml:"WanIPv6Address"`

	// web UI version
	WebUIVersion string `json:"WebUIVersion,omitempty" xml:"WebUIVersion"`

	// work mode
	WorkMode string `json:"WorkMode,omitempty" xml:"workmode"`
}

// Validate validates this device information
func (m *DeviceInformation) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DeviceInformation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DeviceInformation) UnmarshalBinary(b []byte) error {
	var res DeviceInformation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
