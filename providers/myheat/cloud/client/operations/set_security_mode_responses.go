// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/kihamo/boggart/providers/myheat/cloud/models"
)

// SetSecurityModeReader is a Reader for the SetSecurityMode structure.
type SetSecurityModeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SetSecurityModeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSetSecurityModeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewSetSecurityModeDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewSetSecurityModeOK creates a SetSecurityModeOK with default headers values
func NewSetSecurityModeOK() *SetSecurityModeOK {
	return &SetSecurityModeOK{}
}

/*
SetSecurityModeOK describes a response with status code 200, with default header values.

Successful
*/
type SetSecurityModeOK struct {
	Payload *models.Error
}

// IsSuccess returns true when this set security mode o k response has a 2xx status code
func (o *SetSecurityModeOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this set security mode o k response has a 3xx status code
func (o *SetSecurityModeOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this set security mode o k response has a 4xx status code
func (o *SetSecurityModeOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this set security mode o k response has a 5xx status code
func (o *SetSecurityModeOK) IsServerError() bool {
	return false
}

// IsCode returns true when this set security mode o k response a status code equal to that given
func (o *SetSecurityModeOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the set security mode o k response
func (o *SetSecurityModeOK) Code() int {
	return 200
}

func (o *SetSecurityModeOK) Error() string {
	return fmt.Sprintf("[POST /request/?setSecurityMode][%d] setSecurityModeOK  %+v", 200, o.Payload)
}

func (o *SetSecurityModeOK) String() string {
	return fmt.Sprintf("[POST /request/?setSecurityMode][%d] setSecurityModeOK  %+v", 200, o.Payload)
}

func (o *SetSecurityModeOK) GetPayload() *models.Error {
	return o.Payload
}

func (o *SetSecurityModeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSetSecurityModeDefault creates a SetSecurityModeDefault with default headers values
func NewSetSecurityModeDefault(code int) *SetSecurityModeDefault {
	return &SetSecurityModeDefault{
		_statusCode: code,
	}
}

/*
SetSecurityModeDefault describes a response with status code -1, with default header values.

Unexpected error
*/
type SetSecurityModeDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this set security mode default response has a 2xx status code
func (o *SetSecurityModeDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this set security mode default response has a 3xx status code
func (o *SetSecurityModeDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this set security mode default response has a 4xx status code
func (o *SetSecurityModeDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this set security mode default response has a 5xx status code
func (o *SetSecurityModeDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this set security mode default response a status code equal to that given
func (o *SetSecurityModeDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the set security mode default response
func (o *SetSecurityModeDefault) Code() int {
	return o._statusCode
}

func (o *SetSecurityModeDefault) Error() string {
	return fmt.Sprintf("[POST /request/?setSecurityMode][%d] setSecurityMode default  %+v", o._statusCode, o.Payload)
}

func (o *SetSecurityModeDefault) String() string {
	return fmt.Sprintf("[POST /request/?setSecurityMode][%d] setSecurityMode default  %+v", o._statusCode, o.Payload)
}

func (o *SetSecurityModeDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *SetSecurityModeDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
SetSecurityModeBody set security mode body
swagger:model SetSecurityModeBody
*/
type SetSecurityModeBody struct {

	// Идентификатора контроллера
	// Required: true
	DeviceID int64 `json:"deviceId"`

	// возможные значения: a. 0 – снять с охраны b. 1 – поставить на охрану
	// Enum: [0 1]
	Model int64 `json:"model,omitempty"`

	// Идентификатор инженерного оборудования
	// Required: true
	ObjID int64 `json:"objId"`
}

// Validate validates this set security mode body
func (o *SetSecurityModeBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDeviceID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateModel(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateObjID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SetSecurityModeBody) validateDeviceID(formats strfmt.Registry) error {

	if err := validate.Required("request"+"."+"deviceId", "body", int64(o.DeviceID)); err != nil {
		return err
	}

	return nil
}

var setSecurityModeBodyTypeModelPropEnum []interface{}

func init() {
	var res []int64
	if err := json.Unmarshal([]byte(`[0,1]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		setSecurityModeBodyTypeModelPropEnum = append(setSecurityModeBodyTypeModelPropEnum, v)
	}
}

// prop value enum
func (o *SetSecurityModeBody) validateModelEnum(path, location string, value int64) error {
	if err := validate.EnumCase(path, location, value, setSecurityModeBodyTypeModelPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (o *SetSecurityModeBody) validateModel(formats strfmt.Registry) error {
	if swag.IsZero(o.Model) { // not required
		return nil
	}

	// value enum
	if err := o.validateModelEnum("request"+"."+"model", "body", o.Model); err != nil {
		return err
	}

	return nil
}

func (o *SetSecurityModeBody) validateObjID(formats strfmt.Registry) error {

	if err := validate.Required("request"+"."+"objId", "body", int64(o.ObjID)); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this set security mode body based on context it is used
func (o *SetSecurityModeBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *SetSecurityModeBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SetSecurityModeBody) UnmarshalBinary(b []byte) error {
	var res SetSecurityModeBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}