// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/kihamo/boggart/providers/dom24/models"
)

// UserInfoReader is a Reader for the UserInfo structure.
type UserInfoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UserInfoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUserInfoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewUserInfoUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewUserInfoDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUserInfoOK creates a UserInfoOK with default headers values
func NewUserInfoOK() *UserInfoOK {
	return &UserInfoOK{}
}

/*UserInfoOK handles this case with default header values.

Successful operation
*/
type UserInfoOK struct {
	Payload *models.UserInfo
}

func (o *UserInfoOK) Error() string {
	return fmt.Sprintf("[GET /User/Info][%d] userInfoOK  %+v", 200, o.Payload)
}

func (o *UserInfoOK) GetPayload() *models.UserInfo {
	return o.Payload
}

func (o *UserInfoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.UserInfo)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUserInfoUnauthorized creates a UserInfoUnauthorized with default headers values
func NewUserInfoUnauthorized() *UserInfoUnauthorized {
	return &UserInfoUnauthorized{}
}

/*UserInfoUnauthorized handles this case with default header values.

Unauthorized
*/
type UserInfoUnauthorized struct {
	Payload *models.Error
}

func (o *UserInfoUnauthorized) Error() string {
	return fmt.Sprintf("[GET /User/Info][%d] userInfoUnauthorized  %+v", 401, o.Payload)
}

func (o *UserInfoUnauthorized) GetPayload() *models.Error {
	return o.Payload
}

func (o *UserInfoUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUserInfoDefault creates a UserInfoDefault with default headers values
func NewUserInfoDefault(code int) *UserInfoDefault {
	return &UserInfoDefault{
		_statusCode: code,
	}
}

/*UserInfoDefault handles this case with default header values.

Unexpected error
*/
type UserInfoDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the user info default response
func (o *UserInfoDefault) Code() int {
	return o._statusCode
}

func (o *UserInfoDefault) Error() string {
	return fmt.Sprintf("[GET /User/Info][%d] userInfo default  %+v", o._statusCode, o.Payload)
}

func (o *UserInfoDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *UserInfoDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
