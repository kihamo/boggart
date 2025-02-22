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

// GetCurrentByGeographicCoordinatesReader is a Reader for the GetCurrentByGeographicCoordinates structure.
type GetCurrentByGeographicCoordinatesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCurrentByGeographicCoordinatesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetCurrentByGeographicCoordinatesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetCurrentByGeographicCoordinatesTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetCurrentByGeographicCoordinatesDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetCurrentByGeographicCoordinatesOK creates a GetCurrentByGeographicCoordinatesOK with default headers values
func NewGetCurrentByGeographicCoordinatesOK() *GetCurrentByGeographicCoordinatesOK {
	return &GetCurrentByGeographicCoordinatesOK{}
}

/*
GetCurrentByGeographicCoordinatesOK describes a response with status code 200, with default header values.

Successful operation
*/
type GetCurrentByGeographicCoordinatesOK struct {
	Payload *models.Current
}

// IsSuccess returns true when this get current by geographic coordinates o k response has a 2xx status code
func (o *GetCurrentByGeographicCoordinatesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get current by geographic coordinates o k response has a 3xx status code
func (o *GetCurrentByGeographicCoordinatesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get current by geographic coordinates o k response has a 4xx status code
func (o *GetCurrentByGeographicCoordinatesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get current by geographic coordinates o k response has a 5xx status code
func (o *GetCurrentByGeographicCoordinatesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get current by geographic coordinates o k response a status code equal to that given
func (o *GetCurrentByGeographicCoordinatesOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get current by geographic coordinates o k response
func (o *GetCurrentByGeographicCoordinatesOK) Code() int {
	return 200
}

func (o *GetCurrentByGeographicCoordinatesOK) Error() string {
	return fmt.Sprintf("[GET /data/2.5/weather?lat={lat}&lon={lon}][%d] getCurrentByGeographicCoordinatesOK  %+v", 200, o.Payload)
}

func (o *GetCurrentByGeographicCoordinatesOK) String() string {
	return fmt.Sprintf("[GET /data/2.5/weather?lat={lat}&lon={lon}][%d] getCurrentByGeographicCoordinatesOK  %+v", 200, o.Payload)
}

func (o *GetCurrentByGeographicCoordinatesOK) GetPayload() *models.Current {
	return o.Payload
}

func (o *GetCurrentByGeographicCoordinatesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Current)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCurrentByGeographicCoordinatesTooManyRequests creates a GetCurrentByGeographicCoordinatesTooManyRequests with default headers values
func NewGetCurrentByGeographicCoordinatesTooManyRequests() *GetCurrentByGeographicCoordinatesTooManyRequests {
	return &GetCurrentByGeographicCoordinatesTooManyRequests{}
}

/*
GetCurrentByGeographicCoordinatesTooManyRequests describes a response with status code 429, with default header values.

Account is blocked
*/
type GetCurrentByGeographicCoordinatesTooManyRequests struct {
	Payload *models.Error
}

// IsSuccess returns true when this get current by geographic coordinates too many requests response has a 2xx status code
func (o *GetCurrentByGeographicCoordinatesTooManyRequests) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get current by geographic coordinates too many requests response has a 3xx status code
func (o *GetCurrentByGeographicCoordinatesTooManyRequests) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get current by geographic coordinates too many requests response has a 4xx status code
func (o *GetCurrentByGeographicCoordinatesTooManyRequests) IsClientError() bool {
	return true
}

// IsServerError returns true when this get current by geographic coordinates too many requests response has a 5xx status code
func (o *GetCurrentByGeographicCoordinatesTooManyRequests) IsServerError() bool {
	return false
}

// IsCode returns true when this get current by geographic coordinates too many requests response a status code equal to that given
func (o *GetCurrentByGeographicCoordinatesTooManyRequests) IsCode(code int) bool {
	return code == 429
}

// Code gets the status code for the get current by geographic coordinates too many requests response
func (o *GetCurrentByGeographicCoordinatesTooManyRequests) Code() int {
	return 429
}

func (o *GetCurrentByGeographicCoordinatesTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /data/2.5/weather?lat={lat}&lon={lon}][%d] getCurrentByGeographicCoordinatesTooManyRequests  %+v", 429, o.Payload)
}

func (o *GetCurrentByGeographicCoordinatesTooManyRequests) String() string {
	return fmt.Sprintf("[GET /data/2.5/weather?lat={lat}&lon={lon}][%d] getCurrentByGeographicCoordinatesTooManyRequests  %+v", 429, o.Payload)
}

func (o *GetCurrentByGeographicCoordinatesTooManyRequests) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetCurrentByGeographicCoordinatesTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCurrentByGeographicCoordinatesDefault creates a GetCurrentByGeographicCoordinatesDefault with default headers values
func NewGetCurrentByGeographicCoordinatesDefault(code int) *GetCurrentByGeographicCoordinatesDefault {
	return &GetCurrentByGeographicCoordinatesDefault{
		_statusCode: code,
	}
}

/*
GetCurrentByGeographicCoordinatesDefault describes a response with status code -1, with default header values.

Unexpected error
*/
type GetCurrentByGeographicCoordinatesDefault struct {
	_statusCode int

	Payload *models.Error
}

// IsSuccess returns true when this get current by geographic coordinates default response has a 2xx status code
func (o *GetCurrentByGeographicCoordinatesDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get current by geographic coordinates default response has a 3xx status code
func (o *GetCurrentByGeographicCoordinatesDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get current by geographic coordinates default response has a 4xx status code
func (o *GetCurrentByGeographicCoordinatesDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get current by geographic coordinates default response has a 5xx status code
func (o *GetCurrentByGeographicCoordinatesDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get current by geographic coordinates default response a status code equal to that given
func (o *GetCurrentByGeographicCoordinatesDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the get current by geographic coordinates default response
func (o *GetCurrentByGeographicCoordinatesDefault) Code() int {
	return o._statusCode
}

func (o *GetCurrentByGeographicCoordinatesDefault) Error() string {
	return fmt.Sprintf("[GET /data/2.5/weather?lat={lat}&lon={lon}][%d] getCurrentByGeographicCoordinates default  %+v", o._statusCode, o.Payload)
}

func (o *GetCurrentByGeographicCoordinatesDefault) String() string {
	return fmt.Sprintf("[GET /data/2.5/weather?lat={lat}&lon={lon}][%d] getCurrentByGeographicCoordinates default  %+v", o._statusCode, o.Payload)
}

func (o *GetCurrentByGeographicCoordinatesDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetCurrentByGeographicCoordinatesDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
