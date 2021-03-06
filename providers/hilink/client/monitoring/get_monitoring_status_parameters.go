// Code generated by go-swagger; DO NOT EDIT.

package monitoring

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

// NewGetMonitoringStatusParams creates a new GetMonitoringStatusParams object
// with the default values initialized.
func NewGetMonitoringStatusParams() *GetMonitoringStatusParams {

	return &GetMonitoringStatusParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetMonitoringStatusParamsWithTimeout creates a new GetMonitoringStatusParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetMonitoringStatusParamsWithTimeout(timeout time.Duration) *GetMonitoringStatusParams {

	return &GetMonitoringStatusParams{

		timeout: timeout,
	}
}

// NewGetMonitoringStatusParamsWithContext creates a new GetMonitoringStatusParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetMonitoringStatusParamsWithContext(ctx context.Context) *GetMonitoringStatusParams {

	return &GetMonitoringStatusParams{

		Context: ctx,
	}
}

// NewGetMonitoringStatusParamsWithHTTPClient creates a new GetMonitoringStatusParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetMonitoringStatusParamsWithHTTPClient(client *http.Client) *GetMonitoringStatusParams {

	return &GetMonitoringStatusParams{
		HTTPClient: client,
	}
}

/*GetMonitoringStatusParams contains all the parameters to send to the API endpoint
for the get monitoring status operation typically these are written to a http.Request
*/
type GetMonitoringStatusParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get monitoring status params
func (o *GetMonitoringStatusParams) WithTimeout(timeout time.Duration) *GetMonitoringStatusParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get monitoring status params
func (o *GetMonitoringStatusParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get monitoring status params
func (o *GetMonitoringStatusParams) WithContext(ctx context.Context) *GetMonitoringStatusParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get monitoring status params
func (o *GetMonitoringStatusParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get monitoring status params
func (o *GetMonitoringStatusParams) WithHTTPClient(client *http.Client) *GetMonitoringStatusParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get monitoring status params
func (o *GetMonitoringStatusParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetMonitoringStatusParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
