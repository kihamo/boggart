// Code generated by go-swagger; DO NOT EDIT.

package state

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

// NewGetObjStateParams creates a new GetObjStateParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetObjStateParams() *GetObjStateParams {
	return &GetObjStateParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetObjStateParamsWithTimeout creates a new GetObjStateParams object
// with the ability to set a timeout on a request.
func NewGetObjStateParamsWithTimeout(timeout time.Duration) *GetObjStateParams {
	return &GetObjStateParams{
		timeout: timeout,
	}
}

// NewGetObjStateParamsWithContext creates a new GetObjStateParams object
// with the ability to set a context for a request.
func NewGetObjStateParamsWithContext(ctx context.Context) *GetObjStateParams {
	return &GetObjStateParams{
		Context: ctx,
	}
}

// NewGetObjStateParamsWithHTTPClient creates a new GetObjStateParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetObjStateParamsWithHTTPClient(client *http.Client) *GetObjStateParams {
	return &GetObjStateParams{
		HTTPClient: client,
	}
}

/* GetObjStateParams contains all the parameters to send to the API endpoint
   for the get obj state operation.

   Typically these are written to a http.Request.
*/
type GetObjStateParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get obj state params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetObjStateParams) WithDefaults() *GetObjStateParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get obj state params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetObjStateParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get obj state params
func (o *GetObjStateParams) WithTimeout(timeout time.Duration) *GetObjStateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get obj state params
func (o *GetObjStateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get obj state params
func (o *GetObjStateParams) WithContext(ctx context.Context) *GetObjStateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get obj state params
func (o *GetObjStateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get obj state params
func (o *GetObjStateParams) WithHTTPClient(client *http.Client) *GetObjStateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get obj state params
func (o *GetObjStateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetObjStateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
