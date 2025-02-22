// Code generated by go-swagger; DO NOT EDIT.

package playback

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/wiim/httpapi/models"
)

// SetPlayerCmdSwitchModeReader is a Reader for the SetPlayerCmdSwitchMode structure.
type SetPlayerCmdSwitchModeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SetPlayerCmdSwitchModeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSetPlayerCmdSwitchModeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /httpapi.asp?command=setPlayerCmd:switchmode:{mode}] setPlayerCmdSwitchMode", response, response.Code())
	}
}

// NewSetPlayerCmdSwitchModeOK creates a SetPlayerCmdSwitchModeOK with default headers values
func NewSetPlayerCmdSwitchModeOK() *SetPlayerCmdSwitchModeOK {
	return &SetPlayerCmdSwitchModeOK{}
}

/*
SetPlayerCmdSwitchModeOK describes a response with status code 200, with default header values.

Successful
*/
type SetPlayerCmdSwitchModeOK struct {
	Payload models.OperationResultPlainText
}

// IsSuccess returns true when this set player cmd switch mode o k response has a 2xx status code
func (o *SetPlayerCmdSwitchModeOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this set player cmd switch mode o k response has a 3xx status code
func (o *SetPlayerCmdSwitchModeOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this set player cmd switch mode o k response has a 4xx status code
func (o *SetPlayerCmdSwitchModeOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this set player cmd switch mode o k response has a 5xx status code
func (o *SetPlayerCmdSwitchModeOK) IsServerError() bool {
	return false
}

// IsCode returns true when this set player cmd switch mode o k response a status code equal to that given
func (o *SetPlayerCmdSwitchModeOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the set player cmd switch mode o k response
func (o *SetPlayerCmdSwitchModeOK) Code() int {
	return 200
}

func (o *SetPlayerCmdSwitchModeOK) Error() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=setPlayerCmd:switchmode:{mode}][%d] setPlayerCmdSwitchModeOK  %+v", 200, o.Payload)
}

func (o *SetPlayerCmdSwitchModeOK) String() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=setPlayerCmd:switchmode:{mode}][%d] setPlayerCmdSwitchModeOK  %+v", 200, o.Payload)
}

func (o *SetPlayerCmdSwitchModeOK) GetPayload() models.OperationResultPlainText {
	return o.Payload
}

func (o *SetPlayerCmdSwitchModeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
