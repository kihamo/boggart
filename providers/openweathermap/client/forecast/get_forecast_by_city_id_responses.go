// Code generated by go-swagger; DO NOT EDIT.

package forecast

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/openweathermap/models"
)

// GetForecastByCityIDReader is a Reader for the GetForecastByCityID structure.
type GetForecastByCityIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetForecastByCityIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetForecastByCityIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetForecastByCityIDTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetForecastByCityIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetForecastByCityIDOK creates a GetForecastByCityIDOK with default headers values
func NewGetForecastByCityIDOK() *GetForecastByCityIDOK {
	return &GetForecastByCityIDOK{}
}

/*
GetForecastByCityIDOK describes a response with status code 200, with default header values.

Successful operation
*/
type GetForecastByCityIDOK struct {
	Payload *models.Forecast
}

// IsSuccess returns true when this get forecast by city Id o k response has a 2xx status code
func (o *GetForecastByCityIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get forecast by city Id o k response has a 3xx status code
func (o *GetForecastByCityIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get forecast by city Id o k response has a 4xx status code
func (o *GetForecastByCityIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get forecast by city Id o k response has a 5xx status code
func (o *GetForecastByCityIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get forecast by city Id o k response a status code equal to that given
func (o *GetForecastByCityIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get forecast by city Id o k response
func (o *GetForecastByCityIDOK) Code() int {
	return 200
}

func (o *GetForecastByCityIDOK) Error() string {
	return fmt.Sprintf("[GET /data/2.5/forecast?id={id}][%d] getForecastByCityIdOK  %+v", 200, o.Payload)
}

func (o *GetForecastByCityIDOK) String() string {
	return fmt.Sprintf("[GET /data/2.5/forecast?id={id}][%d] getForecastByCityIdOK  %+v", 200, o.Payload)
}

func (o *GetForecastByCityIDOK) GetPayload() *models.Forecast {
	return o.Payload
}

func (o *GetForecastByCityIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Forecast)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetForecastByCityIDTooManyRequests creates a GetForecastByCityIDTooManyRequests with default headers values
func NewGetForecastByCityIDTooManyRequests() *GetForecastByCityIDTooManyRequests {
	return &GetForecastByCityIDTooManyRequests{}
}

/*
GetForecastByCityIDTooManyRequests describes a response with status code 429, with default header values.

Account is blocked
*/
type GetForecastByCityIDTooManyRequests struct {
	Payload *models.Error
}

// IsSuccess returns true when this get forecast by city Id too many requests response has a 2xx status code
func (o *GetForecastByCityIDTooManyRequests) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get forecast by city Id too many requests response has a 3xx status code
func (o *GetForecastByCityIDTooManyRequests) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get forecast by city Id too many requests response has a 4xx status code
func (o *GetForecastByCityIDTooManyRequests) IsClientError() bool {
	return true
}

// IsServerError returns true when this get forecast by city Id too many requests response has a 5xx status code
func (o *GetForecastByCityIDTooManyRequests) IsServerError() bool {
	return false
}

// IsCode returns true when this get forecast by city Id too many requests response a status code equal to that given
func (o *GetForecastByCityIDTooManyRequests) IsCode(code int) bool {
	return code == 429
}

// Code gets the status code for the get forecast by city Id too many requests response
func (o *GetForecastByCityIDTooManyRequests) Code() int {
	return 429
}

func (o *GetForecastByCityIDTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /data/2.5/forecast?id={id}][%d] getForecastByCityIdTooManyRequests  %+v", 429, o.Payload)
}

func (o *GetForecastByCityIDTooManyRequests) String() string {
	return fmt.Sprintf("[GET /data/2.5/forecast?id={id}][%d] getForecastByCityIdTooManyRequests  %+v", 429, o.Payload)
}

func (o *GetForecastByCityIDTooManyRequests) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetForecastByCityIDTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetForecastByCityIDDefault creates a GetForecastByCityIDDefault with default headers values
func NewGetForecastByCityIDDefault(code int) *GetForecastByCityIDDefault {
	return &GetForecastByCityIDDefault{
		_statusCode: code,
	}
}

/*
GetForecastByCityIDDefault describes a response with status code -1, with default header values.

Unexpected error
*/
type GetForecastByCityIDDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this get forecast by city ID default response has a 2xx status code
func (o *GetForecastByCityIDDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get forecast by city ID default response has a 3xx status code
func (o *GetForecastByCityIDDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get forecast by city ID default response has a 4xx status code
func (o *GetForecastByCityIDDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get forecast by city ID default response has a 5xx status code
func (o *GetForecastByCityIDDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get forecast by city ID default response a status code equal to that given
func (o *GetForecastByCityIDDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get forecast by city ID default response
func (o *GetForecastByCityIDDefault) Code() int {
	return o._statusCode
}

func (o *GetForecastByCityIDDefault) Error() string {
	return fmt.Sprintf("[GET /data/2.5/forecast?id={id}][%d] getForecastByCityID default  %+v", o._statusCode, o.Payload)
}

func (o *GetForecastByCityIDDefault) String() string {
	return fmt.Sprintf("[GET /data/2.5/forecast?id={id}][%d] getForecastByCityID default  %+v", o._statusCode, o.Payload)
}

func (o *GetForecastByCityIDDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetForecastByCityIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
