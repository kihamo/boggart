// Code generated by go-swagger; DO NOT EDIT.

package web_server

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/kihamo/boggart/components/boggart/providers/hilink/models"
)

// GetWebServerSessionReader is a Reader for the GetWebServerSession structure.
type GetWebServerSessionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetWebServerSessionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetWebServerSessionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetWebServerSessionOK creates a GetWebServerSessionOK with default headers values
func NewGetWebServerSessionOK() *GetWebServerSessionOK {
	return &GetWebServerSessionOK{}
}

/*GetWebServerSessionOK handles this case with default header values.

Successful operation
*/
type GetWebServerSessionOK struct {
	Payload *models.SessionToken
}

func (o *GetWebServerSessionOK) Error() string {
	return fmt.Sprintf("[GET /webserver/SesTokInfo][%d] getWebServerSessionOK  %+v", 200, o.Payload)
}

func (o *GetWebServerSessionOK) GetPayload() *models.SessionToken {
	return o.Payload
}

func (o *GetWebServerSessionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.SessionToken)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
