// Code generated by go-swagger; DO NOT EDIT.

package system

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/octoprint/models"
)

// GetCommandsBySourceReader is a Reader for the GetCommandsBySource structure.
type GetCommandsBySourceReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCommandsBySourceReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetCommandsBySourceOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetCommandsBySourceNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetCommandsBySourceOK creates a GetCommandsBySourceOK with default headers values
func NewGetCommandsBySourceOK() *GetCommandsBySourceOK {
	return &GetCommandsBySourceOK{}
}

/* GetCommandsBySourceOK describes a response with status code 200, with default header values.

Successful operation
*/
type GetCommandsBySourceOK struct {
	Payload []*models.Command
}

func (o *GetCommandsBySourceOK) Error() string {
	return fmt.Sprintf("[GET /api/system/commands/{source}][%d] getCommandsBySourceOK  %+v", 200, o.Payload)
}
func (o *GetCommandsBySourceOK) GetPayload() []*models.Command {
	return o.Payload
}

func (o *GetCommandsBySourceOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCommandsBySourceNotFound creates a GetCommandsBySourceNotFound with default headers values
func NewGetCommandsBySourceNotFound() *GetCommandsBySourceNotFound {
	return &GetCommandsBySourceNotFound{}
}

/* GetCommandsBySourceNotFound describes a response with status code 404, with default header values.

Not found
*/
type GetCommandsBySourceNotFound struct {
}

func (o *GetCommandsBySourceNotFound) Error() string {
	return fmt.Sprintf("[GET /api/system/commands/{source}][%d] getCommandsBySourceNotFound ", 404)
}

func (o *GetCommandsBySourceNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
