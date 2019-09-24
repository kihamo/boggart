// Code generated by go-swagger; DO NOT EDIT.

package ptz

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

// NewGetPtzChannelsParams creates a new GetPtzChannelsParams object
// with the default values initialized.
func NewGetPtzChannelsParams() *GetPtzChannelsParams {

	return &GetPtzChannelsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetPtzChannelsParamsWithTimeout creates a new GetPtzChannelsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetPtzChannelsParamsWithTimeout(timeout time.Duration) *GetPtzChannelsParams {

	return &GetPtzChannelsParams{

		timeout: timeout,
	}
}

// NewGetPtzChannelsParamsWithContext creates a new GetPtzChannelsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetPtzChannelsParamsWithContext(ctx context.Context) *GetPtzChannelsParams {

	return &GetPtzChannelsParams{

		Context: ctx,
	}
}

// NewGetPtzChannelsParamsWithHTTPClient creates a new GetPtzChannelsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetPtzChannelsParamsWithHTTPClient(client *http.Client) *GetPtzChannelsParams {

	return &GetPtzChannelsParams{
		HTTPClient: client,
	}
}

/*GetPtzChannelsParams contains all the parameters to send to the API endpoint
for the get ptz channels operation typically these are written to a http.Request
*/
type GetPtzChannelsParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get ptz channels params
func (o *GetPtzChannelsParams) WithTimeout(timeout time.Duration) *GetPtzChannelsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get ptz channels params
func (o *GetPtzChannelsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get ptz channels params
func (o *GetPtzChannelsParams) WithContext(ctx context.Context) *GetPtzChannelsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get ptz channels params
func (o *GetPtzChannelsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get ptz channels params
func (o *GetPtzChannelsParams) WithHTTPClient(client *http.Client) *GetPtzChannelsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get ptz channels params
func (o *GetPtzChannelsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetPtzChannelsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}