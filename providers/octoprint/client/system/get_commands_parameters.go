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

// NewGetCommandsParams creates a new GetCommandsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetCommandsParams() *GetCommandsParams {
	return &GetCommandsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetCommandsParamsWithTimeout creates a new GetCommandsParams object
// with the ability to set a timeout on a request.
func NewGetCommandsParamsWithTimeout(timeout time.Duration) *GetCommandsParams {
	return &GetCommandsParams{
		timeout: timeout,
	}
}

// NewGetCommandsParamsWithContext creates a new GetCommandsParams object
// with the ability to set a context for a request.
func NewGetCommandsParamsWithContext(ctx context.Context) *GetCommandsParams {
	return &GetCommandsParams{
		Context: ctx,
	}
}

// NewGetCommandsParamsWithHTTPClient creates a new GetCommandsParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetCommandsParamsWithHTTPClient(client *http.Client) *GetCommandsParams {
	return &GetCommandsParams{
		HTTPClient: client,
	}
}

/*
GetCommandsParams contains all the parameters to send to the API endpoint

	for the get commands operation.

	Typically these are written to a http.Request.
*/
type GetCommandsParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get commands params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetCommandsParams) WithDefaults() *GetCommandsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get commands params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetCommandsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get commands params
func (o *GetCommandsParams) WithTimeout(timeout time.Duration) *GetCommandsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get commands params
func (o *GetCommandsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get commands params
func (o *GetCommandsParams) WithContext(ctx context.Context) *GetCommandsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get commands params
func (o *GetCommandsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get commands params
func (o *GetCommandsParams) WithHTTPClient(client *http.Client) *GetCommandsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get commands params
func (o *GetCommandsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetCommandsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
