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

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetMonitoringTrafficStatisticsParams creates a new GetMonitoringTrafficStatisticsParams object
// with the default values initialized.
func NewGetMonitoringTrafficStatisticsParams() *GetMonitoringTrafficStatisticsParams {

	return &GetMonitoringTrafficStatisticsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetMonitoringTrafficStatisticsParamsWithTimeout creates a new GetMonitoringTrafficStatisticsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetMonitoringTrafficStatisticsParamsWithTimeout(timeout time.Duration) *GetMonitoringTrafficStatisticsParams {

	return &GetMonitoringTrafficStatisticsParams{

		timeout: timeout,
	}
}

// NewGetMonitoringTrafficStatisticsParamsWithContext creates a new GetMonitoringTrafficStatisticsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetMonitoringTrafficStatisticsParamsWithContext(ctx context.Context) *GetMonitoringTrafficStatisticsParams {

	return &GetMonitoringTrafficStatisticsParams{

		Context: ctx,
	}
}

// NewGetMonitoringTrafficStatisticsParamsWithHTTPClient creates a new GetMonitoringTrafficStatisticsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetMonitoringTrafficStatisticsParamsWithHTTPClient(client *http.Client) *GetMonitoringTrafficStatisticsParams {

	return &GetMonitoringTrafficStatisticsParams{
		HTTPClient: client,
	}
}

/*GetMonitoringTrafficStatisticsParams contains all the parameters to send to the API endpoint
for the get monitoring traffic statistics operation typically these are written to a http.Request
*/
type GetMonitoringTrafficStatisticsParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get monitoring traffic statistics params
func (o *GetMonitoringTrafficStatisticsParams) WithTimeout(timeout time.Duration) *GetMonitoringTrafficStatisticsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get monitoring traffic statistics params
func (o *GetMonitoringTrafficStatisticsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get monitoring traffic statistics params
func (o *GetMonitoringTrafficStatisticsParams) WithContext(ctx context.Context) *GetMonitoringTrafficStatisticsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get monitoring traffic statistics params
func (o *GetMonitoringTrafficStatisticsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get monitoring traffic statistics params
func (o *GetMonitoringTrafficStatisticsParams) WithHTTPClient(client *http.Client) *GetMonitoringTrafficStatisticsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get monitoring traffic statistics params
func (o *GetMonitoringTrafficStatisticsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetMonitoringTrafficStatisticsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
