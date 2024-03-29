// Code generated by go-swagger; DO NOT EDIT.

package plugin

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/octoprint/models"
)

// DisplayLayerProgressReader is a Reader for the DisplayLayerProgress structure.
type DisplayLayerProgressReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DisplayLayerProgressReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDisplayLayerProgressOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /plugin/DisplayLayerProgress/values] displayLayerProgress", response, response.Code())
	}
}

// NewDisplayLayerProgressOK creates a DisplayLayerProgressOK with default headers values
func NewDisplayLayerProgressOK() *DisplayLayerProgressOK {
	return &DisplayLayerProgressOK{}
}

/*
DisplayLayerProgressOK describes a response with status code 200, with default header values.

Successful operation
*/
type DisplayLayerProgressOK struct {
	Payload *models.PluginDisplayLayerProgress
}

// IsSuccess returns true when this display layer progress o k response has a 2xx status code
func (o *DisplayLayerProgressOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this display layer progress o k response has a 3xx status code
func (o *DisplayLayerProgressOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this display layer progress o k response has a 4xx status code
func (o *DisplayLayerProgressOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this display layer progress o k response has a 5xx status code
func (o *DisplayLayerProgressOK) IsServerError() bool {
	return false
}

// IsCode returns true when this display layer progress o k response a status code equal to that given
func (o *DisplayLayerProgressOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the display layer progress o k response
func (o *DisplayLayerProgressOK) Code() int {
	return 200
}

func (o *DisplayLayerProgressOK) Error() string {
	return fmt.Sprintf("[GET /plugin/DisplayLayerProgress/values][%d] displayLayerProgressOK  %+v", 200, o.Payload)
}

func (o *DisplayLayerProgressOK) String() string {
	return fmt.Sprintf("[GET /plugin/DisplayLayerProgress/values][%d] displayLayerProgressOK  %+v", 200, o.Payload)
}

func (o *DisplayLayerProgressOK) GetPayload() *models.PluginDisplayLayerProgress {
	return o.Payload
}

func (o *DisplayLayerProgressOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PluginDisplayLayerProgress)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
