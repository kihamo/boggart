// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Document document
// swagger:model Document
type Document struct {

	// receipt
	Receipt *DocumentReceipt `json:"receipt,omitempty"`
}

// Validate validates this document
func (m *Document) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateReceipt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Document) validateReceipt(formats strfmt.Registry) error {

	if swag.IsZero(m.Receipt) { // not required
		return nil
	}

	if m.Receipt != nil {
		if err := m.Receipt.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("receipt")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Document) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Document) UnmarshalBinary(b []byte) error {
	var res Document
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DocumentReceipt document receipt
// swagger:model DocumentReceipt
type DocumentReceipt struct {

	// cash total sum
	CashTotalSum uint64 `json:"cashTotalSum,omitempty"`

	// date time
	// Format: datetime
	DateTime strfmt.DateTime `json:"dateTime,omitempty"`

	// ecash total sum
	EcashTotalSum uint64 `json:"ecashTotalSum,omitempty"`

	// fiscal document number
	FiscalDocumentNumber uint64 `json:"fiscalDocumentNumber,omitempty"`

	// fiscal drive number
	FiscalDriveNumber string `json:"fiscalDriveNumber,omitempty"`

	// fiscal sign
	FiscalSign uint64 `json:"fiscalSign,omitempty"`

	// items
	Items []*DocumentReceiptItemsItems0 `json:"items"`

	// kkt reg Id
	KktRegID string `json:"kktRegId,omitempty"`

	// nds18
	Nds18 uint64 `json:"nds18,omitempty"`

	// operation type
	OperationType uint64 `json:"operationType,omitempty"`

	// operator
	Operator string `json:"operator,omitempty"`

	// raw data
	// Format: byte
	RawData strfmt.Base64 `json:"rawData,omitempty"`

	// receipt code
	ReceiptCode uint64 `json:"receiptCode,omitempty"`

	// request number
	RequestNumber uint64 `json:"requestNumber,omitempty"`

	// shift number
	ShiftNumber uint64 `json:"shiftNumber,omitempty"`

	// taxation type
	TaxationType uint64 `json:"taxationType,omitempty"`

	// total sum
	TotalSum uint64 `json:"totalSum,omitempty"`

	// user inn
	UserInn string `json:"userInn,omitempty"`
}

// Validate validates this document receipt
func (m *DocumentReceipt) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDateTime(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateItems(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRawData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DocumentReceipt) validateDateTime(formats strfmt.Registry) error {

	if swag.IsZero(m.DateTime) { // not required
		return nil
	}

	if err := validate.FormatOf("receipt"+"."+"dateTime", "body", "datetime", m.DateTime.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *DocumentReceipt) validateItems(formats strfmt.Registry) error {

	if swag.IsZero(m.Items) { // not required
		return nil
	}

	for i := 0; i < len(m.Items); i++ {
		if swag.IsZero(m.Items[i]) { // not required
			continue
		}

		if m.Items[i] != nil {
			if err := m.Items[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("receipt" + "." + "items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *DocumentReceipt) validateRawData(formats strfmt.Registry) error {

	if swag.IsZero(m.RawData) { // not required
		return nil
	}

	// Format "byte" (base64 string) is already validated when unmarshalled

	return nil
}

// MarshalBinary interface implementation
func (m *DocumentReceipt) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DocumentReceipt) UnmarshalBinary(b []byte) error {
	var res DocumentReceipt
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DocumentReceiptItemsItems0 document receipt items items0
// swagger:model DocumentReceiptItemsItems0
type DocumentReceiptItemsItems0 struct {

	// name
	Name string `json:"name,omitempty"`

	// nds18
	Nds18 uint64 `json:"nds18,omitempty"`

	// price
	Price uint64 `json:"price,omitempty"`

	// quantity
	Quantity uint64 `json:"quantity,omitempty"`

	// sum
	Sum uint64 `json:"sum,omitempty"`
}

// Validate validates this document receipt items items0
func (m *DocumentReceiptItemsItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DocumentReceiptItemsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DocumentReceiptItemsItems0) UnmarshalBinary(b []byte) error {
	var res DocumentReceiptItemsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}