// Code generated by go-swagger; DO NOT EDIT.

package discovery

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// ScanReader is a Reader for the Scan structure.
type ScanReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ScanReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewScanOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewScanOK creates a ScanOK with default headers values
func NewScanOK() *ScanOK {
	return &ScanOK{}
}

/*ScanOK handles this case with default header values.

OK
*/
type ScanOK struct {
	Payload int32
}

func (o *ScanOK) Error() string {
	return fmt.Sprintf("[POST /discovery/bindings/{bindingId}/scan][%d] scanOK  %+v", 200, o.Payload)
}

func (o *ScanOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}