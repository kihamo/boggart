// Code generated by go-swagger; DO NOT EDIT.

package device

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/wiim/models"
)

// SetShutdownReader is a Reader for the SetShutdown structure.
type SetShutdownReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SetShutdownReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSetShutdownOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /httpapi.asp?command=setShutdown:{seconds}] setShutdown", response, response.Code())
	}
}

// NewSetShutdownOK creates a SetShutdownOK with default headers values
func NewSetShutdownOK() *SetShutdownOK {
	return &SetShutdownOK{}
}

/*
SetShutdownOK describes a response with status code 200, with default header values.

Successful
*/
type SetShutdownOK struct {
	Payload models.OperationResultPlainText
}

// IsSuccess returns true when this set shutdown o k response has a 2xx status code
func (o *SetShutdownOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this set shutdown o k response has a 3xx status code
func (o *SetShutdownOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this set shutdown o k response has a 4xx status code
func (o *SetShutdownOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this set shutdown o k response has a 5xx status code
func (o *SetShutdownOK) IsServerError() bool {
	return false
}

// IsCode returns true when this set shutdown o k response a status code equal to that given
func (o *SetShutdownOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the set shutdown o k response
func (o *SetShutdownOK) Code() int {
	return 200
}

func (o *SetShutdownOK) Error() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=setShutdown:{seconds}][%d] setShutdownOK  %+v", 200, o.Payload)
}

func (o *SetShutdownOK) String() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=setShutdown:{seconds}][%d] setShutdownOK  %+v", 200, o.Payload)
}

func (o *SetShutdownOK) GetPayload() models.OperationResultPlainText {
	return o.Payload
}

func (o *SetShutdownOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
