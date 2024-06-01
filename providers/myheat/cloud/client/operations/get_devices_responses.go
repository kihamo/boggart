// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/kihamo/boggart/providers/myheat/cloud/models"
)

// GetDevicesReader is a Reader for the GetDevices structure.
type GetDevicesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetDevicesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetDevicesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetDevicesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetDevicesOK creates a GetDevicesOK with default headers values
func NewGetDevicesOK() *GetDevicesOK {
	return &GetDevicesOK{}
}

/*
GetDevicesOK describes a response with status code 200, with default header values.

Successful
*/
type GetDevicesOK struct {
	Payload *GetDevicesOKBody
}

// IsSuccess returns true when this get devices o k response has a 2xx status code
func (o *GetDevicesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get devices o k response has a 3xx status code
func (o *GetDevicesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get devices o k response has a 4xx status code
func (o *GetDevicesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get devices o k response has a 5xx status code
func (o *GetDevicesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get devices o k response a status code equal to that given
func (o *GetDevicesOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get devices o k response
func (o *GetDevicesOK) Code() int {
	return 200
}

func (o *GetDevicesOK) Error() string {
	return fmt.Sprintf("[POST /request/?getDevices][%d] getDevicesOK  %+v", 200, o.Payload)
}

func (o *GetDevicesOK) String() string {
	return fmt.Sprintf("[POST /request/?getDevices][%d] getDevicesOK  %+v", 200, o.Payload)
}

func (o *GetDevicesOK) GetPayload() *GetDevicesOKBody {
	return o.Payload
}

func (o *GetDevicesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetDevicesOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetDevicesDefault creates a GetDevicesDefault with default headers values
func NewGetDevicesDefault(code int) *GetDevicesDefault {
	return &GetDevicesDefault{
		_statusCode: code,
	}
}

/*
GetDevicesDefault describes a response with status code -1, with default header values.

Unexpected error
*/
type GetDevicesDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this get devices default response has a 2xx status code
func (o *GetDevicesDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get devices default response has a 3xx status code
func (o *GetDevicesDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get devices default response has a 4xx status code
func (o *GetDevicesDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get devices default response has a 5xx status code
func (o *GetDevicesDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get devices default response a status code equal to that given
func (o *GetDevicesDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get devices default response
func (o *GetDevicesDefault) Code() int {
	return o._statusCode
}

func (o *GetDevicesDefault) Error() string {
	return fmt.Sprintf("[POST /request/?getDevices][%d] getDevices default  %+v", o._statusCode, o.Payload)
}

func (o *GetDevicesDefault) String() string {
	return fmt.Sprintf("[POST /request/?getDevices][%d] getDevices default  %+v", o._statusCode, o.Payload)
}

func (o *GetDevicesDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetDevicesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
GetDevicesOKBody get devices o k body
swagger:model GetDevicesOKBody
*/
type GetDevicesOKBody struct {

	// data
	Data *models.ResponseDevices `json:"data,omitempty"`

	// error
	Error int64 `json:"error,omitempty"`

	// refresh page
	RefreshPage bool `json:"refreshPage,omitempty"`
}

// Validate validates this get devices o k body
func (o *GetDevicesOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetDevicesOKBody) validateData(formats strfmt.Registry) error {
	if swag.IsZero(o.Data) { // not required
		return nil
	}

	if o.Data != nil {
		if err := o.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getDevicesOK" + "." + "data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getDevicesOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this get devices o k body based on the context it is used
func (o *GetDevicesOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetDevicesOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if o.Data != nil {

		if swag.IsZero(o.Data) { // not required
			return nil
		}

		if err := o.Data.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getDevicesOK" + "." + "data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getDevicesOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetDevicesOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDevicesOKBody) UnmarshalBinary(b []byte) error {
	var res GetDevicesOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}