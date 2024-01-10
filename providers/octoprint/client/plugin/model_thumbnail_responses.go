// Code generated by go-swagger; DO NOT EDIT.

package plugin

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// ModelThumbnailReader is a Reader for the ModelThumbnail structure.
type ModelThumbnailReader struct {
	formats strfmt.Registry
	writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *ModelThumbnailReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewModelThumbnailOK(o.writer)
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /plugin/UltimakerFormatPackage/thumbnail/{name}.png] modelThumbnail", response, response.Code())
	}
}

// NewModelThumbnailOK creates a ModelThumbnailOK with default headers values
func NewModelThumbnailOK(writer io.Writer) *ModelThumbnailOK {
	return &ModelThumbnailOK{

		Payload: writer,
	}
}

/*
ModelThumbnailOK describes a response with status code 200, with default header values.

Successful operation
*/
type ModelThumbnailOK struct {
	Payload io.Writer
}

// IsSuccess returns true when this model thumbnail o k response has a 2xx status code
func (o *ModelThumbnailOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this model thumbnail o k response has a 3xx status code
func (o *ModelThumbnailOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this model thumbnail o k response has a 4xx status code
func (o *ModelThumbnailOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this model thumbnail o k response has a 5xx status code
func (o *ModelThumbnailOK) IsServerError() bool {
	return false
}

// IsCode returns true when this model thumbnail o k response a status code equal to that given
func (o *ModelThumbnailOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the model thumbnail o k response
func (o *ModelThumbnailOK) Code() int {
	return 200
}

func (o *ModelThumbnailOK) Error() string {
	return fmt.Sprintf("[GET /plugin/UltimakerFormatPackage/thumbnail/{name}.png][%d] modelThumbnailOK  %+v", 200, o.Payload)
}

func (o *ModelThumbnailOK) String() string {
	return fmt.Sprintf("[GET /plugin/UltimakerFormatPackage/thumbnail/{name}.png][%d] modelThumbnailOK  %+v", 200, o.Payload)
}

func (o *ModelThumbnailOK) GetPayload() io.Writer {
	return o.Payload
}

func (o *ModelThumbnailOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
