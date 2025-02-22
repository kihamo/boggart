// Code generated by go-swagger; DO NOT EDIT.

package presets

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewMCUKeyShortClickParams creates a new MCUKeyShortClickParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewMCUKeyShortClickParams() *MCUKeyShortClickParams {
	return &MCUKeyShortClickParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewMCUKeyShortClickParamsWithTimeout creates a new MCUKeyShortClickParams object
// with the ability to set a timeout on a request.
func NewMCUKeyShortClickParamsWithTimeout(timeout time.Duration) *MCUKeyShortClickParams {
	return &MCUKeyShortClickParams{
		timeout: timeout,
	}
}

// NewMCUKeyShortClickParamsWithContext creates a new MCUKeyShortClickParams object
// with the ability to set a context for a request.
func NewMCUKeyShortClickParamsWithContext(ctx context.Context) *MCUKeyShortClickParams {
	return &MCUKeyShortClickParams{
		Context: ctx,
	}
}

// NewMCUKeyShortClickParamsWithHTTPClient creates a new MCUKeyShortClickParams object
// with the ability to set a custom HTTPClient for a request.
func NewMCUKeyShortClickParamsWithHTTPClient(client *http.Client) *MCUKeyShortClickParams {
	return &MCUKeyShortClickParams{
		HTTPClient: client,
	}
}

/*
MCUKeyShortClickParams contains all the parameters to send to the API endpoint

	for the m c u key short click operation.

	Typically these are written to a http.Request.
*/
type MCUKeyShortClickParams struct {

	/* Number.

	   The numeric value of the required Preset Value range is from 0 - 12

	*/
	Number int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the m c u key short click params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *MCUKeyShortClickParams) WithDefaults() *MCUKeyShortClickParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the m c u key short click params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *MCUKeyShortClickParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the m c u key short click params
func (o *MCUKeyShortClickParams) WithTimeout(timeout time.Duration) *MCUKeyShortClickParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the m c u key short click params
func (o *MCUKeyShortClickParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the m c u key short click params
func (o *MCUKeyShortClickParams) WithContext(ctx context.Context) *MCUKeyShortClickParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the m c u key short click params
func (o *MCUKeyShortClickParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the m c u key short click params
func (o *MCUKeyShortClickParams) WithHTTPClient(client *http.Client) *MCUKeyShortClickParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the m c u key short click params
func (o *MCUKeyShortClickParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithNumber adds the number to the m c u key short click params
func (o *MCUKeyShortClickParams) WithNumber(number int64) *MCUKeyShortClickParams {
	o.SetNumber(number)
	return o
}

// SetNumber adds the number to the m c u key short click params
func (o *MCUKeyShortClickParams) SetNumber(number int64) {
	o.Number = number
}

// WriteToRequest writes these params to a swagger request
func (o *MCUKeyShortClickParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param number
	if err := r.SetPathParam("number", swag.FormatInt64(o.Number)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
