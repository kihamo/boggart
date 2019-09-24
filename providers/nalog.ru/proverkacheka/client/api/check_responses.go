// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// CheckReader is a Reader for the Check structure.
type CheckReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CheckReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewCheckNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewCheckNoContent creates a CheckNoContent with default headers values
func NewCheckNoContent() *CheckNoContent {
	return &CheckNoContent{}
}

/*CheckNoContent handles this case with default header values.

Successful operation
*/
type CheckNoContent struct {
}

func (o *CheckNoContent) Error() string {
	return fmt.Sprintf("[GET /ofds/*/inns/*/fss/{fiscalDriveNumber}/operations/1/tickets/{fiscalDocumentNumber}][%d] checkNoContent ", 204)
}

func (o *CheckNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}