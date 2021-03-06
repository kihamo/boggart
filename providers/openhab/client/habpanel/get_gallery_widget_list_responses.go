// Code generated by go-swagger; DO NOT EDIT.

package habpanel

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// GetGalleryWidgetListReader is a Reader for the GetGalleryWidgetList structure.
type GetGalleryWidgetListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetGalleryWidgetListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetGalleryWidgetListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetGalleryWidgetListNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetGalleryWidgetListOK creates a GetGalleryWidgetListOK with default headers values
func NewGetGalleryWidgetListOK() *GetGalleryWidgetListOK {
	return &GetGalleryWidgetListOK{}
}

/*GetGalleryWidgetListOK handles this case with default header values.

OK
*/
type GetGalleryWidgetListOK struct {
	Payload string
}

func (o *GetGalleryWidgetListOK) Error() string {
	return fmt.Sprintf("[GET /habpanel/gallery/{galleryName}/widgets][%d] getGalleryWidgetListOK  %+v", 200, o.Payload)
}

func (o *GetGalleryWidgetListOK) GetPayload() string {
	return o.Payload
}

func (o *GetGalleryWidgetListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetGalleryWidgetListNotFound creates a GetGalleryWidgetListNotFound with default headers values
func NewGetGalleryWidgetListNotFound() *GetGalleryWidgetListNotFound {
	return &GetGalleryWidgetListNotFound{}
}

/*GetGalleryWidgetListNotFound handles this case with default header values.

Unknown gallery
*/
type GetGalleryWidgetListNotFound struct {
}

func (o *GetGalleryWidgetListNotFound) Error() string {
	return fmt.Sprintf("[GET /habpanel/gallery/{galleryName}/widgets][%d] getGalleryWidgetListNotFound ", 404)
}

func (o *GetGalleryWidgetListNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
