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

// GetCompressLogFileReader is a Reader for the GetCompressLogFile structure.
type GetCompressLogFileReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCompressLogFileReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetCompressLogFileOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetCompressLogFileDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetCompressLogFileOK creates a GetCompressLogFileOK with default headers values
func NewGetCompressLogFileOK() *GetCompressLogFileOK {
	return &GetCompressLogFileOK{}
}

/*GetCompressLogFileOK handles this case with default header values.

Successful operation
*/
type GetCompressLogFileOK struct {
	Payload *models.DeviceCompressLogFile
}

func (o *GetCompressLogFileOK) Error() string {
	return fmt.Sprintf("[GET /device/compresslogfile][%d] getCompressLogFileOK  %+v", 200, o.Payload)
}

func (o *GetCompressLogFileOK) GetPayload() *models.DeviceCompressLogFile {
	return o.Payload
}

func (o *GetCompressLogFileOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DeviceCompressLogFile)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCompressLogFileDefault creates a GetCompressLogFileDefault with default headers values
func NewGetCompressLogFileDefault(code int) *GetCompressLogFileDefault {
	return &GetCompressLogFileDefault{
		_statusCode: code,
	}
}

/*GetCompressLogFileDefault handles this case with default header values.

Unexpected error
*/
type GetCompressLogFileDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get compress log file default response
func (o *GetCompressLogFileDefault) Code() int {
	return o._statusCode
}

func (o *GetCompressLogFileDefault) Error() string {
	return fmt.Sprintf("[GET /device/compresslogfile][%d] getCompressLogFile default  %+v", o._statusCode, o.Payload)
}

func (o *GetCompressLogFileDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetCompressLogFileDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
