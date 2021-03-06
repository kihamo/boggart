// Code generated by go-swagger; DO NOT EDIT.

package links

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// IsAutomaticReader is a Reader for the IsAutomatic structure.
type IsAutomaticReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *IsAutomaticReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewIsAutomaticOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewIsAutomaticOK creates a IsAutomaticOK with default headers values
func NewIsAutomaticOK() *IsAutomaticOK {
	return &IsAutomaticOK{}
}

/*IsAutomaticOK handles this case with default header values.

OK
*/
type IsAutomaticOK struct {
	Payload bool
}

func (o *IsAutomaticOK) Error() string {
	return fmt.Sprintf("[GET /links/auto][%d] isAutomaticOK  %+v", 200, o.Payload)
}

func (o *IsAutomaticOK) GetPayload() bool {
	return o.Payload
}

func (o *IsAutomaticOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
