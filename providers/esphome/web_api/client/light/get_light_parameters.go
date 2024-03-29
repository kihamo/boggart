// Code generated by go-swagger; DO NOT EDIT.

package light

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

// NewGetLightParams creates a new GetLightParams object
// with the default values initialized.
func NewGetLightParams() *GetLightParams {
	var ()
	return &GetLightParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetLightParamsWithTimeout creates a new GetLightParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetLightParamsWithTimeout(timeout time.Duration) *GetLightParams {
	var ()
	return &GetLightParams{

		timeout: timeout,
	}
}

// NewGetLightParamsWithContext creates a new GetLightParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetLightParamsWithContext(ctx context.Context) *GetLightParams {
	var ()
	return &GetLightParams{

		Context: ctx,
	}
}

// NewGetLightParamsWithHTTPClient creates a new GetLightParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetLightParamsWithHTTPClient(client *http.Client) *GetLightParams {
	var ()
	return &GetLightParams{
		HTTPClient: client,
	}
}

/*GetLightParams contains all the parameters to send to the API endpoint
for the get light operation typically these are written to a http.Request
*/
type GetLightParams struct {

	/*ID
	  Light ID

	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get light params
func (o *GetLightParams) WithTimeout(timeout time.Duration) *GetLightParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get light params
func (o *GetLightParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get light params
func (o *GetLightParams) WithContext(ctx context.Context) *GetLightParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get light params
func (o *GetLightParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get light params
func (o *GetLightParams) WithHTTPClient(client *http.Client) *GetLightParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get light params
func (o *GetLightParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get light params
func (o *GetLightParams) WithID(id string) *GetLightParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get light params
func (o *GetLightParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetLightParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
