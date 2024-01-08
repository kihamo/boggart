// Code generated by go-swagger; DO NOT EDIT.

package system

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/hikvision/models"
)

// SetTimeReader is a Reader for the SetTime structure.
type SetTimeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SetTimeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSetTimeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[PUT /System/time] setTime", response, response.Code())
	}
}

// NewSetTimeOK creates a SetTimeOK with default headers values
func NewSetTimeOK() *SetTimeOK {
	return &SetTimeOK{}
}

/*
SetTimeOK describes a response with status code 200, with default header values.

Successful operation
*/
type SetTimeOK struct {
	Payload *models.Status
}

// IsSuccess returns true when this set time o k response has a 2xx status code
func (o *SetTimeOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this set time o k response has a 3xx status code
func (o *SetTimeOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this set time o k response has a 4xx status code
func (o *SetTimeOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this set time o k response has a 5xx status code
func (o *SetTimeOK) IsServerError() bool {
	return false
}

// IsCode returns true when this set time o k response a status code equal to that given
func (o *SetTimeOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the set time o k response
func (o *SetTimeOK) Code() int {
	return 200
}

func (o *SetTimeOK) Error() string {
	return fmt.Sprintf("[PUT /System/time][%d] setTimeOK  %+v", 200, o.Payload)
}

func (o *SetTimeOK) String() string {
	return fmt.Sprintf("[PUT /System/time][%d] setTimeOK  %+v", 200, o.Payload)
}

func (o *SetTimeOK) GetPayload() *models.Status {
	return o.Payload
}

func (o *SetTimeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Status)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
