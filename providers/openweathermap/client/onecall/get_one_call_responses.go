// Code generated by go-swagger; DO NOT EDIT.

package onecall

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/openweathermap/models"
)

// GetOneCallReader is a Reader for the GetOneCall structure.
type GetOneCallReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetOneCallReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetOneCallOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 429:
		result := NewGetOneCallTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetOneCallDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetOneCallOK creates a GetOneCallOK with default headers values
func NewGetOneCallOK() *GetOneCallOK {
	return &GetOneCallOK{}
}

/*GetOneCallOK handles this case with default header values.

Successful operation
*/
type GetOneCallOK struct {
	Payload *models.OneCall
}

func (o *GetOneCallOK) Error() string {
	return fmt.Sprintf("[GET /data/2.5/onecall?lat={lat}&lon={lon}][%d] getOneCallOK  %+v", 200, o.Payload)
}

func (o *GetOneCallOK) GetPayload() *models.OneCall {
	return o.Payload
}

func (o *GetOneCallOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.OneCall)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetOneCallTooManyRequests creates a GetOneCallTooManyRequests with default headers values
func NewGetOneCallTooManyRequests() *GetOneCallTooManyRequests {
	return &GetOneCallTooManyRequests{}
}

/*GetOneCallTooManyRequests handles this case with default header values.

Account is blocked
*/
type GetOneCallTooManyRequests struct {
	Payload *models.Error
}

func (o *GetOneCallTooManyRequests) Error() string {
	return fmt.Sprintf("[GET /data/2.5/onecall?lat={lat}&lon={lon}][%d] getOneCallTooManyRequests  %+v", 429, o.Payload)
}

func (o *GetOneCallTooManyRequests) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetOneCallTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetOneCallDefault creates a GetOneCallDefault with default headers values
func NewGetOneCallDefault(code int) *GetOneCallDefault {
	return &GetOneCallDefault{
		_statusCode: code,
	}
}

/*GetOneCallDefault handles this case with default header values.

Unexpected error
*/
type GetOneCallDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get one call default response
func (o *GetOneCallDefault) Code() int {
	return o._statusCode
}

func (o *GetOneCallDefault) Error() string {
	return fmt.Sprintf("[GET /data/2.5/onecall?lat={lat}&lon={lon}][%d] getOneCall default  %+v", o._statusCode, o.Payload)
}

func (o *GetOneCallDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetOneCallDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
