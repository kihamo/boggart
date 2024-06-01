// Code generated by go-swagger; DO NOT EDIT.

package operations

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
)

// NewGetDevicesParams creates a new GetDevicesParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetDevicesParams() *GetDevicesParams {
	return &GetDevicesParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetDevicesParamsWithTimeout creates a new GetDevicesParams object
// with the ability to set a timeout on a request.
func NewGetDevicesParamsWithTimeout(timeout time.Duration) *GetDevicesParams {
	return &GetDevicesParams{
		timeout: timeout,
	}
}

// NewGetDevicesParamsWithContext creates a new GetDevicesParams object
// with the ability to set a context for a request.
func NewGetDevicesParamsWithContext(ctx context.Context) *GetDevicesParams {
	return &GetDevicesParams{
		Context: ctx,
	}
}

// NewGetDevicesParamsWithHTTPClient creates a new GetDevicesParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetDevicesParamsWithHTTPClient(client *http.Client) *GetDevicesParams {
	return &GetDevicesParams{
		HTTPClient: client,
	}
}

/*
GetDevicesParams contains all the parameters to send to the API endpoint

	for the get devices operation.

	Typically these are written to a http.Request.
*/
type GetDevicesParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get devices params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDevicesParams) WithDefaults() *GetDevicesParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get devices params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetDevicesParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get devices params
func (o *GetDevicesParams) WithTimeout(timeout time.Duration) *GetDevicesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get devices params
func (o *GetDevicesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get devices params
func (o *GetDevicesParams) WithContext(ctx context.Context) *GetDevicesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get devices params
func (o *GetDevicesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get devices params
func (o *GetDevicesParams) WithHTTPClient(client *http.Client) *GetDevicesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get devices params
func (o *GetDevicesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetDevicesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}