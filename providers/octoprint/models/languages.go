// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Languages languages
//
// swagger:model Languages
type Languages struct {

	// language packs
	LanguagePacks map[string]LanguagesLanguagePacksAnon `json:"language_packs,omitempty"`
}

// Validate validates this languages
func (m *Languages) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLanguagePacks(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Languages) validateLanguagePacks(formats strfmt.Registry) error {

	if swag.IsZero(m.LanguagePacks) { // not required
		return nil
	}

	for k := range m.LanguagePacks {

		if swag.IsZero(m.LanguagePacks[k]) { // not required
			continue
		}
		if val, ok := m.LanguagePacks[k]; ok {
			if err := val.Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Languages) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Languages) UnmarshalBinary(b []byte) error {
	var res Languages
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// LanguagesLanguagePacksAnon languages language packs anon
//
// swagger:model LanguagesLanguagePacksAnon
type LanguagesLanguagePacksAnon struct {

	// display
	Display string `json:"display,omitempty"`

	// identifier
	Identifier string `json:"identifier,omitempty"`

	// languages
	Languages []*LanguagesLanguagePacksAnonLanguagesItems0 `json:"languages"`
}

// Validate validates this languages language packs anon
func (m *LanguagesLanguagePacksAnon) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLanguages(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *LanguagesLanguagePacksAnon) validateLanguages(formats strfmt.Registry) error {

	if swag.IsZero(m.Languages) { // not required
		return nil
	}

	for i := 0; i < len(m.Languages); i++ {
		if swag.IsZero(m.Languages[i]) { // not required
			continue
		}

		if m.Languages[i] != nil {
			if err := m.Languages[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("languages" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *LanguagesLanguagePacksAnon) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LanguagesLanguagePacksAnon) UnmarshalBinary(b []byte) error {
	var res LanguagesLanguagePacksAnon
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// LanguagesLanguagePacksAnonLanguagesItems0 languages language packs anon languages items0
//
// swagger:model LanguagesLanguagePacksAnonLanguagesItems0
type LanguagesLanguagePacksAnonLanguagesItems0 struct {

	// author
	Author string `json:"author,omitempty"`

	// last update
	LastUpdate float64 `json:"last_update,omitempty"`

	// locale
	Locale string `json:"locale,omitempty"`

	// locale display
	LocaleDisplay string `json:"locale_display,omitempty"`

	// locale english
	LocaleEnglish string `json:"locale_english,omitempty"`
}

// Validate validates this languages language packs anon languages items0
func (m *LanguagesLanguagePacksAnonLanguagesItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *LanguagesLanguagePacksAnonLanguagesItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LanguagesLanguagePacksAnonLanguagesItems0) UnmarshalBinary(b []byte) error {
	var res LanguagesLanguagePacksAnonLanguagesItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
