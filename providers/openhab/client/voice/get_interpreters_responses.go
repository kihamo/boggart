// Code generated by go-swagger; DO NOT EDIT.

package voice

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// GetInterpretersReader is a Reader for the GetInterpreters structure.
type GetInterpretersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetInterpretersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetInterpretersOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetInterpretersOK creates a GetInterpretersOK with default headers values
func NewGetInterpretersOK() *GetInterpretersOK {
	return &GetInterpretersOK{}
}

/*GetInterpretersOK handles this case with default header values.

OK
*/
type GetInterpretersOK struct {
}

func (o *GetInterpretersOK) Error() string {
	return fmt.Sprintf("[GET /voice/interpreters][%d] getInterpretersOK ", 200)
}

func (o *GetInterpretersOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
