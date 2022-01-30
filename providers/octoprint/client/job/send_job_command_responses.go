// Code generated by go-swagger; DO NOT EDIT.

package job

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SendJobCommandReader is a Reader for the SendJobCommand structure.
type SendJobCommandReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SendJobCommandReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewSendJobCommandNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 409:
		result := NewSendJobCommandConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewSendJobCommandNoContent creates a SendJobCommandNoContent with default headers values
func NewSendJobCommandNoContent() *SendJobCommandNoContent {
	return &SendJobCommandNoContent{}
}

/*SendJobCommandNoContent handles this case with default header values.

Successful operation
*/
type SendJobCommandNoContent struct {
}

func (o *SendJobCommandNoContent) Error() string {
	return fmt.Sprintf("[POST /api/job][%d] sendJobCommandNoContent ", 204)
}

func (o *SendJobCommandNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewSendJobCommandConflict creates a SendJobCommandConflict with default headers values
func NewSendJobCommandConflict() *SendJobCommandConflict {
	return &SendJobCommandConflict{}
}

/*SendJobCommandConflict handles this case with default header values.

If the printer is not operational or the current print job state does not match the preconditions for the command
*/
type SendJobCommandConflict struct {
}

func (o *SendJobCommandConflict) Error() string {
	return fmt.Sprintf("[POST /api/job][%d] sendJobCommandConflict ", 409)
}

func (o *SendJobCommandConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*SendJobCommandBody send job command body
swagger:model SendJobCommandBody
*/
type SendJobCommandBody struct {

	// action
	// Enum: [pause resume toggle]
	Action string `json:"action,omitempty"`

	// command
	// Enum: [start cancel restart pause resume toggle]
	Command string `json:"command,omitempty"`
}

// Validate validates this send job command body
func (o *SendJobCommandBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAction(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateCommand(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var sendJobCommandBodyTypeActionPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["pause","resume","toggle"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		sendJobCommandBodyTypeActionPropEnum = append(sendJobCommandBodyTypeActionPropEnum, v)
	}
}

const (

	// SendJobCommandBodyActionPause captures enum value "pause"
	SendJobCommandBodyActionPause string = "pause"

	// SendJobCommandBodyActionResume captures enum value "resume"
	SendJobCommandBodyActionResume string = "resume"

	// SendJobCommandBodyActionToggle captures enum value "toggle"
	SendJobCommandBodyActionToggle string = "toggle"
)

// prop value enum
func (o *SendJobCommandBody) validateActionEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, sendJobCommandBodyTypeActionPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (o *SendJobCommandBody) validateAction(formats strfmt.Registry) error {

	if swag.IsZero(o.Action) { // not required
		return nil
	}

	// value enum
	if err := o.validateActionEnum("body"+"."+"action", "body", o.Action); err != nil {
		return err
	}

	return nil
}

var sendJobCommandBodyTypeCommandPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["start","cancel","restart","pause","resume","toggle"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		sendJobCommandBodyTypeCommandPropEnum = append(sendJobCommandBodyTypeCommandPropEnum, v)
	}
}

const (

	// SendJobCommandBodyCommandStart captures enum value "start"
	SendJobCommandBodyCommandStart string = "start"

	// SendJobCommandBodyCommandCancel captures enum value "cancel"
	SendJobCommandBodyCommandCancel string = "cancel"

	// SendJobCommandBodyCommandRestart captures enum value "restart"
	SendJobCommandBodyCommandRestart string = "restart"

	// SendJobCommandBodyCommandPause captures enum value "pause"
	SendJobCommandBodyCommandPause string = "pause"

	// SendJobCommandBodyCommandResume captures enum value "resume"
	SendJobCommandBodyCommandResume string = "resume"

	// SendJobCommandBodyCommandToggle captures enum value "toggle"
	SendJobCommandBodyCommandToggle string = "toggle"
)

// prop value enum
func (o *SendJobCommandBody) validateCommandEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, sendJobCommandBodyTypeCommandPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (o *SendJobCommandBody) validateCommand(formats strfmt.Registry) error {

	if swag.IsZero(o.Command) { // not required
		return nil
	}

	// value enum
	if err := o.validateCommandEnum("body"+"."+"command", "body", o.Command); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SendJobCommandBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SendJobCommandBody) UnmarshalBinary(b []byte) error {
	var res SendJobCommandBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
