// Code generated by go-swagger; DO NOT EDIT.

package info

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/kihamo/boggart/components/boggart/providers/google/home/models"
)

// GetSupportedLocalesReader is a Reader for the GetSupportedLocales structure.
type GetSupportedLocalesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSupportedLocalesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetSupportedLocalesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetSupportedLocalesOK creates a GetSupportedLocalesOK with default headers values
func NewGetSupportedLocalesOK() *GetSupportedLocalesOK {
	return &GetSupportedLocalesOK{}
}

/*GetSupportedLocalesOK handles this case with default header values.

returned supported locales list
*/
type GetSupportedLocalesOK struct {
	Payload []*models.Locale
}

func (o *GetSupportedLocalesOK) Error() string {
	return fmt.Sprintf("[GET /setup/supported_locales][%d] getSupportedLocalesOK  %+v", 200, o.Payload)
}

func (o *GetSupportedLocalesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}