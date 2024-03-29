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

// NewSetSecurityModeParams creates a new SetSecurityModeParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSetSecurityModeParams() *SetSecurityModeParams {
	return &SetSecurityModeParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSetSecurityModeParamsWithTimeout creates a new SetSecurityModeParams object
// with the ability to set a timeout on a request.
func NewSetSecurityModeParamsWithTimeout(timeout time.Duration) *SetSecurityModeParams {
	return &SetSecurityModeParams{
		timeout: timeout,
	}
}

// NewSetSecurityModeParamsWithContext creates a new SetSecurityModeParams object
// with the ability to set a context for a request.
func NewSetSecurityModeParamsWithContext(ctx context.Context) *SetSecurityModeParams {
	return &SetSecurityModeParams{
		Context: ctx,
	}
}

// NewSetSecurityModeParamsWithHTTPClient creates a new SetSecurityModeParams object
// with the ability to set a custom HTTPClient for a request.
func NewSetSecurityModeParamsWithHTTPClient(client *http.Client) *SetSecurityModeParams {
	return &SetSecurityModeParams{
		HTTPClient: client,
	}
}

/*
SetSecurityModeParams contains all the parameters to send to the API endpoint

	for the set security mode operation.

	Typically these are written to a http.Request.
*/
type SetSecurityModeParams struct {

	// Request.
	Request SetSecurityModeBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the set security mode params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetSecurityModeParams) WithDefaults() *SetSecurityModeParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the set security mode params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetSecurityModeParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the set security mode params
func (o *SetSecurityModeParams) WithTimeout(timeout time.Duration) *SetSecurityModeParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set security mode params
func (o *SetSecurityModeParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set security mode params
func (o *SetSecurityModeParams) WithContext(ctx context.Context) *SetSecurityModeParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set security mode params
func (o *SetSecurityModeParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set security mode params
func (o *SetSecurityModeParams) WithHTTPClient(client *http.Client) *SetSecurityModeParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set security mode params
func (o *SetSecurityModeParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRequest adds the request to the set security mode params
func (o *SetSecurityModeParams) WithRequest(request SetSecurityModeBody) *SetSecurityModeParams {
	o.SetRequest(request)
	return o
}

// SetRequest adds the request to the set security mode params
func (o *SetSecurityModeParams) SetRequest(request SetSecurityModeBody) {
	o.Request = request
}

// WriteToRequest writes these params to a swagger request
func (o *SetSecurityModeParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if err := r.SetBodyParam(o.Request); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
