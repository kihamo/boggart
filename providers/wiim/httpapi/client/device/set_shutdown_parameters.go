// Code generated by go-swagger; DO NOT EDIT.

package device

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

// NewSetShutdownParams creates a new SetShutdownParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSetShutdownParams() *SetShutdownParams {
	return &SetShutdownParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSetShutdownParamsWithTimeout creates a new SetShutdownParams object
// with the ability to set a timeout on a request.
func NewSetShutdownParamsWithTimeout(timeout time.Duration) *SetShutdownParams {
	return &SetShutdownParams{
		timeout: timeout,
	}
}

// NewSetShutdownParamsWithContext creates a new SetShutdownParams object
// with the ability to set a context for a request.
func NewSetShutdownParamsWithContext(ctx context.Context) *SetShutdownParams {
	return &SetShutdownParams{
		Context: ctx,
	}
}

// NewSetShutdownParamsWithHTTPClient creates a new SetShutdownParams object
// with the ability to set a custom HTTPClient for a request.
func NewSetShutdownParamsWithHTTPClient(client *http.Client) *SetShutdownParams {
	return &SetShutdownParams{
		HTTPClient: client,
	}
}

/*
SetShutdownParams contains all the parameters to send to the API endpoint

	for the set shutdown operation.

	Typically these are written to a http.Request.
*/
type SetShutdownParams struct {

	/* Seconds.

	   0: shutdown immediately -1: cancel the previous shutdown timer

	*/
	Seconds int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the set shutdown params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetShutdownParams) WithDefaults() *SetShutdownParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the set shutdown params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetShutdownParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the set shutdown params
func (o *SetShutdownParams) WithTimeout(timeout time.Duration) *SetShutdownParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set shutdown params
func (o *SetShutdownParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set shutdown params
func (o *SetShutdownParams) WithContext(ctx context.Context) *SetShutdownParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set shutdown params
func (o *SetShutdownParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set shutdown params
func (o *SetShutdownParams) WithHTTPClient(client *http.Client) *SetShutdownParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set shutdown params
func (o *SetShutdownParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithSeconds adds the seconds to the set shutdown params
func (o *SetShutdownParams) WithSeconds(seconds int64) *SetShutdownParams {
	o.SetSeconds(seconds)
	return o
}

// SetSeconds adds the seconds to the set shutdown params
func (o *SetShutdownParams) SetSeconds(seconds int64) {
	o.Seconds = seconds
}

// WriteToRequest writes these params to a swagger request
func (o *SetShutdownParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param seconds
	if err := r.SetPathParam("seconds", swag.FormatInt64(o.Seconds)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
