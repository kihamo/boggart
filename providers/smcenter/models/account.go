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

// Account account
//
// swagger:model Account
type Account struct {

	// access o s s
	AccessOSS bool `json:"accessOSS,omitempty"`

	// accounts
	Accounts []*AccountAccountsItems0 `json:"accounts"`

	// acx
	Acx string `json:"acx,omitempty"`

	// address
	Address string `json:"address,omitempty"`

	// birthday
	Birthday string `json:"birthday,omitempty"`

	// company phone
	CompanyPhone string `json:"companyPhone,omitempty"`

	// email
	Email string `json:"email,omitempty"`

	// fio
	Fio string `json:"fio,omitempty"`

	// is dispatcher
	IsDispatcher bool `json:"isDispatcher,omitempty"`

	// login
	Login string `json:"login,omitempty"`

	// phone
	Phone string `json:"phone,omitempty"`

	// user settings
	UserSettings string `json:"userSettings,omitempty"`
}

// Validate validates this account
func (m *Account) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccounts(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Account) validateAccounts(formats strfmt.Registry) error {
	if swag.IsZero(m.Accounts) { // not required
		return nil
	}

	for i := 0; i < len(m.Accounts); i++ {
		if swag.IsZero(m.Accounts[i]) { // not required
			continue
		}

		if m.Accounts[i] != nil {
			if err := m.Accounts[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("accounts" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("accounts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this account based on the context it is used
func (m *Account) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAccounts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Account) contextValidateAccounts(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Accounts); i++ {

		if m.Accounts[i] != nil {

			if swag.IsZero(m.Accounts[i]) { // not required
				return nil
			}

			if err := m.Accounts[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("accounts" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("accounts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Account) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Account) UnmarshalBinary(b []byte) error {
	var res Account
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// AccountAccountsItems0 account accounts items0
//
// swagger:model AccountAccountsItems0
type AccountAccountsItems0 struct {

	// address
	Address string `json:"address,omitempty"`

	// cn
	Cn string `json:"cn,omitempty"`

	// company
	Company string `json:"company,omitempty"`

	// fio
	Fio string `json:"fio,omitempty"`

	// id
	ID uint64 `json:"id,omitempty"`

	// ident
	Ident string `json:"ident,omitempty"`

	// meters access flag
	MetersAccessFlag bool `json:"metersAccessFlag,omitempty"`

	// meters end day
	MetersEndDay uint64 `json:"metersEndDay,omitempty"`

	// meters period end is current
	MetersPeriodEndIsCurrent bool `json:"metersPeriodEndIsCurrent,omitempty"`

	// meters period start is current
	MetersPeriodStartIsCurrent bool `json:"metersPeriodStartIsCurrent,omitempty"`

	// meters start day
	MetersStartDay uint64 `json:"metersStartDay,omitempty"`

	// phone
	Phone string `json:"phone,omitempty"`
}

// Validate validates this account accounts items0
func (m *AccountAccountsItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this account accounts items0 based on context it is used
func (m *AccountAccountsItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AccountAccountsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AccountAccountsItems0) UnmarshalBinary(b []byte) error {
	var res AccountAccountsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
