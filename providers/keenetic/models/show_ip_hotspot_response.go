// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ShowIPHotspotResponse show IP hotspot response
//
// swagger:model ShowIPHotspotResponse
type ShowIPHotspotResponse struct {

	// host
	Host []*ShowIPHotspotResponseHostItems0 `json:"host"`
}

// Validate validates this show IP hotspot response
func (m *ShowIPHotspotResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHost(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ShowIPHotspotResponse) validateHost(formats strfmt.Registry) error {
	if swag.IsZero(m.Host) { // not required
		return nil
	}

	for i := 0; i < len(m.Host); i++ {
		if swag.IsZero(m.Host[i]) { // not required
			continue
		}

		if m.Host[i] != nil {
			if err := m.Host[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("host" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("host" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this show IP hotspot response based on the context it is used
func (m *ShowIPHotspotResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateHost(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ShowIPHotspotResponse) contextValidateHost(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Host); i++ {

		if m.Host[i] != nil {
			if err := m.Host[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("host" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("host" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ShowIPHotspotResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ShowIPHotspotResponse) UnmarshalBinary(b []byte) error {
	var res ShowIPHotspotResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ShowIPHotspotResponseHostItems0 show IP hotspot response host items0
//
// swagger:model ShowIPHotspotResponseHostItems0
type ShowIPHotspotResponseHostItems0 struct {

	// 11
	Nr11 []string `json:"_11"`

	// access
	Access string `json:"access,omitempty"`

	// active
	Active bool `json:"active,omitempty"`

	// ap
	Ap string `json:"ap,omitempty"`

	// authenticated
	Authenticated bool `json:"authenticated,omitempty"`

	// dhcp
	Dhcp *ShowIPHotspotResponseHostItems0Dhcp `json:"dhcp,omitempty"`

	// dl mu
	DlMu bool `json:"dl-mu,omitempty"`

	// ebf
	Ebf bool `json:"ebf,omitempty"`

	// first seen
	FirstSeen int64 `json:"first-seen,omitempty"`

	// gi
	Gi int64 `json:"gi,omitempty"`

	// hostname
	Hostname string `json:"hostname,omitempty"`

	// ht
	Ht int64 `json:"ht,omitempty"`

	// interface
	Interface *ShowIPHotspotResponseHostItems0Interface `json:"interface,omitempty"`

	// ip
	IP string `json:"ip,omitempty"`

	// last seen
	LastSeen int64 `json:"last-seen,omitempty"`

	// link
	Link string `json:"link,omitempty"`

	// mac
	Mac string `json:"mac,omitempty"`

	// mcs
	Mcs int64 `json:"mcs,omitempty"`

	// mode
	Mode string `json:"mode,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// priority
	Priority int64 `json:"priority,omitempty"`

	// registered
	Registered bool `json:"registered,omitempty"`

	// rssi
	Rssi int64 `json:"rssi,omitempty"`

	// rxbytes
	Rxbytes int64 `json:"rxbytes,omitempty"`

	// schedule
	Schedule string `json:"schedule,omitempty"`

	// security
	Security string `json:"security,omitempty"`

	// ssid
	Ssid string `json:"ssid,omitempty"`

	// traffic shape
	TrafficShape *ShowIPHotspotResponseHostItems0TrafficShape `json:"traffic-shape,omitempty"`

	// txbytes
	Txbytes int64 `json:"txbytes,omitempty"`

	// txrate
	Txrate int64 `json:"txrate,omitempty"`

	// txss
	Txss int64 `json:"txss,omitempty"`

	// uptime
	Uptime int64 `json:"uptime,omitempty"`

	// via
	Via string `json:"via,omitempty"`
}

// Validate validates this show IP hotspot response host items0
func (m *ShowIPHotspotResponseHostItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDhcp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInterface(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTrafficShape(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ShowIPHotspotResponseHostItems0) validateDhcp(formats strfmt.Registry) error {
	if swag.IsZero(m.Dhcp) { // not required
		return nil
	}

	if m.Dhcp != nil {
		if err := m.Dhcp.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("dhcp")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("dhcp")
			}
			return err
		}
	}

	return nil
}

func (m *ShowIPHotspotResponseHostItems0) validateInterface(formats strfmt.Registry) error {
	if swag.IsZero(m.Interface) { // not required
		return nil
	}

	if m.Interface != nil {
		if err := m.Interface.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("interface")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("interface")
			}
			return err
		}
	}

	return nil
}

func (m *ShowIPHotspotResponseHostItems0) validateTrafficShape(formats strfmt.Registry) error {
	if swag.IsZero(m.TrafficShape) { // not required
		return nil
	}

	if m.TrafficShape != nil {
		if err := m.TrafficShape.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("traffic-shape")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("traffic-shape")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this show IP hotspot response host items0 based on the context it is used
func (m *ShowIPHotspotResponseHostItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDhcp(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateInterface(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTrafficShape(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ShowIPHotspotResponseHostItems0) contextValidateDhcp(ctx context.Context, formats strfmt.Registry) error {

	if m.Dhcp != nil {
		if err := m.Dhcp.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("dhcp")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("dhcp")
			}
			return err
		}
	}

	return nil
}

func (m *ShowIPHotspotResponseHostItems0) contextValidateInterface(ctx context.Context, formats strfmt.Registry) error {

	if m.Interface != nil {
		if err := m.Interface.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("interface")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("interface")
			}
			return err
		}
	}

	return nil
}

func (m *ShowIPHotspotResponseHostItems0) contextValidateTrafficShape(ctx context.Context, formats strfmt.Registry) error {

	if m.TrafficShape != nil {
		if err := m.TrafficShape.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("traffic-shape")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("traffic-shape")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ShowIPHotspotResponseHostItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ShowIPHotspotResponseHostItems0) UnmarshalBinary(b []byte) error {
	var res ShowIPHotspotResponseHostItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ShowIPHotspotResponseHostItems0Dhcp show IP hotspot response host items0 dhcp
//
// swagger:model ShowIPHotspotResponseHostItems0Dhcp
type ShowIPHotspotResponseHostItems0Dhcp struct {

	// static
	Static bool `json:"static,omitempty"`
}

// Validate validates this show IP hotspot response host items0 dhcp
func (m *ShowIPHotspotResponseHostItems0Dhcp) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this show IP hotspot response host items0 dhcp based on context it is used
func (m *ShowIPHotspotResponseHostItems0Dhcp) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ShowIPHotspotResponseHostItems0Dhcp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ShowIPHotspotResponseHostItems0Dhcp) UnmarshalBinary(b []byte) error {
	var res ShowIPHotspotResponseHostItems0Dhcp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ShowIPHotspotResponseHostItems0Interface show IP hotspot response host items0 interface
//
// swagger:model ShowIPHotspotResponseHostItems0Interface
type ShowIPHotspotResponseHostItems0Interface struct {

	// description
	Description string `json:"description,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this show IP hotspot response host items0 interface
func (m *ShowIPHotspotResponseHostItems0Interface) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this show IP hotspot response host items0 interface based on context it is used
func (m *ShowIPHotspotResponseHostItems0Interface) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ShowIPHotspotResponseHostItems0Interface) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ShowIPHotspotResponseHostItems0Interface) UnmarshalBinary(b []byte) error {
	var res ShowIPHotspotResponseHostItems0Interface
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ShowIPHotspotResponseHostItems0TrafficShape show IP hotspot response host items0 traffic shape
//
// swagger:model ShowIPHotspotResponseHostItems0TrafficShape
type ShowIPHotspotResponseHostItems0TrafficShape struct {

	// mode
	Mode string `json:"mode,omitempty"`

	// rx
	Rx int64 `json:"rx,omitempty"`

	// schedule
	Schedule string `json:"schedule,omitempty"`

	// tx
	Tx int64 `json:"tx,omitempty"`
}

// Validate validates this show IP hotspot response host items0 traffic shape
func (m *ShowIPHotspotResponseHostItems0TrafficShape) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this show IP hotspot response host items0 traffic shape based on context it is used
func (m *ShowIPHotspotResponseHostItems0TrafficShape) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ShowIPHotspotResponseHostItems0TrafficShape) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ShowIPHotspotResponseHostItems0TrafficShape) UnmarshalBinary(b []byte) error {
	var res ShowIPHotspotResponseHostItems0TrafficShape
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
