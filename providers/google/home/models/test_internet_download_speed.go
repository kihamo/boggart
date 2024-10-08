// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// TestInternetDownloadSpeed test internet download speed
// swagger:model TestInternetDownloadSpeed
type TestInternetDownloadSpeed struct {

	// bytes received
	// Required: true
	BytesReceived *int64 `json:"bytes_received"`

	// response code
	// Required: true
	ResponseCode *int64 `json:"response_code"`

	// time for data fetch
	// Required: true
	TimeForDataFetch *int64 `json:"time_for_data_fetch"`

	// time for http response
	// Required: true
	TimeForHTTPResponse *int64 `json:"time_for_http_response"`
}

// Validate validates this test internet download speed
func (m *TestInternetDownloadSpeed) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBytesReceived(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResponseCode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimeForDataFetch(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTimeForHTTPResponse(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TestInternetDownloadSpeed) validateBytesReceived(formats strfmt.Registry) error {

	if err := validate.Required("bytes_received", "body", m.BytesReceived); err != nil {
		return err
	}

	return nil
}

func (m *TestInternetDownloadSpeed) validateResponseCode(formats strfmt.Registry) error {

	if err := validate.Required("response_code", "body", m.ResponseCode); err != nil {
		return err
	}

	return nil
}

func (m *TestInternetDownloadSpeed) validateTimeForDataFetch(formats strfmt.Registry) error {

	if err := validate.Required("time_for_data_fetch", "body", m.TimeForDataFetch); err != nil {
		return err
	}

	return nil
}

func (m *TestInternetDownloadSpeed) validateTimeForHTTPResponse(formats strfmt.Registry) error {

	if err := validate.Required("time_for_http_response", "body", m.TimeForHTTPResponse); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TestInternetDownloadSpeed) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TestInternetDownloadSpeed) UnmarshalBinary(b []byte) error {
	var res TestInternetDownloadSpeed
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
