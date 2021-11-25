// Code generated by go-swagger; DO NOT EDIT.

package feed

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

// NewGetByCityParams creates a new GetByCityParams object
// with the default values initialized.
func NewGetByCityParams() *GetByCityParams {
	var ()
	return &GetByCityParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetByCityParamsWithTimeout creates a new GetByCityParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetByCityParamsWithTimeout(timeout time.Duration) *GetByCityParams {
	var ()
	return &GetByCityParams{

		timeout: timeout,
	}
}

// NewGetByCityParamsWithContext creates a new GetByCityParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetByCityParamsWithContext(ctx context.Context) *GetByCityParams {
	var ()
	return &GetByCityParams{

		Context: ctx,
	}
}

// NewGetByCityParamsWithHTTPClient creates a new GetByCityParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetByCityParamsWithHTTPClient(client *http.Client) *GetByCityParams {
	var ()
	return &GetByCityParams{
		HTTPClient: client,
	}
}

/*GetByCityParams contains all the parameters to send to the API endpoint
for the get by city operation typically these are written to a http.Request
*/
type GetByCityParams struct {

	/*City
	  Name of the city (eg beijing), or id (eg @7397)

	*/
	City string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get by city params
func (o *GetByCityParams) WithTimeout(timeout time.Duration) *GetByCityParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get by city params
func (o *GetByCityParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get by city params
func (o *GetByCityParams) WithContext(ctx context.Context) *GetByCityParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get by city params
func (o *GetByCityParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get by city params
func (o *GetByCityParams) WithHTTPClient(client *http.Client) *GetByCityParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get by city params
func (o *GetByCityParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCity adds the city to the get by city params
func (o *GetByCityParams) WithCity(city string) *GetByCityParams {
	o.SetCity(city)
	return o
}

// SetCity adds the city to the get by city params
func (o *GetByCityParams) SetCity(city string) {
	o.City = city
}

// WriteToRequest writes these params to a swagger request
func (o *GetByCityParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param city
	if err := r.SetPathParam("city", o.City); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
