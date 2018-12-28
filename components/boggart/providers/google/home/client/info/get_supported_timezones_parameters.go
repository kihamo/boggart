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

// NewGetSupportedTimezonesParams creates a new GetSupportedTimezonesParams object
// with the default values initialized.
func NewGetSupportedTimezonesParams() *GetSupportedTimezonesParams {

	return &GetSupportedTimezonesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetSupportedTimezonesParamsWithTimeout creates a new GetSupportedTimezonesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetSupportedTimezonesParamsWithTimeout(timeout time.Duration) *GetSupportedTimezonesParams {

	return &GetSupportedTimezonesParams{

		timeout: timeout,
	}
}

// NewGetSupportedTimezonesParamsWithContext creates a new GetSupportedTimezonesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetSupportedTimezonesParamsWithContext(ctx context.Context) *GetSupportedTimezonesParams {

	return &GetSupportedTimezonesParams{

		Context: ctx,
	}
}

// NewGetSupportedTimezonesParamsWithHTTPClient creates a new GetSupportedTimezonesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetSupportedTimezonesParamsWithHTTPClient(client *http.Client) *GetSupportedTimezonesParams {

	return &GetSupportedTimezonesParams{
		HTTPClient: client,
	}
}

/*GetSupportedTimezonesParams contains all the parameters to send to the API endpoint
for the get supported timezones operation typically these are written to a http.Request
*/
type GetSupportedTimezonesParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get supported timezones params
func (o *GetSupportedTimezonesParams) WithTimeout(timeout time.Duration) *GetSupportedTimezonesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get supported timezones params
func (o *GetSupportedTimezonesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get supported timezones params
func (o *GetSupportedTimezonesParams) WithContext(ctx context.Context) *GetSupportedTimezonesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get supported timezones params
func (o *GetSupportedTimezonesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get supported timezones params
func (o *GetSupportedTimezonesParams) WithHTTPClient(client *http.Client) *GetSupportedTimezonesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get supported timezones params
func (o *GetSupportedTimezonesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetSupportedTimezonesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
