// Code generated by go-swagger; DO NOT EDIT.

package requests

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

// NewGetRequestDetailsParams creates a new GetRequestDetailsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetRequestDetailsParams() *GetRequestDetailsParams {
	return &GetRequestDetailsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetRequestDetailsParamsWithTimeout creates a new GetRequestDetailsParams object
// with the ability to set a timeout on a request.
func NewGetRequestDetailsParamsWithTimeout(timeout time.Duration) *GetRequestDetailsParams {
	return &GetRequestDetailsParams{
		timeout: timeout,
	}
}

// NewGetRequestDetailsParamsWithContext creates a new GetRequestDetailsParams object
// with the ability to set a context for a request.
func NewGetRequestDetailsParamsWithContext(ctx context.Context) *GetRequestDetailsParams {
	return &GetRequestDetailsParams{
		Context: ctx,
	}
}

// NewGetRequestDetailsParamsWithHTTPClient creates a new GetRequestDetailsParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetRequestDetailsParamsWithHTTPClient(client *http.Client) *GetRequestDetailsParams {
	return &GetRequestDetailsParams{
		HTTPClient: client,
	}
}

/*
GetRequestDetailsParams contains all the parameters to send to the API endpoint

	for the get request details operation.

	Typically these are written to a http.Request.
*/
type GetRequestDetailsParams struct {

	/* ID.

	   Request ID

	   Format: uint64
	*/
	ID uint64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get request details params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetRequestDetailsParams) WithDefaults() *GetRequestDetailsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get request details params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetRequestDetailsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get request details params
func (o *GetRequestDetailsParams) WithTimeout(timeout time.Duration) *GetRequestDetailsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get request details params
func (o *GetRequestDetailsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get request details params
func (o *GetRequestDetailsParams) WithContext(ctx context.Context) *GetRequestDetailsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get request details params
func (o *GetRequestDetailsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get request details params
func (o *GetRequestDetailsParams) WithHTTPClient(client *http.Client) *GetRequestDetailsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get request details params
func (o *GetRequestDetailsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get request details params
func (o *GetRequestDetailsParams) WithID(id uint64) *GetRequestDetailsParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get request details params
func (o *GetRequestDetailsParams) SetID(id uint64) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetRequestDetailsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatUint64(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
