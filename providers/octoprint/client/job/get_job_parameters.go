// Code generated by go-swagger; DO NOT EDIT.

package job

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

// NewGetJobParams creates a new GetJobParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetJobParams() *GetJobParams {
	return &GetJobParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetJobParamsWithTimeout creates a new GetJobParams object
// with the ability to set a timeout on a request.
func NewGetJobParamsWithTimeout(timeout time.Duration) *GetJobParams {
	return &GetJobParams{
		timeout: timeout,
	}
}

// NewGetJobParamsWithContext creates a new GetJobParams object
// with the ability to set a context for a request.
func NewGetJobParamsWithContext(ctx context.Context) *GetJobParams {
	return &GetJobParams{
		Context: ctx,
	}
}

// NewGetJobParamsWithHTTPClient creates a new GetJobParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetJobParamsWithHTTPClient(client *http.Client) *GetJobParams {
	return &GetJobParams{
		HTTPClient: client,
	}
}

/* GetJobParams contains all the parameters to send to the API endpoint
   for the get job operation.

   Typically these are written to a http.Request.
*/
type GetJobParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get job params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetJobParams) WithDefaults() *GetJobParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get job params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetJobParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get job params
func (o *GetJobParams) WithTimeout(timeout time.Duration) *GetJobParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get job params
func (o *GetJobParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get job params
func (o *GetJobParams) WithContext(ctx context.Context) *GetJobParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get job params
func (o *GetJobParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get job params
func (o *GetJobParams) WithHTTPClient(client *http.Client) *GetJobParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get job params
func (o *GetJobParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetJobParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
