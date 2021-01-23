// Code generated by go-swagger; DO NOT EDIT.

package auth

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

// NewLoginParams creates a new LoginParams object
// with the default values initialized.
func NewLoginParams() *LoginParams {
	var ()
	return &LoginParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewLoginParamsWithTimeout creates a new LoginParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewLoginParamsWithTimeout(timeout time.Duration) *LoginParams {
	var ()
	return &LoginParams{

		timeout: timeout,
	}
}

// NewLoginParamsWithContext creates a new LoginParams object
// with the default values initialized, and the ability to set a context for a request
func NewLoginParamsWithContext(ctx context.Context) *LoginParams {
	var ()
	return &LoginParams{

		Context: ctx,
	}
}

// NewLoginParamsWithHTTPClient creates a new LoginParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewLoginParamsWithHTTPClient(client *http.Client) *LoginParams {
	var ()
	return &LoginParams{
		HTTPClient: client,
	}
}

/*LoginParams contains all the parameters to send to the API endpoint
for the login operation typically these are written to a http.Request
*/
type LoginParams struct {

	/*Request*/
	Request LoginBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the login params
func (o *LoginParams) WithTimeout(timeout time.Duration) *LoginParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the login params
func (o *LoginParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the login params
func (o *LoginParams) WithContext(ctx context.Context) *LoginParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the login params
func (o *LoginParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the login params
func (o *LoginParams) WithHTTPClient(client *http.Client) *LoginParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the login params
func (o *LoginParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRequest adds the request to the login params
func (o *LoginParams) WithRequest(request LoginBody) *LoginParams {
	o.SetRequest(request)
	return o
}

// SetRequest adds the request to the login params
func (o *LoginParams) SetRequest(request LoginBody) {
	o.Request = request
}

// WriteToRequest writes these params to a swagger request
func (o *LoginParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.Request); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
