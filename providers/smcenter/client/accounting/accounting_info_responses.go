// Code generated by go-swagger; DO NOT EDIT.

package accounting

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/kihamo/boggart/providers/smcenter/models"
)

// AccountingInfoReader is a Reader for the AccountingInfo structure.
type AccountingInfoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AccountingInfoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewAccountingInfoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewAccountingInfoUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewAccountingInfoDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewAccountingInfoOK creates a AccountingInfoOK with default headers values
func NewAccountingInfoOK() *AccountingInfoOK {
	return &AccountingInfoOK{}
}

/* AccountingInfoOK describes a response with status code 200, with default header values.

Successful operation
*/
type AccountingInfoOK struct {
	Payload *AccountingInfoOKBody
}

func (o *AccountingInfoOK) Error() string {
	return fmt.Sprintf("[GET /Accounting/Info][%d] accountingInfoOK  %+v", 200, o.Payload)
}
func (o *AccountingInfoOK) GetPayload() *AccountingInfoOKBody {
	return o.Payload
}

func (o *AccountingInfoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(AccountingInfoOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAccountingInfoUnauthorized creates a AccountingInfoUnauthorized with default headers values
func NewAccountingInfoUnauthorized() *AccountingInfoUnauthorized {
	return &AccountingInfoUnauthorized{}
}

/* AccountingInfoUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type AccountingInfoUnauthorized struct {
	Payload *models.Error
}

func (o *AccountingInfoUnauthorized) Error() string {
	return fmt.Sprintf("[GET /Accounting/Info][%d] accountingInfoUnauthorized  %+v", 401, o.Payload)
}
func (o *AccountingInfoUnauthorized) GetPayload() *models.Error {
	return o.Payload
}

func (o *AccountingInfoUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAccountingInfoDefault creates a AccountingInfoDefault with default headers values
func NewAccountingInfoDefault(code int) *AccountingInfoDefault {
	return &AccountingInfoDefault{
		_statusCode: code,
	}
}

/* AccountingInfoDefault describes a response with status code -1, with default header values.

Unexpected error
*/
type AccountingInfoDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the accounting info default response
func (o *AccountingInfoDefault) Code() int {
	return o._statusCode
}

func (o *AccountingInfoDefault) Error() string {
	return fmt.Sprintf("[GET /Accounting/Info][%d] accountingInfo default  %+v", o._statusCode, o.Payload)
}
func (o *AccountingInfoDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *AccountingInfoDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*AccountingInfoOKBody accounting info o k body
swagger:model AccountingInfoOKBody
*/
type AccountingInfoOKBody struct {

	// data
	Data []*models.AccountingInfo `json:"Data"`
}

// Validate validates this accounting info o k body
func (o *AccountingInfoOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *AccountingInfoOKBody) validateData(formats strfmt.Registry) error {
	if swag.IsZero(o.Data) { // not required
		return nil
	}

	for i := 0; i < len(o.Data); i++ {
		if swag.IsZero(o.Data[i]) { // not required
			continue
		}

		if o.Data[i] != nil {
			if err := o.Data[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("accountingInfoOK" + "." + "Data" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("accountingInfoOK" + "." + "Data" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this accounting info o k body based on the context it is used
func (o *AccountingInfoOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *AccountingInfoOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Data); i++ {

		if o.Data[i] != nil {
			if err := o.Data[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("accountingInfoOK" + "." + "Data" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("accountingInfoOK" + "." + "Data" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *AccountingInfoOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AccountingInfoOKBody) UnmarshalBinary(b []byte) error {
	var res AccountingInfoOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
