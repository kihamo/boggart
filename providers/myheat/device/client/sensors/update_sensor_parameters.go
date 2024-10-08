// Code generated by go-swagger; DO NOT EDIT.

package sensors

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

	"github.com/kihamo/boggart/providers/myheat/device/models"
)

// NewUpdateSensorParams creates a new UpdateSensorParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateSensorParams() *UpdateSensorParams {
	return &UpdateSensorParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateSensorParamsWithTimeout creates a new UpdateSensorParams object
// with the ability to set a timeout on a request.
func NewUpdateSensorParamsWithTimeout(timeout time.Duration) *UpdateSensorParams {
	return &UpdateSensorParams{
		timeout: timeout,
	}
}

// NewUpdateSensorParamsWithContext creates a new UpdateSensorParams object
// with the ability to set a context for a request.
func NewUpdateSensorParamsWithContext(ctx context.Context) *UpdateSensorParams {
	return &UpdateSensorParams{
		Context: ctx,
	}
}

// NewUpdateSensorParamsWithHTTPClient creates a new UpdateSensorParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateSensorParamsWithHTTPClient(client *http.Client) *UpdateSensorParams {
	return &UpdateSensorParams{
		HTTPClient: client,
	}
}

/* UpdateSensorParams contains all the parameters to send to the API endpoint
   for the update sensor operation.

   Typically these are written to a http.Request.
*/
type UpdateSensorParams struct {

	// Request.
	Request *models.UpdateSensorRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update sensor params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateSensorParams) WithDefaults() *UpdateSensorParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update sensor params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateSensorParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update sensor params
func (o *UpdateSensorParams) WithTimeout(timeout time.Duration) *UpdateSensorParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update sensor params
func (o *UpdateSensorParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update sensor params
func (o *UpdateSensorParams) WithContext(ctx context.Context) *UpdateSensorParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update sensor params
func (o *UpdateSensorParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update sensor params
func (o *UpdateSensorParams) WithHTTPClient(client *http.Client) *UpdateSensorParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update sensor params
func (o *UpdateSensorParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRequest adds the request to the update sensor params
func (o *UpdateSensorParams) WithRequest(request *models.UpdateSensorRequest) *UpdateSensorParams {
	o.SetRequest(request)
	return o
}

// SetRequest adds the request to the update sensor params
func (o *UpdateSensorParams) SetRequest(request *models.UpdateSensorRequest) {
	o.Request = request
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateSensorParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Request != nil {
		if err := r.SetBodyParam(o.Request); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
