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

// NewSetEngGoalParams creates a new SetEngGoalParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSetEngGoalParams() *SetEngGoalParams {
	return &SetEngGoalParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSetEngGoalParamsWithTimeout creates a new SetEngGoalParams object
// with the ability to set a timeout on a request.
func NewSetEngGoalParamsWithTimeout(timeout time.Duration) *SetEngGoalParams {
	return &SetEngGoalParams{
		timeout: timeout,
	}
}

// NewSetEngGoalParamsWithContext creates a new SetEngGoalParams object
// with the ability to set a context for a request.
func NewSetEngGoalParamsWithContext(ctx context.Context) *SetEngGoalParams {
	return &SetEngGoalParams{
		Context: ctx,
	}
}

// NewSetEngGoalParamsWithHTTPClient creates a new SetEngGoalParams object
// with the ability to set a custom HTTPClient for a request.
func NewSetEngGoalParamsWithHTTPClient(client *http.Client) *SetEngGoalParams {
	return &SetEngGoalParams{
		HTTPClient: client,
	}
}

/*
SetEngGoalParams contains all the parameters to send to the API endpoint

	for the set eng goal operation.

	Typically these are written to a http.Request.
*/
type SetEngGoalParams struct {

	// Request.
	Request SetEngGoalBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the set eng goal params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetEngGoalParams) WithDefaults() *SetEngGoalParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the set eng goal params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetEngGoalParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the set eng goal params
func (o *SetEngGoalParams) WithTimeout(timeout time.Duration) *SetEngGoalParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set eng goal params
func (o *SetEngGoalParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set eng goal params
func (o *SetEngGoalParams) WithContext(ctx context.Context) *SetEngGoalParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set eng goal params
func (o *SetEngGoalParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set eng goal params
func (o *SetEngGoalParams) WithHTTPClient(client *http.Client) *SetEngGoalParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set eng goal params
func (o *SetEngGoalParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRequest adds the request to the set eng goal params
func (o *SetEngGoalParams) WithRequest(request SetEngGoalBody) *SetEngGoalParams {
	o.SetRequest(request)
	return o
}

// SetRequest adds the request to the set eng goal params
func (o *SetEngGoalParams) SetRequest(request SetEngGoalBody) {
	o.Request = request
}

// WriteToRequest writes these params to a swagger request
func (o *SetEngGoalParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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