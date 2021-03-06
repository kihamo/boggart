// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// Debt debt
// swagger:model Debt
type Debt struct {

	// address
	Address string `json:"Address,omitempty"`

	// debt actual date
	DebtActualDate string `json:"DebtActualDate,omitempty"`

	// ident
	Ident string `json:"Ident,omitempty"`

	// meters access flag
	MetersAccessFlag bool `json:"MetersAccessFlag,omitempty"`

	// meters end day
	MetersEndDay int64 `json:"MetersEndDay,omitempty"`

	// meters start day
	MetersStartDay int64 `json:"MetersStartDay,omitempty"`

	// sum
	Sum float64 `json:"Sum,omitempty"`

	// sum fine
	SumFine float64 `json:"SumFine,omitempty"`
}

// Validate validates this debt
func (m *Debt) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Debt) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Debt) UnmarshalBinary(b []byte) error {
	var res Debt
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
