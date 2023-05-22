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

// ShowIPHotspotReader is a Reader for the ShowIPHotspot structure.
type ShowIPHotspotReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ShowIPHotspotReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewShowIPHotspotOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewShowIPHotspotUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewShowIPHotspotOK creates a ShowIPHotspotOK with default headers values
func NewShowIPHotspotOK() *ShowIPHotspotOK {
	return &ShowIPHotspotOK{}
}

/* ShowIPHotspotOK describes a response with status code 200, with default header values.

Successful operation
*/
type ShowIPHotspotOK struct {
	Payload *models.ShowIPHotspotResponse
}

func (o *ShowIPHotspotOK) Error() string {
	return fmt.Sprintf("[GET /rci/show/ip/hotspot][%d] showIpHotspotOK  %+v", 200, o.Payload)
}
func (o *ShowIPHotspotOK) GetPayload() *models.ShowIPHotspotResponse {
	return o.Payload
}

func (o *ShowIPHotspotOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ShowIPHotspotResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewShowIPHotspotUnauthorized creates a ShowIPHotspotUnauthorized with default headers values
func NewShowIPHotspotUnauthorized() *ShowIPHotspotUnauthorized {
	return &ShowIPHotspotUnauthorized{}
}

/* ShowIPHotspotUnauthorized describes a response with status code 401, with default header values.

Unauthorized operation
*/
type ShowIPHotspotUnauthorized struct {
	WWWAuthenticate string
	XNDMChallenge   string
	XNDMRealm       string
}

func (o *ShowIPHotspotUnauthorized) Error() string {
	return fmt.Sprintf("[GET /rci/show/ip/hotspot][%d] showIpHotspotUnauthorized ", 401)
}

func (o *ShowIPHotspotUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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