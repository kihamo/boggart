// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// OperationResultPlainText operation result plain text
//
// swagger:model OperationResultPlainText
type OperationResultPlainText string

func NewOperationResultPlainText(value OperationResultPlainText) *OperationResultPlainText {
	return &value
}

// Pointer returns a pointer to a freshly-allocated OperationResultPlainText.
func (m OperationResultPlainText) Pointer() *OperationResultPlainText {
	return &m
}

const (

	// OperationResultPlainTextFAIL captures enum value "FAIL"
	OperationResultPlainTextFAIL OperationResultPlainText = "FAIL"

	// OperationResultPlainTextOK captures enum value "OK"
	OperationResultPlainTextOK OperationResultPlainText = "OK"
)

// for schema
var operationResultPlainTextEnum []interface{}

func init() {
	var res []OperationResultPlainText
	if err := json.Unmarshal([]byte(`["FAIL","OK"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		operationResultPlainTextEnum = append(operationResultPlainTextEnum, v)
	}
}

func (m OperationResultPlainText) validateOperationResultPlainTextEnum(path, location string, value OperationResultPlainText) error {
	if err := validate.EnumCase(path, location, value, operationResultPlainTextEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this operation result plain text
func (m OperationResultPlainText) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateOperationResultPlainTextEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this operation result plain text based on context it is used
func (m OperationResultPlainText) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
