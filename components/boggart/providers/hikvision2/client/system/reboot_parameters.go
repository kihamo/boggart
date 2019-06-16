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

	strfmt "github.com/go-openapi/strfmt"
)

// NewRebootParams creates a new RebootParams object
// with the default values initialized.
func NewRebootParams() *RebootParams {

	return &RebootParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRebootParamsWithTimeout creates a new RebootParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRebootParamsWithTimeout(timeout time.Duration) *RebootParams {

	return &RebootParams{

		timeout: timeout,
	}
}

// NewRebootParamsWithContext creates a new RebootParams object
// with the default values initialized, and the ability to set a context for a request
func NewRebootParamsWithContext(ctx context.Context) *RebootParams {

	return &RebootParams{

		Context: ctx,
	}
}

// NewRebootParamsWithHTTPClient creates a new RebootParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewRebootParamsWithHTTPClient(client *http.Client) *RebootParams {

	return &RebootParams{
		HTTPClient: client,
	}
}

/*RebootParams contains all the parameters to send to the API endpoint
for the reboot operation typically these are written to a http.Request
*/
type RebootParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the reboot params
func (o *RebootParams) WithTimeout(timeout time.Duration) *RebootParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the reboot params
func (o *RebootParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the reboot params
func (o *RebootParams) WithContext(ctx context.Context) *RebootParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the reboot params
func (o *RebootParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the reboot params
func (o *RebootParams) WithHTTPClient(client *http.Client) *RebootParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the reboot params
func (o *RebootParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *RebootParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
