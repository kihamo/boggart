// Code generated by go-swagger; DO NOT EDIT.

package eq

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// EQGetListReader is a Reader for the EQGetList structure.
type EQGetListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *EQGetListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewEQGetListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /httpapi.asp?command=EQGetList] EQGetList", response, response.Code())
	}
}

// NewEQGetListOK creates a EQGetListOK with default headers values
func NewEQGetListOK() *EQGetListOK {
	return &EQGetListOK{}
}

/*
EQGetListOK describes a response with status code 200, with default header values.

Successful
*/
type EQGetListOK struct {
	Payload []string
}

// IsSuccess returns true when this e q get list o k response has a 2xx status code
func (o *EQGetListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this e q get list o k response has a 3xx status code
func (o *EQGetListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this e q get list o k response has a 4xx status code
func (o *EQGetListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this e q get list o k response has a 5xx status code
func (o *EQGetListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this e q get list o k response a status code equal to that given
func (o *EQGetListOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the e q get list o k response
func (o *EQGetListOK) Code() int {
	return 200
}

func (o *EQGetListOK) Error() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=EQGetList][%d] eQGetListOK  %+v", 200, o.Payload)
}

func (o *EQGetListOK) String() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=EQGetList][%d] eQGetListOK  %+v", 200, o.Payload)
}

func (o *EQGetListOK) GetPayload() []string {
	return o.Payload
}

func (o *EQGetListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}