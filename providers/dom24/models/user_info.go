// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// UserInfo user info
// swagger:model UserInfo
type UserInfo struct {

	// accounts
	Accounts []*UserInfoAccountsItems0 `json:"Accounts"`

	// can pay
	CanPay bool `json:"CanPay,omitempty"`

	// email
	Email string `json:"Email,omitempty"`

	// f i o
	FIO string `json:"FIO,omitempty"`

	// login
	Login string `json:"Login,omitempty"`

	// phone
	Phone string `json:"Phone,omitempty"`
}

// Validate validates this user info
func (m *UserInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccounts(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *UserInfo) validateAccounts(formats strfmt.Registry) error {

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
					return ve.ValidateName("Accounts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *UserInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UserInfo) UnmarshalBinary(b []byte) error {
	var res UserInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// UserInfoAccountsItems0 user info accounts items0
// swagger:model UserInfoAccountsItems0
type UserInfoAccountsItems0 struct {

	// address
	Address string `json:"Address,omitempty"`

	// company
	Company string `json:"Company,omitempty"`

	// f i o
	FIO string `json:"FIO,omitempty"`

	// Id
	ID uint64 `json:"Id,omitempty"`

	// ident
	Ident string `json:"Ident,omitempty"`
}

// Validate validates this user info accounts items0
func (m *UserInfoAccountsItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *UserInfoAccountsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UserInfoAccountsItems0) UnmarshalBinary(b []byte) error {
	var res UserInfoAccountsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
