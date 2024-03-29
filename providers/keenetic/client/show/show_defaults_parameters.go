// Code generated by go-swagger; DO NOT EDIT.

package show

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

// NewShowDefaultsParams creates a new ShowDefaultsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewShowDefaultsParams() *ShowDefaultsParams {
	return &ShowDefaultsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewShowDefaultsParamsWithTimeout creates a new ShowDefaultsParams object
// with the ability to set a timeout on a request.
func NewShowDefaultsParamsWithTimeout(timeout time.Duration) *ShowDefaultsParams {
	return &ShowDefaultsParams{
		timeout: timeout,
	}
}

// NewShowDefaultsParamsWithContext creates a new ShowDefaultsParams object
// with the ability to set a context for a request.
func NewShowDefaultsParamsWithContext(ctx context.Context) *ShowDefaultsParams {
	return &ShowDefaultsParams{
		Context: ctx,
	}
}

// NewShowDefaultsParamsWithHTTPClient creates a new ShowDefaultsParams object
// with the ability to set a custom HTTPClient for a request.
func NewShowDefaultsParamsWithHTTPClient(client *http.Client) *ShowDefaultsParams {
	return &ShowDefaultsParams{
		HTTPClient: client,
	}
}

/* ShowDefaultsParams contains all the parameters to send to the API endpoint
   for the show defaults operation.

   Typically these are written to a http.Request.
*/
type ShowDefaultsParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the show defaults params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ShowDefaultsParams) WithDefaults() *ShowDefaultsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the show defaults params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ShowDefaultsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the show defaults params
func (o *ShowDefaultsParams) WithTimeout(timeout time.Duration) *ShowDefaultsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the show defaults params
func (o *ShowDefaultsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the show defaults params
func (o *ShowDefaultsParams) WithContext(ctx context.Context) *ShowDefaultsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the show defaults params
func (o *ShowDefaultsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the show defaults params
func (o *ShowDefaultsParams) WithHTTPClient(client *http.Client) *ShowDefaultsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the show defaults params
func (o *ShowDefaultsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *ShowDefaultsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
