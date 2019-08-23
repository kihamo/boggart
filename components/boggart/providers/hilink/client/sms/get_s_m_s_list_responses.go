// Code generated by go-swagger; DO NOT EDIT.

package sms

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/kihamo/boggart/components/boggart/providers/hilink/models"
)

// GetSMSListReader is a Reader for the GetSMSList structure.
type GetSMSListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSMSListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetSMSListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetSMSListOK creates a GetSMSListOK with default headers values
func NewGetSMSListOK() *GetSMSListOK {
	return &GetSMSListOK{}
}

/*GetSMSListOK handles this case with default header values.

Successful operation
*/
type GetSMSListOK struct {
	Payload *models.SMSList
}

func (o *GetSMSListOK) Error() string {
	return fmt.Sprintf("[POST /sms/sms-list][%d] getSMSListOK  %+v", 200, o.Payload)
}

func (o *GetSMSListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.SMSList)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}