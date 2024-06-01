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

// SetNtpServerReader is a Reader for the SetNtpServer structure.
type SetNtpServerReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SetNtpServerReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSetNtpServerOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[PUT /System/time/ntpServers/{id}] setNtpServer", response, response.Code())
	}
}

// NewSetNtpServerOK creates a SetNtpServerOK with default headers values
func NewSetNtpServerOK() *SetNtpServerOK {
	return &SetNtpServerOK{}
}

/*
SetNtpServerOK describes a response with status code 200, with default header values.

Successful operation
*/
type SetNtpServerOK struct {
	Payload *models.Status
}

// IsSuccess returns true when this set ntp server o k response has a 2xx status code
func (o *SetNtpServerOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this set ntp server o k response has a 3xx status code
func (o *SetNtpServerOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this set ntp server o k response has a 4xx status code
func (o *SetNtpServerOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this set ntp server o k response has a 5xx status code
func (o *SetNtpServerOK) IsServerError() bool {
	return false
}

// IsCode returns true when this set ntp server o k response a status code equal to that given
func (o *SetNtpServerOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the set ntp server o k response
func (o *SetNtpServerOK) Code() int {
	return 200
}

func (o *SetNtpServerOK) Error() string {
	return fmt.Sprintf("[PUT /System/time/ntpServers/{id}][%d] setNtpServerOK  %+v", 200, o.Payload)
}

func (o *SetNtpServerOK) String() string {
	return fmt.Sprintf("[PUT /System/time/ntpServers/{id}][%d] setNtpServerOK  %+v", 200, o.Payload)
}

func (o *SetNtpServerOK) GetPayload() *models.Status {
	return o.Payload
}

func (o *SetNtpServerOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Status)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}