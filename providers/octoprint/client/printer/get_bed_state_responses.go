// Code generated by go-swagger; DO NOT EDIT.

package printer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/octoprint/models"
)

// GetBedStateReader is a Reader for the GetBedState structure.
type GetBedStateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetBedStateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetBedStateOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 409:
		result := NewGetBedStateConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /api/printer/bed] getBedState", response, response.Code())
	}
}

// NewGetBedStateOK creates a GetBedStateOK with default headers values
func NewGetBedStateOK() *GetBedStateOK {
	return &GetBedStateOK{}
}

/*
GetBedStateOK describes a response with status code 200, with default header values.

Successful operation
*/
type GetBedStateOK struct {
	Payload *models.BedState
}

// IsSuccess returns true when this get bed state o k response has a 2xx status code
func (o *GetBedStateOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get bed state o k response has a 3xx status code
func (o *GetBedStateOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get bed state o k response has a 4xx status code
func (o *GetBedStateOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get bed state o k response has a 5xx status code
func (o *GetBedStateOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get bed state o k response a status code equal to that given
func (o *GetBedStateOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get bed state o k response
func (o *GetBedStateOK) Code() int {
	return 200
}

func (o *GetBedStateOK) Error() string {
	return fmt.Sprintf("[GET /api/printer/bed][%d] getBedStateOK  %+v", 200, o.Payload)
}

func (o *GetBedStateOK) String() string {
	return fmt.Sprintf("[GET /api/printer/bed][%d] getBedStateOK  %+v", 200, o.Payload)
}

func (o *GetBedStateOK) GetPayload() *models.BedState {
	return o.Payload
}

func (o *GetBedStateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.BedState)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetBedStateConflict creates a GetBedStateConflict with default headers values
func NewGetBedStateConflict() *GetBedStateConflict {
	return &GetBedStateConflict{}
}

/*
GetBedStateConflict describes a response with status code 409, with default header values.

If the printer is not operational
*/
type GetBedStateConflict struct {
}

// IsSuccess returns true when this get bed state conflict response has a 2xx status code
func (o *GetBedStateConflict) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get bed state conflict response has a 3xx status code
func (o *GetBedStateConflict) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get bed state conflict response has a 4xx status code
func (o *GetBedStateConflict) IsClientError() bool {
	return true
}

// IsServerError returns true when this get bed state conflict response has a 5xx status code
func (o *GetBedStateConflict) IsServerError() bool {
	return false
}

// IsCode returns true when this get bed state conflict response a status code equal to that given
func (o *GetBedStateConflict) IsCode(code int) bool {
	return code == 409
}

// Code gets the status code for the get bed state conflict response
func (o *GetBedStateConflict) Code() int {
	return 409
}

func (o *GetBedStateConflict) Error() string {
	return fmt.Sprintf("[GET /api/printer/bed][%d] getBedStateConflict ", 409)
}

func (o *GetBedStateConflict) String() string {
	return fmt.Sprintf("[GET /api/printer/bed][%d] getBedStateConflict ", 409)
}

func (o *GetBedStateConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
