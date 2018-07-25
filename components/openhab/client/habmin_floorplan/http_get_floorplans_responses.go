// Code generated by go-swagger; DO NOT EDIT.

package habmin_floorplan

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// HTTPGetFloorplansReader is a Reader for the HTTPGetFloorplans structure.
type HTTPGetFloorplansReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *HTTPGetFloorplansReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {

	result := NewHTTPGetFloorplansDefault(response.Code())
	if err := result.readResponse(response, consumer, o.formats); err != nil {
		return nil, err
	}
	if response.Code()/100 == 2 {
		return result, nil
	}
	return nil, result

}

// NewHTTPGetFloorplansDefault creates a HTTPGetFloorplansDefault with default headers values
func NewHTTPGetFloorplansDefault(code int) *HTTPGetFloorplansDefault {
	return &HTTPGetFloorplansDefault{
		_statusCode: code,
	}
}

/*HTTPGetFloorplansDefault handles this case with default header values.

successful operation
*/
type HTTPGetFloorplansDefault struct {
	_statusCode int
}

// Code gets the status code for the http get floorplans default response
func (o *HTTPGetFloorplansDefault) Code() int {
	return o._statusCode
}

func (o *HTTPGetFloorplansDefault) Error() string {
	return fmt.Sprintf("[GET /habmin/floorplan][%d] httpGetFloorplans default ", o._statusCode)
}

func (o *HTTPGetFloorplansDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}