// Code generated by go-swagger; DO NOT EDIT.

package weather

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/openweathermap/models"
)

// GetCurrentByCityNameReader is a Reader for the GetCurrentByCityName structure.
type GetCurrentByCityNameReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCurrentByCityNameReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetCurrentByCityNameOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetCurrentByCityNameTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetCurrentByCityNameDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetCurrentByCityNameOK creates a GetCurrentByCityNameOK with default headers values
func NewGetCurrentByCityNameOK() *GetCurrentByCityNameOK {
	return &GetCurrentByCityNameOK{}
}

/*
GetCurrentByCityNameOK describes a response with status code 200, with default header values.

Successful operation
*/
type GetCurrentByCityNameOK struct {
	Payload *models.Current
}

// IsSuccess returns true when this get current by city name o k response has a 2xx status code
func (o *GetCurrentByCityNameOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get current by city name o k response has a 3xx status code
func (o *GetCurrentByCityNameOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get current by city name o k response has a 4xx status code
func (o *GetCurrentByCityNameOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get current by city name o k response has a 5xx status code
func (o *GetCurrentByCityNameOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get current by city name o k response a status code equal to that given
func (o *GetCurrentByCityNameOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get current by city name o k response
func (o *GetCurrentByCityNameOK) Code() int {
	return 200
}

func (o *GetCurrentByCityNameOK) Error() string {
	return fmt.Sprintf("[GET /data/2.5/weather?q={q}][%d] getCurrentByCityNameOK  %+v", 200, o.Payload)
}

func (o *GetCurrentByCityNameOK) String() string {
	return fmt.Sprintf("[GET /data/2.5/weather?q={q}][%d] getCurrentByCityNameOK  %+v", 200, o.Payload)
}

func (o *GetCurrentByCityNameOK) GetPayload() *models.Current {
	return o.Payload
}

func (o *GetCurrentByCityNameOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Current)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCurrentByCityNameTooManyRequests creates a GetCurrentByCityNameTooManyRequests with default headers values
func NewGetCurrentByCityNameTooManyRequests() *GetCurrentByCityNameTooManyRequests {
	return &GetCurrentByCityNameTooManyRequests{}
}

/*
GetCurrentByCityNameTooManyRequests describes a response with status code 429, with default header values.

Account is blocked
*/
type GetCurrentByCityNameTooManyRequests struct {
	Payload *models.Error
}

// IsSuccess returns true when this get current by city name too many requests response has a 2xx status code
func (o *GetCurrentByCityNameTooManyRequests) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get current by city name too many requests response has a 3xx status code
func (o *GetCurrentByCityNameTooManyRequests) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get current by city name too many requests response has a 4xx status code
func (o *GetCurrentByCityNameTooManyRequests) IsClientError() bool {
	return true
}

// IsServerError returns true when this get current by city name too many requests response has a 5xx status code
func (o *GetCurrentByCityNameTooManyRequests) IsServerError() bool {
	return false
}

// IsCode returns true when this get current by city name too many requests response a status code equal to that given
func (o *GetCurrentByCityNameTooManyRequests) IsCode(code int) bool {
	return code == 429
}

// Code gets the status code for the get current by city name too many requests response
func (o *GetCurrentByCityNameTooManyRequests) Code() int {
	return 429
}

func (o *GetCurrentByCityNameTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /data/2.5/weather?q={q}][%d] getCurrentByCityNameTooManyRequests  %+v", 429, o.Payload)
}

func (o *GetCurrentByCityNameTooManyRequests) String() string {
	return fmt.Sprintf("[GET /data/2.5/weather?q={q}][%d] getCurrentByCityNameTooManyRequests  %+v", 429, o.Payload)
}

func (o *GetCurrentByCityNameTooManyRequests) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetCurrentByCityNameTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCurrentByCityNameDefault creates a GetCurrentByCityNameDefault with default headers values
func NewGetCurrentByCityNameDefault(code int) *GetCurrentByCityNameDefault {
	return &GetCurrentByCityNameDefault{
		_statusCode: code,
	}
}

/*
GetCurrentByCityNameDefault describes a response with status code -1, with default header values.

Unexpected error
*/
type GetCurrentByCityNameDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this get current by city name default response has a 2xx status code
func (o *GetCurrentByCityNameDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get current by city name default response has a 3xx status code
func (o *GetCurrentByCityNameDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get current by city name default response has a 4xx status code
func (o *GetCurrentByCityNameDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get current by city name default response has a 5xx status code
func (o *GetCurrentByCityNameDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get current by city name default response a status code equal to that given
func (o *GetCurrentByCityNameDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get current by city name default response
func (o *GetCurrentByCityNameDefault) Code() int {
	return o._statusCode
}

func (o *GetCurrentByCityNameDefault) Error() string {
	return fmt.Sprintf("[GET /data/2.5/weather?q={q}][%d] getCurrentByCityName default  %+v", o._statusCode, o.Payload)
}

func (o *GetCurrentByCityNameDefault) String() string {
	return fmt.Sprintf("[GET /data/2.5/weather?q={q}][%d] getCurrentByCityName default  %+v", o._statusCode, o.Payload)
}

func (o *GetCurrentByCityNameDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetCurrentByCityNameDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
