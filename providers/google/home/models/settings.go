// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// Settings settings
// swagger:model Settings
type Settings struct {

	// closed caption
	ClosedCaption interface{} `json:"closed_caption,omitempty"`

	// control notifications
	ControlNotifications int64 `json:"control_notifications,omitempty"`

	// country code
	CountryCode string `json:"country_code,omitempty"`

	// locale
	Locale string `json:"locale,omitempty"`

	// network standby
	NetworkStandby int64 `json:"network_standby,omitempty"`

	// system sound effects
	SystemSoundEffects bool `json:"system_sound_effects,omitempty"`

	// time format
	TimeFormat int64 `json:"time_format,omitempty"`

	// timezone
	Timezone string `json:"timezone,omitempty"`

	// wake on cast
	WakeOnCast int64 `json:"wake_on_cast,omitempty"`
}

// Validate validates this settings
func (m *Settings) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Settings) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Settings) UnmarshalBinary(b []byte) error {
	var res Settings
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
