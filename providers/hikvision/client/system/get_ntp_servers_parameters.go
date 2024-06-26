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

// NewGetNtpServersParams creates a new GetNtpServersParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetNtpServersParams() *GetNtpServersParams {
	return &GetNtpServersParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetNtpServersParamsWithTimeout creates a new GetNtpServersParams object
// with the ability to set a timeout on a request.
func NewGetNtpServersParamsWithTimeout(timeout time.Duration) *GetNtpServersParams {
	return &GetNtpServersParams{
		timeout: timeout,
	}
}

// NewGetNtpServersParamsWithContext creates a new GetNtpServersParams object
// with the ability to set a context for a request.
func NewGetNtpServersParamsWithContext(ctx context.Context) *GetNtpServersParams {
	return &GetNtpServersParams{
		Context: ctx,
	}
}

// NewGetNtpServersParamsWithHTTPClient creates a new GetNtpServersParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetNtpServersParamsWithHTTPClient(client *http.Client) *GetNtpServersParams {
	return &GetNtpServersParams{
		HTTPClient: client,
	}
}

/*
GetNtpServersParams contains all the parameters to send to the API endpoint

	for the get ntp servers operation.

	Typically these are written to a http.Request.
*/
type GetNtpServersParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get ntp servers params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetNtpServersParams) WithDefaults() *GetNtpServersParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get ntp servers params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetNtpServersParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get ntp servers params
func (o *GetNtpServersParams) WithTimeout(timeout time.Duration) *GetNtpServersParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get ntp servers params
func (o *GetNtpServersParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get ntp servers params
func (o *GetNtpServersParams) WithContext(ctx context.Context) *GetNtpServersParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get ntp servers params
func (o *GetNtpServersParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get ntp servers params
func (o *GetNtpServersParams) WithHTTPClient(client *http.Client) *GetNtpServersParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get ntp servers params
func (o *GetNtpServersParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetNtpServersParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
