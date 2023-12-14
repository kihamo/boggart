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

// AccountingInfo accounting info
//
// swagger:model AccountingInfo
type AccountingInfo struct {

	// account ID
	AccountID uint64 `json:"AccountID,omitempty"`

	// account type
	AccountType string `json:"AccountType,omitempty"`

	// address
	Address string `json:"Address,omitempty"`

	// bills
	Bills []*Bill `json:"Bills"`

	// bonus balance
	BonusBalance float64 `json:"BonusBalance,omitempty"`

	// comission
	Comission float64 `json:"Comission,omitempty"`

	// debt actual date
	DebtActualDate string `json:"DebtActualDate,omitempty"`

	// dont show insurance
	DontShowInsurance bool `json:"DontShowInsurance,omitempty"`

	// house Id
	HouseID uint64 `json:"HouseId,omitempty"`

	// i n n
	INN string `json:"INN,omitempty"`

	// ident
	Ident string `json:"Ident,omitempty"`

	// insurance sum
	InsuranceSum float64 `json:"InsuranceSum,omitempty"`

	// payments
	Payments []*Payment `json:"Payments"`

	// sum
	Sum float64 `json:"Sum,omitempty"`

	// sum fine
	SumFine float64 `json:"SumFine,omitempty"`

	// total sum
	TotalSum float64 `json:"TotalSum,omitempty"`
}

// Validate validates this accounting info
func (m *AccountingInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBills(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePayments(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AccountingInfo) validateBills(formats strfmt.Registry) error {
	if swag.IsZero(m.Bills) { // not required
		return nil
	}

	for i := 0; i < len(m.Bills); i++ {
		if swag.IsZero(m.Bills[i]) { // not required
			continue
		}

		if m.Bills[i] != nil {
			if err := m.Bills[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Bills" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("Bills" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *AccountingInfo) validatePayments(formats strfmt.Registry) error {
	if swag.IsZero(m.Payments) { // not required
		return nil
	}

	for i := 0; i < len(m.Payments); i++ {
		if swag.IsZero(m.Payments[i]) { // not required
			continue
		}

		if m.Payments[i] != nil {
			if err := m.Payments[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Payments" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("Payments" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this accounting info based on the context it is used
func (m *AccountingInfo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBills(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePayments(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AccountingInfo) contextValidateBills(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Bills); i++ {

		if m.Bills[i] != nil {
			if err := m.Bills[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Bills" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("Bills" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *AccountingInfo) contextValidatePayments(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Payments); i++ {

		if m.Payments[i] != nil {
			if err := m.Payments[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("Payments" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("Payments" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *AccountingInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AccountingInfo) UnmarshalBinary(b []byte) error {
	var res AccountingInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
