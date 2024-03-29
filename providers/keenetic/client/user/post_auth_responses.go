// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PostAuthReader is a Reader for the PostAuth structure.
type PostAuthReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostAuthReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostAuthOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPostAuthOK creates a PostAuthOK with default headers values
func NewPostAuthOK() *PostAuthOK {
	return &PostAuthOK{}
}

/* PostAuthOK describes a response with status code 200, with default header values.

Successful operation
*/
type PostAuthOK struct {
}

func (o *PostAuthOK) Error() string {
	return fmt.Sprintf("[POST /auth][%d] postAuthOK ", 200)
}

func (o *PostAuthOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*PostAuthBody post auth body
swagger:model PostAuthBody
*/
type PostAuthBody struct {

	// login
	Login string `json:"login,omitempty"`

	// password
	Password string `json:"password,omitempty"`
}

// Validate validates this post auth body
func (o *PostAuthBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post auth body based on context it is used
func (o *PostAuthBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostAuthBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostAuthBody) UnmarshalBinary(b []byte) error {
	var res PostAuthBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
