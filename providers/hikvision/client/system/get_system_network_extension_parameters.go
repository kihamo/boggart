// Code generated by go-swagger; DO NOT EDIT.

package system

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

// NewGetSystemNetworkExtensionParams creates a new GetSystemNetworkExtensionParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetSystemNetworkExtensionParams() *GetSystemNetworkExtensionParams {
	return &GetSystemNetworkExtensionParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetSystemNetworkExtensionParamsWithTimeout creates a new GetSystemNetworkExtensionParams object
// with the ability to set a timeout on a request.
func NewGetSystemNetworkExtensionParamsWithTimeout(timeout time.Duration) *GetSystemNetworkExtensionParams {
	return &GetSystemNetworkExtensionParams{
		timeout: timeout,
	}
}

// NewGetSystemNetworkExtensionParamsWithContext creates a new GetSystemNetworkExtensionParams object
// with the ability to set a context for a request.
func NewGetSystemNetworkExtensionParamsWithContext(ctx context.Context) *GetSystemNetworkExtensionParams {
	return &GetSystemNetworkExtensionParams{
		Context: ctx,
	}
}

// NewGetSystemNetworkExtensionParamsWithHTTPClient creates a new GetSystemNetworkExtensionParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetSystemNetworkExtensionParamsWithHTTPClient(client *http.Client) *GetSystemNetworkExtensionParams {
	return &GetSystemNetworkExtensionParams{
		HTTPClient: client,
	}
}

/* GetSystemNetworkExtensionParams contains all the parameters to send to the API endpoint
   for the get system network extension operation.

   Typically these are written to a http.Request.
*/
type GetSystemNetworkExtensionParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get system network extension params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetSystemNetworkExtensionParams) WithDefaults() *GetSystemNetworkExtensionParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get system network extension params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetSystemNetworkExtensionParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get system network extension params
func (o *GetSystemNetworkExtensionParams) WithTimeout(timeout time.Duration) *GetSystemNetworkExtensionParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get system network extension params
func (o *GetSystemNetworkExtensionParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get system network extension params
func (o *GetSystemNetworkExtensionParams) WithContext(ctx context.Context) *GetSystemNetworkExtensionParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get system network extension params
func (o *GetSystemNetworkExtensionParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get system network extension params
func (o *GetSystemNetworkExtensionParams) WithHTTPClient(client *http.Client) *GetSystemNetworkExtensionParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get system network extension params
func (o *GetSystemNetworkExtensionParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetSystemNetworkExtensionParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}