// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// GetAuthReader is a Reader for the GetAuth structure.
type GetAuthReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAuthReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAuthOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetAuthUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetAuthOK creates a GetAuthOK with default headers values
func NewGetAuthOK() *GetAuthOK {
	return &GetAuthOK{}
}

/* GetAuthOK describes a response with status code 200, with default header values.

Successful operation
*/
type GetAuthOK struct {
}

func (o *GetAuthOK) Error() string {
	return fmt.Sprintf("[GET /auth][%d] getAuthOK ", 200)
}

func (o *GetAuthOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetAuthUnauthorized creates a GetAuthUnauthorized with default headers values
func NewGetAuthUnauthorized() *GetAuthUnauthorized {
	return &GetAuthUnauthorized{}
}

/* GetAuthUnauthorized describes a response with status code 401, with default header values.

Unauthorized operation
*/
type GetAuthUnauthorized struct {
	WWWAuthenticate string
	XNDMChallenge   string
	XNDMRealm       string
}

func (o *GetAuthUnauthorized) Error() string {
	return fmt.Sprintf("[GET /auth][%d] getAuthUnauthorized ", 401)
}

func (o *GetAuthUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

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
