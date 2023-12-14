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

// MobileAppSettings mobile app settings
//
// swagger:model MobileAppSettings
type MobileAppSettings struct {

	// dont show debt
	DontShowDebt bool `json:"DontShowDebt,omitempty"`

	// mockup count
	MockupCount int64 `json:"MockupCount,omitempty"`

	// mockup max height
	MockupMaxHeight int64 `json:"MockupMaxHeight,omitempty"`

	// mockup max width
	MockupMaxWidth int64 `json:"MockupMaxWidth,omitempty"`

	// address
	Address string `json:"address,omitempty"`

	// ads code android
	AdsCodeAndroid string `json:"adsCodeAndroid,omitempty"`

	// ads code i o s
	AdsCodeIOS string `json:"adsCodeIOS,omitempty"`

	// ads type
	AdsType int64 `json:"adsType,omitempty"`

	// app icon file
	AppIconFile string `json:"appIconFile,omitempty"`

	// app link android
	AppLinkAndroid string `json:"appLinkAndroid,omitempty"`

	// app link i o s
	AppLinkIOS string `json:"appLinkIOS,omitempty"`

	// app theme
	AppTheme string `json:"appTheme,omitempty"`

	// block user auth
	BlockUserAuth bool `json:"blockUserAuth,omitempty"`

	// bonus oferta file
	BonusOfertaFile string `json:"bonusOfertaFile,omitempty"`

	// choose ident by house
	ChooseIdentByHouse bool `json:"chooseIdentByHouse,omitempty"`

	// color
	Color string `json:"color,omitempty"`

	// disable commenting requests
	DisableCommentingRequests bool `json:"disableCommentingRequests,omitempty"`

	// districts exists
	DistrictsExists bool `json:"districtsExists,omitempty"`

	// enable o s s
	EnableOSS bool `json:"enableOSS,omitempty"`

	// houses exists
	HousesExists bool `json:"housesExists,omitempty"`

	// language
	Language string `json:"language,omitempty"`

	// main name
	MainName string `json:"main_name,omitempty"`

	// menu
	Menu []*MobileAppSettingsMenuItems0 `json:"menu"`

	// phone
	Phone string `json:"phone,omitempty"`

	// premises exists
	PremisesExists bool `json:"premisesExists,omitempty"`

	// register without s m s
	RegisterWithoutSMS bool `json:"registerWithoutSMS,omitempty"`

	// require birth date
	RequireBirthDate bool `json:"requireBirthDate,omitempty"`

	// service percent
	ServicePercent float64 `json:"servicePercent,omitempty"`

	// show ads
	ShowAds bool `json:"showAds,omitempty"`

	// show our service
	ShowOurService bool `json:"showOurService,omitempty"`

	// site icon file
	SiteIconFile string `json:"siteIconFile,omitempty"`

	// start screen
	StartScreen string `json:"startScreen,omitempty"`

	// streets exists
	StreetsExists bool `json:"streetsExists,omitempty"`

	// use account pin code
	UseAccountPinCode bool `json:"useAccountPinCode,omitempty"`

	// use bonus system
	UseBonusSystem bool `json:"useBonusSystem,omitempty"`

	// use dispatcher auth
	UseDispatcherAuth bool `json:"useDispatcherAuth,omitempty"`

	// сheck crash system
	СheckCrashSystem bool `json:"сheckCrashSystem,omitempty"`
}

// Validate validates this mobile app settings
func (m *MobileAppSettings) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMenu(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MobileAppSettings) validateMenu(formats strfmt.Registry) error {
	if swag.IsZero(m.Menu) { // not required
		return nil
	}

	for i := 0; i < len(m.Menu); i++ {
		if swag.IsZero(m.Menu[i]) { // not required
			continue
		}

		if m.Menu[i] != nil {
			if err := m.Menu[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("menu" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("menu" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this mobile app settings based on the context it is used
func (m *MobileAppSettings) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMenu(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *MobileAppSettings) contextValidateMenu(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Menu); i++ {

		if m.Menu[i] != nil {
			if err := m.Menu[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("menu" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("menu" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *MobileAppSettings) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MobileAppSettings) UnmarshalBinary(b []byte) error {
	var res MobileAppSettings
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// MobileAppSettingsMenuItems0 mobile app settings menu items0
//
// swagger:model MobileAppSettingsMenuItems0
type MobileAppSettingsMenuItems0 struct {

	// id
	ID int64 `json:"id,omitempty"`

	// name app
	NameApp string `json:"name_app,omitempty"`

	// simple name
	SimpleName string `json:"simple_name,omitempty"`

	// visible
	Visible int64 `json:"visible,omitempty"`
}

// Validate validates this mobile app settings menu items0
func (m *MobileAppSettingsMenuItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this mobile app settings menu items0 based on context it is used
func (m *MobileAppSettingsMenuItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MobileAppSettingsMenuItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MobileAppSettingsMenuItems0) UnmarshalBinary(b []byte) error {
	var res MobileAppSettingsMenuItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
