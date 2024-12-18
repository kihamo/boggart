// Code generated by go-swagger; DO NOT EDIT.

package presets

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/wiim/httpapi/models"
)

// GetPresetInfoReader is a Reader for the GetPresetInfo structure.
type GetPresetInfoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetPresetInfoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetPresetInfoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /httpapi.asp?command=getPresetInfo] getPresetInfo", response, response.Code())
	}
}

// NewGetPresetInfoOK creates a GetPresetInfoOK with default headers values
func NewGetPresetInfoOK() *GetPresetInfoOK {
	return &GetPresetInfoOK{}
}

/*
GetPresetInfoOK describes a response with status code 200, with default header values.

Successful
*/
type GetPresetInfoOK struct {
	Payload *models.PresetInfo
}

// IsSuccess returns true when this get preset info o k response has a 2xx status code
func (o *GetPresetInfoOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get preset info o k response has a 3xx status code
func (o *GetPresetInfoOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get preset info o k response has a 4xx status code
func (o *GetPresetInfoOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get preset info o k response has a 5xx status code
func (o *GetPresetInfoOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get preset info o k response a status code equal to that given
func (o *GetPresetInfoOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get preset info o k response
func (o *GetPresetInfoOK) Code() int {
	return 200
}

func (o *GetPresetInfoOK) Error() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=getPresetInfo][%d] getPresetInfoOK  %+v", 200, o.Payload)
}

func (o *GetPresetInfoOK) String() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=getPresetInfo][%d] getPresetInfoOK  %+v", 200, o.Payload)
}

func (o *GetPresetInfoOK) GetPayload() *models.PresetInfo {
	return o.Payload
}

func (o *GetPresetInfoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.PresetInfo)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
