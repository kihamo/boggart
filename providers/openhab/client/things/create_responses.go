// Code generated by go-swagger; DO NOT EDIT.

package things

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// CreateReader is a Reader for the Create structure.
type CreateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewCreateBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewCreateConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewCreateCreated creates a CreateCreated with default headers values
func NewCreateCreated() *CreateCreated {
	return &CreateCreated{}
}

/*CreateCreated handles this case with default header values.

Created
*/
type CreateCreated struct {
	Payload string
}

func (o *CreateCreated) Error() string {
	return fmt.Sprintf("[POST /things][%d] createCreated  %+v", 201, o.Payload)
}

func (o *CreateCreated) GetPayload() string {
	return o.Payload
}

func (o *CreateCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateBadRequest creates a CreateBadRequest with default headers values
func NewCreateBadRequest() *CreateBadRequest {
	return &CreateBadRequest{}
}

/*CreateBadRequest handles this case with default header values.

A uid must be provided, if no binding can create a thing of this type.
*/
type CreateBadRequest struct {
}

func (o *CreateBadRequest) Error() string {
	return fmt.Sprintf("[POST /things][%d] createBadRequest ", 400)
}

func (o *CreateBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateConflict creates a CreateConflict with default headers values
func NewCreateConflict() *CreateConflict {
	return &CreateConflict{}
}

/*CreateConflict handles this case with default header values.

A thing with the same uid already exists.
*/
type CreateConflict struct {
}

func (o *CreateConflict) Error() string {
	return fmt.Sprintf("[POST /things][%d] createConflict ", 409)
}

func (o *CreateConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
