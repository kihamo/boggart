// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// EurekaInfo eureka info
// swagger:model EurekaInfo
type EurekaInfo struct {

	// audio
	Audio *Audio `json:"audio,omitempty"`

	// bssid
	Bssid string `json:"bssid,omitempty"`

	// build info
	BuildInfo *BuildInfo `json:"build_info,omitempty"`

	// build version
	BuildVersion string `json:"build_version,omitempty"`

	// cast build revision
	CastBuildRevision string `json:"cast_build_revision,omitempty"`

	// closed caption
	ClosedCaption interface{} `json:"closed_caption,omitempty"`

	// connected
	Connected bool `json:"connected,omitempty"`

	// detail
	Detail *Detail `json:"detail,omitempty"`

	// device info
	DeviceInfo *DeviceInfo `json:"device_info,omitempty"`

	// ethernet connected
	EthernetConnected bool `json:"ethernet_connected,omitempty"`

	// has update
	HasUpdate bool `json:"has_update,omitempty"`

	// hotspot bssid
	HotspotBssid string `json:"hotspot_bssid,omitempty"`

	// ip address
	IPAddress string `json:"ip_address,omitempty"`

	// locale
	Locale string `json:"locale,omitempty"`

	// location
	Location *Location `json:"location,omitempty"`

	// mac address
	MacAddress string `json:"mac_address,omitempty"`

	// multizone
	Multizone Multizone `json:"multizone,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// net
	Net *Net `json:"net,omitempty"`

	// night mode params
	NightModeParams NightModeParams `json:"night_mode_params,omitempty"`

	// noise level
	NoiseLevel int64 `json:"noise_level,omitempty"`

	// opt in
	OptIn *OptIn `json:"opt_in,omitempty"`

	// proxy
	Proxy *Proxy `json:"proxy,omitempty"`

	// public key
	PublicKey string `json:"public_key,omitempty"`

	// release track
	ReleaseTrack string `json:"release_track,omitempty"`

	// settings
	Settings *Settings `json:"settings,omitempty"`

	// setup
	Setup *Setup `json:"setup,omitempty"`

	// setup state
	SetupState int64 `json:"setup_state,omitempty"`

	// setup stats
	SetupStats *Stats `json:"setup_stats,omitempty"`

	// signal level
	SignalLevel int64 `json:"signal_level,omitempty"`

	// ssdp udn
	SsdpUdn string `json:"ssdp_udn,omitempty"`

	// ssid
	Ssid string `json:"ssid,omitempty"`

	// time format
	TimeFormat int64 `json:"time_format,omitempty"`

	// timezone
	Timezone string `json:"timezone,omitempty"`

	// tos accepted
	TosAccepted bool `json:"tos_accepted,omitempty"`

	// uptime
	Uptime float32 `json:"uptime,omitempty"`

	// user eq
	UserEq UserEq `json:"user_eq,omitempty"`

	// version
	Version int64 `json:"version,omitempty"`

	// wpa configured
	WpaConfigured bool `json:"wpa_configured,omitempty"`

	// wpa id
	WpaID int64 `json:"wpa_id,omitempty"`

	// wpa state
	WpaState int64 `json:"wpa_state,omitempty"`
}

// Validate validates this eureka info
func (m *EurekaInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAudio(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBuildInfo(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDetail(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDeviceInfo(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLocation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNet(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOptIn(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProxy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSettings(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSetup(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSetupStats(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EurekaInfo) validateAudio(formats strfmt.Registry) error {

	if swag.IsZero(m.Audio) { // not required
		return nil
	}

	if m.Audio != nil {
		if err := m.Audio.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("audio")
			}
			return err
		}
	}

	return nil
}

func (m *EurekaInfo) validateBuildInfo(formats strfmt.Registry) error {

	if swag.IsZero(m.BuildInfo) { // not required
		return nil
	}

	if m.BuildInfo != nil {
		if err := m.BuildInfo.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("build_info")
			}
			return err
		}
	}

	return nil
}

func (m *EurekaInfo) validateDetail(formats strfmt.Registry) error {

	if swag.IsZero(m.Detail) { // not required
		return nil
	}

	if m.Detail != nil {
		if err := m.Detail.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("detail")
			}
			return err
		}
	}

	return nil
}

func (m *EurekaInfo) validateDeviceInfo(formats strfmt.Registry) error {

	if swag.IsZero(m.DeviceInfo) { // not required
		return nil
	}

	if m.DeviceInfo != nil {
		if err := m.DeviceInfo.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("device_info")
			}
			return err
		}
	}

	return nil
}

func (m *EurekaInfo) validateLocation(formats strfmt.Registry) error {

	if swag.IsZero(m.Location) { // not required
		return nil
	}

	if m.Location != nil {
		if err := m.Location.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("location")
			}
			return err
		}
	}

	return nil
}

func (m *EurekaInfo) validateNet(formats strfmt.Registry) error {

	if swag.IsZero(m.Net) { // not required
		return nil
	}

	if m.Net != nil {
		if err := m.Net.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("net")
			}
			return err
		}
	}

	return nil
}

func (m *EurekaInfo) validateOptIn(formats strfmt.Registry) error {

	if swag.IsZero(m.OptIn) { // not required
		return nil
	}

	if m.OptIn != nil {
		if err := m.OptIn.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("opt_in")
			}
			return err
		}
	}

	return nil
}

func (m *EurekaInfo) validateProxy(formats strfmt.Registry) error {

	if swag.IsZero(m.Proxy) { // not required
		return nil
	}

	if m.Proxy != nil {
		if err := m.Proxy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("proxy")
			}
			return err
		}
	}

	return nil
}

func (m *EurekaInfo) validateSettings(formats strfmt.Registry) error {

	if swag.IsZero(m.Settings) { // not required
		return nil
	}

	if m.Settings != nil {
		if err := m.Settings.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("settings")
			}
			return err
		}
	}

	return nil
}

func (m *EurekaInfo) validateSetup(formats strfmt.Registry) error {

	if swag.IsZero(m.Setup) { // not required
		return nil
	}

	if m.Setup != nil {
		if err := m.Setup.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("setup")
			}
			return err
		}
	}

	return nil
}

func (m *EurekaInfo) validateSetupStats(formats strfmt.Registry) error {

	if swag.IsZero(m.SetupStats) { // not required
		return nil
	}

	if m.SetupStats != nil {
		if err := m.SetupStats.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("setup_stats")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *EurekaInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EurekaInfo) UnmarshalBinary(b []byte) error {
	var res EurekaInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
