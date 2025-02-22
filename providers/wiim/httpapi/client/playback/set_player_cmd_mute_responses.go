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

// SetPlayerCmdMuteReader is a Reader for the SetPlayerCmdMute structure.
type SetPlayerCmdMuteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SetPlayerCmdMuteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewSetPlayerCmdMuteOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /httpapi.asp?command=setPlayerCmd:mute:{mute}] setPlayerCmdMute", response, response.Code())
	}
}

// NewSetPlayerCmdMuteOK creates a SetPlayerCmdMuteOK with default headers values
func NewSetPlayerCmdMuteOK() *SetPlayerCmdMuteOK {
	return &SetPlayerCmdMuteOK{}
}

/*
SetPlayerCmdMuteOK describes a response with status code 200, with default header values.

Successful
*/
type SetPlayerCmdMuteOK struct {
	Payload models.OperationResultPlainText
}

// IsSuccess returns true when this set player cmd mute o k response has a 2xx status code
func (o *SetPlayerCmdMuteOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this set player cmd mute o k response has a 3xx status code
func (o *SetPlayerCmdMuteOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this set player cmd mute o k response has a 4xx status code
func (o *SetPlayerCmdMuteOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this set player cmd mute o k response has a 5xx status code
func (o *SetPlayerCmdMuteOK) IsServerError() bool {
	return false
}

// IsCode returns true when this set player cmd mute o k response a status code equal to that given
func (o *SetPlayerCmdMuteOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the set player cmd mute o k response
func (o *SetPlayerCmdMuteOK) Code() int {
	return 200
}

func (o *SetPlayerCmdMuteOK) Error() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=setPlayerCmd:mute:{mute}][%d] setPlayerCmdMuteOK  %+v", 200, o.Payload)
}

func (o *SetPlayerCmdMuteOK) String() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=setPlayerCmd:mute:{mute}][%d] setPlayerCmdMuteOK  %+v", 200, o.Payload)
}

func (o *SetPlayerCmdMuteOK) GetPayload() models.OperationResultPlainText {
	return o.Payload
}

func (o *SetPlayerCmdMuteOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
