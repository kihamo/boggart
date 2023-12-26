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
	"github.com/go-openapi/validate"

	"github.com/kihamo/boggart/providers/myheat/cloud/models"
)

// GetDeviceInfoReader is a Reader for the GetDeviceInfo structure.
type GetDeviceInfoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetDeviceInfoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetDeviceInfoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetDeviceInfoDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetDeviceInfoOK creates a GetDeviceInfoOK with default headers values
func NewGetDeviceInfoOK() *GetDeviceInfoOK {
	return &GetDeviceInfoOK{}
}

/*
GetDeviceInfoOK describes a response with status code 200, with default header values.

Successful
*/
type GetDeviceInfoOK struct {
	Payload *GetDeviceInfoOKBody
}

// IsSuccess returns true when this get device info o k response has a 2xx status code
func (o *GetDeviceInfoOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get device info o k response has a 3xx status code
func (o *GetDeviceInfoOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get device info o k response has a 4xx status code
func (o *GetDeviceInfoOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get device info o k response has a 5xx status code
func (o *GetDeviceInfoOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get device info o k response a status code equal to that given
func (o *GetDeviceInfoOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get device info o k response
func (o *GetDeviceInfoOK) Code() int {
	return 200
}

func (o *GetDeviceInfoOK) Error() string {
	return fmt.Sprintf("[POST /request/?getDeviceInfo][%d] getDeviceInfoOK  %+v", 200, o.Payload)
}

func (o *GetDeviceInfoOK) String() string {
	return fmt.Sprintf("[POST /request/?getDeviceInfo][%d] getDeviceInfoOK  %+v", 200, o.Payload)
}

func (o *GetDeviceInfoOK) GetPayload() *GetDeviceInfoOKBody {
	return o.Payload
}

func (o *GetDeviceInfoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetDeviceInfoOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetDeviceInfoDefault creates a GetDeviceInfoDefault with default headers values
func NewGetDeviceInfoDefault(code int) *GetDeviceInfoDefault {
	return &GetDeviceInfoDefault{
		_statusCode: code,
	}
}

/*
GetDeviceInfoDefault describes a response with status code -1, with default header values.

Unexpected error
*/
type GetDeviceInfoDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this get device info default response has a 2xx status code
func (o *GetDeviceInfoDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get device info default response has a 3xx status code
func (o *GetDeviceInfoDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get device info default response has a 4xx status code
func (o *GetDeviceInfoDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get device info default response has a 5xx status code
func (o *GetDeviceInfoDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get device info default response a status code equal to that given
func (o *GetDeviceInfoDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get device info default response
func (o *GetDeviceInfoDefault) Code() int {
	return o._statusCode
}

func (o *GetDeviceInfoDefault) Error() string {
	return fmt.Sprintf("[POST /request/?getDeviceInfo][%d] getDeviceInfo default  %+v", o._statusCode, o.Payload)
}

func (o *GetDeviceInfoDefault) String() string {
	return fmt.Sprintf("[POST /request/?getDeviceInfo][%d] getDeviceInfo default  %+v", o._statusCode, o.Payload)
}

func (o *GetDeviceInfoDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetDeviceInfoDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
GetDeviceInfoBody get device info body
swagger:model GetDeviceInfoBody
*/
type GetDeviceInfoBody struct {

	// device Id
	// Required: true
	DeviceID int64 `json:"deviceId"`
}

// Validate validates this get device info body
func (o *GetDeviceInfoBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDeviceID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetDeviceInfoBody) validateDeviceID(formats strfmt.Registry) error {

	if err := validate.Required("request"+"."+"deviceId", "body", int64(o.DeviceID)); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this get device info body based on context it is used
func (o *GetDeviceInfoBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetDeviceInfoBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDeviceInfoBody) UnmarshalBinary(b []byte) error {
	var res GetDeviceInfoBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*
GetDeviceInfoOKBody get device info o k body
swagger:model GetDeviceInfoOKBody
*/
type GetDeviceInfoOKBody struct {

	// data
	Data *models.ResponseDeviceInfo `json:"data,omitempty"`

	// error
	Error int64 `json:"error,omitempty"`

	// refresh page
	RefreshPage bool `json:"refreshPage,omitempty"`
}

// Validate validates this get device info o k body
func (o *GetDeviceInfoOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetDeviceInfoOKBody) validateData(formats strfmt.Registry) error {
	if swag.IsZero(o.Data) { // not required
		return nil
	}

	if o.Data != nil {
		if err := o.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getDeviceInfoOK" + "." + "data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getDeviceInfoOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this get device info o k body based on the context it is used
func (o *GetDeviceInfoOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetDeviceInfoOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if o.Data != nil {

		if swag.IsZero(o.Data) { // not required
			return nil
		}

		if err := o.Data.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getDeviceInfoOK" + "." + "data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getDeviceInfoOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetDeviceInfoOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetDeviceInfoOKBody) UnmarshalBinary(b []byte) error {
	var res GetDeviceInfoOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
