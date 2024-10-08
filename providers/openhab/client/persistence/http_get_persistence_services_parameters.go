// Code generated by go-swagger; DO NOT EDIT.

package persistence

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

// NewHTTPGetPersistenceServicesParams creates a new HTTPGetPersistenceServicesParams object
// with the default values initialized.
func NewHTTPGetPersistenceServicesParams() *HTTPGetPersistenceServicesParams {
	var ()
	return &HTTPGetPersistenceServicesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewHTTPGetPersistenceServicesParamsWithTimeout creates a new HTTPGetPersistenceServicesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewHTTPGetPersistenceServicesParamsWithTimeout(timeout time.Duration) *HTTPGetPersistenceServicesParams {
	var ()
	return &HTTPGetPersistenceServicesParams{

		timeout: timeout,
	}
}

// NewHTTPGetPersistenceServicesParamsWithContext creates a new HTTPGetPersistenceServicesParams object
// with the default values initialized, and the ability to set a context for a request
func NewHTTPGetPersistenceServicesParamsWithContext(ctx context.Context) *HTTPGetPersistenceServicesParams {
	var ()
	return &HTTPGetPersistenceServicesParams{

		Context: ctx,
	}
}

// NewHTTPGetPersistenceServicesParamsWithHTTPClient creates a new HTTPGetPersistenceServicesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewHTTPGetPersistenceServicesParamsWithHTTPClient(client *http.Client) *HTTPGetPersistenceServicesParams {
	var ()
	return &HTTPGetPersistenceServicesParams{
		HTTPClient: client,
	}
}

/*HTTPGetPersistenceServicesParams contains all the parameters to send to the API endpoint
for the http get persistence services operation typically these are written to a http.Request
*/
type HTTPGetPersistenceServicesParams struct {

	/*AcceptLanguage
	  Accept-Language

	*/
	AcceptLanguage *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the http get persistence services params
func (o *HTTPGetPersistenceServicesParams) WithTimeout(timeout time.Duration) *HTTPGetPersistenceServicesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the http get persistence services params
func (o *HTTPGetPersistenceServicesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the http get persistence services params
func (o *HTTPGetPersistenceServicesParams) WithContext(ctx context.Context) *HTTPGetPersistenceServicesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the http get persistence services params
func (o *HTTPGetPersistenceServicesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the http get persistence services params
func (o *HTTPGetPersistenceServicesParams) WithHTTPClient(client *http.Client) *HTTPGetPersistenceServicesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the http get persistence services params
func (o *HTTPGetPersistenceServicesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAcceptLanguage adds the acceptLanguage to the http get persistence services params
func (o *HTTPGetPersistenceServicesParams) WithAcceptLanguage(acceptLanguage *string) *HTTPGetPersistenceServicesParams {
	o.SetAcceptLanguage(acceptLanguage)
	return o
}

// SetAcceptLanguage adds the acceptLanguage to the http get persistence services params
func (o *HTTPGetPersistenceServicesParams) SetAcceptLanguage(acceptLanguage *string) {
	o.AcceptLanguage = acceptLanguage
}

// WriteToRequest writes these params to a swagger request
func (o *HTTPGetPersistenceServicesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.AcceptLanguage != nil {

		// header param Accept-Language
		if err := r.SetHeaderParam("Accept-Language", *o.AcceptLanguage); err != nil {
			return err
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
