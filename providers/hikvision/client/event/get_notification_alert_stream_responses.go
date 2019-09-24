// Code generated by go-swagger; DO NOT EDIT.

package event

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// GetNotificationAlertStreamReader is a Reader for the GetNotificationAlertStream structure.
type GetNotificationAlertStreamReader struct {
	formats strfmt.Registry
	writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *GetNotificationAlertStreamReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetNotificationAlertStreamOK(o.writer)
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetNotificationAlertStreamOK creates a GetNotificationAlertStreamOK with default headers values
func NewGetNotificationAlertStreamOK(writer io.Writer) *GetNotificationAlertStreamOK {
	return &GetNotificationAlertStreamOK{
		Payload: writer,
	}
}

/*GetNotificationAlertStreamOK handles this case with default header values.

Successful operation
*/
type GetNotificationAlertStreamOK struct {
	Payload io.Writer
}

func (o *GetNotificationAlertStreamOK) Error() string {
	return fmt.Sprintf("[GET /Event/notification/alertStream][%d] getNotificationAlertStreamOK  %+v", 200, o.Payload)
}

func (o *GetNotificationAlertStreamOK) GetPayload() io.Writer {
	return o.Payload
}

func (o *GetNotificationAlertStreamOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}