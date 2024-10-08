// Code generated by go-swagger; DO NOT EDIT.

package items

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// RemoveMemberReader is a Reader for the RemoveMember structure.
type RemoveMemberReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RemoveMemberReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRemoveMemberOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewRemoveMemberNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 405:
		result := NewRemoveMemberMethodNotAllowed()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewRemoveMemberOK creates a RemoveMemberOK with default headers values
func NewRemoveMemberOK() *RemoveMemberOK {
	return &RemoveMemberOK{}
}

/*RemoveMemberOK handles this case with default header values.

OK
*/
type RemoveMemberOK struct {
}

func (o *RemoveMemberOK) Error() string {
	return fmt.Sprintf("[DELETE /items/{itemName}/members/{memberItemName}][%d] removeMemberOK ", 200)
}

func (o *RemoveMemberOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewRemoveMemberNotFound creates a RemoveMemberNotFound with default headers values
func NewRemoveMemberNotFound() *RemoveMemberNotFound {
	return &RemoveMemberNotFound{}
}

/*RemoveMemberNotFound handles this case with default header values.

Item or member item not found or item is not of type group item.
*/
type RemoveMemberNotFound struct {
}

func (o *RemoveMemberNotFound) Error() string {
	return fmt.Sprintf("[DELETE /items/{itemName}/members/{memberItemName}][%d] removeMemberNotFound ", 404)
}

func (o *RemoveMemberNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewRemoveMemberMethodNotAllowed creates a RemoveMemberMethodNotAllowed with default headers values
func NewRemoveMemberMethodNotAllowed() *RemoveMemberMethodNotAllowed {
	return &RemoveMemberMethodNotAllowed{}
}

/*RemoveMemberMethodNotAllowed handles this case with default header values.

Member item is not editable.
*/
type RemoveMemberMethodNotAllowed struct {
}

func (o *RemoveMemberMethodNotAllowed) Error() string {
	return fmt.Sprintf("[DELETE /items/{itemName}/members/{memberItemName}][%d] removeMemberMethodNotAllowed ", 405)
}

func (o *RemoveMemberMethodNotAllowed) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
