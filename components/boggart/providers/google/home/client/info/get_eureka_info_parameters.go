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

// NewGetEurekaInfoParams creates a new GetEurekaInfoParams object
// with the default values initialized.
func NewGetEurekaInfoParams() *GetEurekaInfoParams {
	var ()
	return &GetEurekaInfoParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetEurekaInfoParamsWithTimeout creates a new GetEurekaInfoParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetEurekaInfoParamsWithTimeout(timeout time.Duration) *GetEurekaInfoParams {
	var ()
	return &GetEurekaInfoParams{

		timeout: timeout,
	}
}

// NewGetEurekaInfoParamsWithContext creates a new GetEurekaInfoParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetEurekaInfoParamsWithContext(ctx context.Context) *GetEurekaInfoParams {
	var ()
	return &GetEurekaInfoParams{

		Context: ctx,
	}
}

// NewGetEurekaInfoParamsWithHTTPClient creates a new GetEurekaInfoParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetEurekaInfoParamsWithHTTPClient(client *http.Client) *GetEurekaInfoParams {
	var ()
	return &GetEurekaInfoParams{
		HTTPClient: client,
	}
}

/*GetEurekaInfoParams contains all the parameters to send to the API endpoint
for the get eureka info operation typically these are written to a http.Request
*/
type GetEurekaInfoParams struct {

	/*Options
	  Set detail mode

	*/
	Options *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get eureka info params
func (o *GetEurekaInfoParams) WithTimeout(timeout time.Duration) *GetEurekaInfoParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get eureka info params
func (o *GetEurekaInfoParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get eureka info params
func (o *GetEurekaInfoParams) WithContext(ctx context.Context) *GetEurekaInfoParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get eureka info params
func (o *GetEurekaInfoParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get eureka info params
func (o *GetEurekaInfoParams) WithHTTPClient(client *http.Client) *GetEurekaInfoParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get eureka info params
func (o *GetEurekaInfoParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithOptions adds the options to the get eureka info params
func (o *GetEurekaInfoParams) WithOptions(options *string) *GetEurekaInfoParams {
	o.SetOptions(options)
	return o
}

// SetOptions adds the options to the get eureka info params
func (o *GetEurekaInfoParams) SetOptions(options *string) {
	o.Options = options
}

// WriteToRequest writes these params to a swagger request
func (o *GetEurekaInfoParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Options != nil {

		// query param options
		var qrOptions string
		if o.Options != nil {
			qrOptions = *o.Options
		}
		qOptions := qrOptions
		if qOptions != "" {
			if err := r.SetQueryParam("options", qOptions); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
