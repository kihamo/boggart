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

/*GetCurrentByCityNameOK handles this case with default header values.

Successful operation
*/
type GetCurrentByCityNameOK struct {
	Payload *models.Current
}

func (o *GetCurrentByCityNameOK) Error() string {
	return fmt.Sprintf("[GET /data/2.5/weather?q={q}&mode=json][%d] getCurrentByCityNameOK  %+v", 200, o.Payload)
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

// NewGetCurrentByCityNameDefault creates a GetCurrentByCityNameDefault with default headers values
func NewGetCurrentByCityNameDefault(code int) *GetCurrentByCityNameDefault {
	return &GetCurrentByCityNameDefault{
		_statusCode: code,
	}
}

/*GetCurrentByCityNameDefault handles this case with default header values.

Unexpected error
*/
type GetCurrentByCityNameDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get current by city name default response
func (o *GetCurrentByCityNameDefault) Code() int {
	return o._statusCode
}

func (o *GetCurrentByCityNameDefault) Error() string {
	return fmt.Sprintf("[GET /data/2.5/weather?q={q}&mode=json][%d] getCurrentByCityName default  %+v", o._statusCode, o.Payload)
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
