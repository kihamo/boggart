// Code generated by go-swagger; DO NOT EDIT.

package streaming

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// GetStreamingPictureReader is a Reader for the GetStreamingPicture structure.
type GetStreamingPictureReader struct {
	formats strfmt.Registry
	writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *GetStreamingPictureReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewGetStreamingPictureOK(o.writer)
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetStreamingPictureOK creates a GetStreamingPictureOK with default headers values
func NewGetStreamingPictureOK(writer io.Writer) *GetStreamingPictureOK {
	return &GetStreamingPictureOK{
		Payload: writer,
	}
}

/*GetStreamingPictureOK handles this case with default header values.

Successful operation
*/
type GetStreamingPictureOK struct {
	Payload io.Writer
}

func (o *GetStreamingPictureOK) Error() string {
	return fmt.Sprintf("[GET /Streaming/channels/{channel}/picture][%d] getStreamingPictureOK  %+v", 200, o.Payload)
}

func (o *GetStreamingPictureOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}