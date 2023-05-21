// Code generated by go-swagger; DO NOT EDIT.

package show

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/keenetic/models"
)

// ShowSystemReader is a Reader for the ShowSystem structure.
type ShowSystemReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ShowSystemReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewShowSystemOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewShowSystemUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewShowSystemOK creates a ShowSystemOK with default headers values
func NewShowSystemOK() *ShowSystemOK {
	return &ShowSystemOK{}
}

/* ShowSystemOK describes a response with status code 200, with default header values.

Successful operation
*/
type ShowSystemOK struct {
	Payload *models.ShowSystemResponse
}

func (o *ShowSystemOK) Error() string {
	return fmt.Sprintf("[GET /rci/show/system][%d] showSystemOK  %+v", 200, o.Payload)
}
func (o *ShowSystemOK) GetPayload() *models.ShowSystemResponse {
	return o.Payload
}

func (o *ShowSystemOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ShowSystemResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewShowSystemUnauthorized creates a ShowSystemUnauthorized with default headers values
func NewShowSystemUnauthorized() *ShowSystemUnauthorized {
	return &ShowSystemUnauthorized{}
}

/* ShowSystemUnauthorized describes a response with status code 401, with default header values.

Unauthorized operation
*/
type ShowSystemUnauthorized struct {
	WWWAuthenticate string
	XNDMChallenge   string
	XNDMRealm       string
}

func (o *ShowSystemUnauthorized) Error() string {
	return fmt.Sprintf("[GET /rci/show/system][%d] showSystemUnauthorized ", 401)
}

func (o *ShowSystemUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header WWW-Authenticate
	hdrWWWAuthenticate := response.GetHeader("WWW-Authenticate")

	if hdrWWWAuthenticate != "" {
		o.WWWAuthenticate = hdrWWWAuthenticate
	}

	// hydrates response header X-NDM-Challenge
	hdrXNDMChallenge := response.GetHeader("X-NDM-Challenge")

	if hdrXNDMChallenge != "" {
		o.XNDMChallenge = hdrXNDMChallenge
	}

	// hydrates response header X-NDM-Realm
	hdrXNDMRealm := response.GetHeader("X-NDM-Realm")

	if hdrXNDMRealm != "" {
		o.XNDMRealm = hdrXNDMRealm
	}

	return nil
}
