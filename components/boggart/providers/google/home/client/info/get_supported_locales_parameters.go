// Code generated by go-swagger; DO NOT EDIT.

package info

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetSupportedLocalesParams creates a new GetSupportedLocalesParams object
// with the default values initialized.
func NewGetSupportedLocalesParams() *GetSupportedLocalesParams {

	return &GetSupportedLocalesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetSupportedLocalesParamsWithTimeout creates a new GetSupportedLocalesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetSupportedLocalesParamsWithTimeout(timeout time.Duration) *GetSupportedLocalesParams {

	return &GetSupportedLocalesParams{

		timeout: timeout,
	}
}

// NewGetSupportedLocalesParamsWithContext creates a new GetSupportedLocalesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetSupportedLocalesParamsWithContext(ctx context.Context) *GetSupportedLocalesParams {

	return &GetSupportedLocalesParams{

		Context: ctx,
	}
}

// NewGetSupportedLocalesParamsWithHTTPClient creates a new GetSupportedLocalesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetSupportedLocalesParamsWithHTTPClient(client *http.Client) *GetSupportedLocalesParams {

	return &GetSupportedLocalesParams{
		HTTPClient: client,
	}
}

/*GetSupportedLocalesParams contains all the parameters to send to the API endpoint
for the get supported locales operation typically these are written to a http.Request
*/
type GetSupportedLocalesParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get supported locales params
func (o *GetSupportedLocalesParams) WithTimeout(timeout time.Duration) *GetSupportedLocalesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get supported locales params
func (o *GetSupportedLocalesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get supported locales params
func (o *GetSupportedLocalesParams) WithContext(ctx context.Context) *GetSupportedLocalesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get supported locales params
func (o *GetSupportedLocalesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get supported locales params
func (o *GetSupportedLocalesParams) WithHTTPClient(client *http.Client) *GetSupportedLocalesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get supported locales params
func (o *GetSupportedLocalesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetSupportedLocalesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
