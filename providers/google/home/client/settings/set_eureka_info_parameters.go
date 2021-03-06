// Code generated by go-swagger; DO NOT EDIT.

package settings

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

	models "github.com/kihamo/boggart/providers/google/home/models"
)

// NewSetEurekaInfoParams creates a new SetEurekaInfoParams object
// with the default values initialized.
func NewSetEurekaInfoParams() *SetEurekaInfoParams {
	var ()
	return &SetEurekaInfoParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewSetEurekaInfoParamsWithTimeout creates a new SetEurekaInfoParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewSetEurekaInfoParamsWithTimeout(timeout time.Duration) *SetEurekaInfoParams {
	var ()
	return &SetEurekaInfoParams{

		timeout: timeout,
	}
}

// NewSetEurekaInfoParamsWithContext creates a new SetEurekaInfoParams object
// with the default values initialized, and the ability to set a context for a request
func NewSetEurekaInfoParamsWithContext(ctx context.Context) *SetEurekaInfoParams {
	var ()
	return &SetEurekaInfoParams{

		Context: ctx,
	}
}

// NewSetEurekaInfoParamsWithHTTPClient creates a new SetEurekaInfoParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewSetEurekaInfoParamsWithHTTPClient(client *http.Client) *SetEurekaInfoParams {
	var ()
	return &SetEurekaInfoParams{
		HTTPClient: client,
	}
}

/*SetEurekaInfoParams contains all the parameters to send to the API endpoint
for the set eureka info operation typically these are written to a http.Request
*/
type SetEurekaInfoParams struct {

	/*Body
	  List fields for set

	*/
	Body *models.SetEurekaInfo

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the set eureka info params
func (o *SetEurekaInfoParams) WithTimeout(timeout time.Duration) *SetEurekaInfoParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set eureka info params
func (o *SetEurekaInfoParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set eureka info params
func (o *SetEurekaInfoParams) WithContext(ctx context.Context) *SetEurekaInfoParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set eureka info params
func (o *SetEurekaInfoParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set eureka info params
func (o *SetEurekaInfoParams) WithHTTPClient(client *http.Client) *SetEurekaInfoParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set eureka info params
func (o *SetEurekaInfoParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the set eureka info params
func (o *SetEurekaInfoParams) WithBody(body *models.SetEurekaInfo) *SetEurekaInfoParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the set eureka info params
func (o *SetEurekaInfoParams) SetBody(body *models.SetEurekaInfo) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *SetEurekaInfoParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
