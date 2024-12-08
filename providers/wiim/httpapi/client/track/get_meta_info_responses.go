// Code generated by go-swagger; DO NOT EDIT.

package track

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/kihamo/boggart/providers/wiim/httpapi/models"
)

// GetMetaInfoReader is a Reader for the GetMetaInfo structure.
type GetMetaInfoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetMetaInfoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetMetaInfoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /httpapi.asp?command=getMetaInfo] getMetaInfo", response, response.Code())
	}
}

// NewGetMetaInfoOK creates a GetMetaInfoOK with default headers values
func NewGetMetaInfoOK() *GetMetaInfoOK {
	return &GetMetaInfoOK{}
}

/*
GetMetaInfoOK describes a response with status code 200, with default header values.

Successful
*/
type GetMetaInfoOK struct {
	Payload *models.MetaInfo
}

// IsSuccess returns true when this get meta info o k response has a 2xx status code
func (o *GetMetaInfoOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get meta info o k response has a 3xx status code
func (o *GetMetaInfoOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get meta info o k response has a 4xx status code
func (o *GetMetaInfoOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get meta info o k response has a 5xx status code
func (o *GetMetaInfoOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get meta info o k response a status code equal to that given
func (o *GetMetaInfoOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get meta info o k response
func (o *GetMetaInfoOK) Code() int {
	return 200
}

func (o *GetMetaInfoOK) Error() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=getMetaInfo][%d] getMetaInfoOK  %+v", 200, o.Payload)
}

func (o *GetMetaInfoOK) String() string {
	return fmt.Sprintf("[GET /httpapi.asp?command=getMetaInfo][%d] getMetaInfoOK  %+v", 200, o.Payload)
}

func (o *GetMetaInfoOK) GetPayload() *models.MetaInfo {
	return o.Payload
}

func (o *GetMetaInfoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.MetaInfo)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}