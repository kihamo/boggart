// Code generated by go-swagger; DO NOT EDIT.

package generate

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// GenerateReader is a Reader for the Generate structure.
type GenerateReader struct {
	formats strfmt.Registry
	writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *GenerateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGenerateOK(o.writer)
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGenerateBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 423:
		result := NewGenerateLocked()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGenerateDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGenerateOK creates a GenerateOK with default headers values
func NewGenerateOK(writer io.Writer) *GenerateOK {
	return &GenerateOK{
		Payload: writer,
	}
}

/*GenerateOK handles this case with default header values.

Successful operation
*/
type GenerateOK struct {
	Payload io.Writer
}

func (o *GenerateOK) Error() string {
	return fmt.Sprintf("[GET /generate][%d] generateOK  %+v", 200, o.Payload)
}

func (o *GenerateOK) GetPayload() io.Writer {
	return o.Payload
}

func (o *GenerateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGenerateBadRequest creates a GenerateBadRequest with default headers values
func NewGenerateBadRequest() *GenerateBadRequest {
	return &GenerateBadRequest{}
}

/*GenerateBadRequest handles this case with default header values.

Bad request
*/
type GenerateBadRequest struct {
	Payload string
}

func (o *GenerateBadRequest) Error() string {
	return fmt.Sprintf("[GET /generate][%d] generateBadRequest  %+v", 400, o.Payload)
}

func (o *GenerateBadRequest) GetPayload() string {
	return o.Payload
}

func (o *GenerateBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGenerateLocked creates a GenerateLocked with default headers values
func NewGenerateLocked() *GenerateLocked {
	return &GenerateLocked{}
}

/*GenerateLocked handles this case with default header values.

API key is locked, please contact Yandex support team
*/
type GenerateLocked struct {
	Payload string
}

func (o *GenerateLocked) Error() string {
	return fmt.Sprintf("[GET /generate][%d] generateLocked  %+v", 423, o.Payload)
}

func (o *GenerateLocked) GetPayload() string {
	return o.Payload
}

func (o *GenerateLocked) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGenerateDefault creates a GenerateDefault with default headers values
func NewGenerateDefault(code int) *GenerateDefault {
	return &GenerateDefault{
		_statusCode: code,
	}
}

/*GenerateDefault handles this case with default header values.

Unexpected error
*/
type GenerateDefault struct {
	_statusCode int

	Payload string
}

// Code gets the status code for the generate default response
func (o *GenerateDefault) Code() int {
	return o._statusCode
}

func (o *GenerateDefault) Error() string {
	return fmt.Sprintf("[GET /generate][%d] generate default  %+v", o._statusCode, o.Payload)
}

func (o *GenerateDefault) GetPayload() string {
	return o.Payload
}

func (o *GenerateDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
