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

// GetSMSCountReader is a Reader for the GetSMSCount structure.
type GetSMSCountReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSMSCountReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetSMSCountOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetSMSCountOK creates a GetSMSCountOK with default headers values
func NewGetSMSCountOK() *GetSMSCountOK {
	return &GetSMSCountOK{}
}

/*GetSMSCountOK handles this case with default header values.

Successful operation
*/
type GetSMSCountOK struct {
	Payload *models.SMSCount
}

func (o *GetSMSCountOK) Error() string {
	return fmt.Sprintf("[GET /sms/sms-count][%d] getSMSCountOK  %+v", 200, o.Payload)
}

func (o *GetSMSCountOK) GetPayload() *models.SMSCount {
	return o.Payload
}

func (o *GetSMSCountOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.SMSCount)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
