// Code generated by go-swagger; DO NOT EDIT.

package device

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/kihamo/boggart/components/boggart/providers/hilink/models"
)

// GetDeviceBasicInformationReader is a Reader for the GetDeviceBasicInformation structure.
type GetDeviceBasicInformationReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetDeviceBasicInformationReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetDeviceBasicInformationOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetDeviceBasicInformationOK creates a GetDeviceBasicInformationOK with default headers values
func NewGetDeviceBasicInformationOK() *GetDeviceBasicInformationOK {
	return &GetDeviceBasicInformationOK{}
}

/*GetDeviceBasicInformationOK handles this case with default header values.

Successful operation
*/
type GetDeviceBasicInformationOK struct {
	Payload *models.DeviceBasicInformation
}

func (o *GetDeviceBasicInformationOK) Error() string {
	return fmt.Sprintf("[GET /device/basic_information][%d] getDeviceBasicInformationOK  %+v", 200, o.Payload)
}

func (o *GetDeviceBasicInformationOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DeviceBasicInformation)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}