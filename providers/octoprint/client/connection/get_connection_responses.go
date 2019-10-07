// Code generated by go-swagger; DO NOT EDIT.

package connection

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/kihamo/boggart/providers/octoprint/models"
)

// GetConnectionReader is a Reader for the GetConnection structure.
type GetConnectionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetConnectionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetConnectionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetConnectionOK creates a GetConnectionOK with default headers values
func NewGetConnectionOK() *GetConnectionOK {
	return &GetConnectionOK{}
}

/*GetConnectionOK handles this case with default header values.

Successful operation
*/
type GetConnectionOK struct {
	Payload *models.Connection
}

func (o *GetConnectionOK) Error() string {
	return fmt.Sprintf("[GET /connection][%d] getConnectionOK  %+v", 200, o.Payload)
}

func (o *GetConnectionOK) GetPayload() *models.Connection {
	return o.Payload
}

func (o *GetConnectionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Connection)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
