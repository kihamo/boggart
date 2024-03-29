// Code generated by go-swagger; DO NOT EDIT.

package audio

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// GetDefaultSinkReader is a Reader for the GetDefaultSink structure.
type GetDefaultSinkReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetDefaultSinkReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetDefaultSinkOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetDefaultSinkNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetDefaultSinkOK creates a GetDefaultSinkOK with default headers values
func NewGetDefaultSinkOK() *GetDefaultSinkOK {
	return &GetDefaultSinkOK{}
}

/*GetDefaultSinkOK handles this case with default header values.

OK
*/
type GetDefaultSinkOK struct {
}

func (o *GetDefaultSinkOK) Error() string {
	return fmt.Sprintf("[GET /audio/defaultsink][%d] getDefaultSinkOK ", 200)
}

func (o *GetDefaultSinkOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetDefaultSinkNotFound creates a GetDefaultSinkNotFound with default headers values
func NewGetDefaultSinkNotFound() *GetDefaultSinkNotFound {
	return &GetDefaultSinkNotFound{}
}

/*GetDefaultSinkNotFound handles this case with default header values.

Sink not found
*/
type GetDefaultSinkNotFound struct {
}

func (o *GetDefaultSinkNotFound) Error() string {
	return fmt.Sprintf("[GET /audio/defaultsink][%d] getDefaultSinkNotFound ", 404)
}

func (o *GetDefaultSinkNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
