// Code generated by go-swagger; DO NOT EDIT.

package info

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

// NewPostTestInternetDownloadSpeedParams creates a new PostTestInternetDownloadSpeedParams object
// with the default values initialized.
func NewPostTestInternetDownloadSpeedParams() *PostTestInternetDownloadSpeedParams {
	var ()
	return &PostTestInternetDownloadSpeedParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostTestInternetDownloadSpeedParamsWithTimeout creates a new PostTestInternetDownloadSpeedParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostTestInternetDownloadSpeedParamsWithTimeout(timeout time.Duration) *PostTestInternetDownloadSpeedParams {
	var ()
	return &PostTestInternetDownloadSpeedParams{

		timeout: timeout,
	}
}

// NewPostTestInternetDownloadSpeedParamsWithContext creates a new PostTestInternetDownloadSpeedParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostTestInternetDownloadSpeedParamsWithContext(ctx context.Context) *PostTestInternetDownloadSpeedParams {
	var ()
	return &PostTestInternetDownloadSpeedParams{

		Context: ctx,
	}
}

// NewPostTestInternetDownloadSpeedParamsWithHTTPClient creates a new PostTestInternetDownloadSpeedParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostTestInternetDownloadSpeedParamsWithHTTPClient(client *http.Client) *PostTestInternetDownloadSpeedParams {
	var ()
	return &PostTestInternetDownloadSpeedParams{
		HTTPClient: client,
	}
}

/*PostTestInternetDownloadSpeedParams contains all the parameters to send to the API endpoint
for the post test internet download speed operation typically these are written to a http.Request
*/
type PostTestInternetDownloadSpeedParams struct {

	/*Body
	  URL for test

	*/
	Body PostTestInternetDownloadSpeedBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post test internet download speed params
func (o *PostTestInternetDownloadSpeedParams) WithTimeout(timeout time.Duration) *PostTestInternetDownloadSpeedParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post test internet download speed params
func (o *PostTestInternetDownloadSpeedParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post test internet download speed params
func (o *PostTestInternetDownloadSpeedParams) WithContext(ctx context.Context) *PostTestInternetDownloadSpeedParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post test internet download speed params
func (o *PostTestInternetDownloadSpeedParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post test internet download speed params
func (o *PostTestInternetDownloadSpeedParams) WithHTTPClient(client *http.Client) *PostTestInternetDownloadSpeedParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post test internet download speed params
func (o *PostTestInternetDownloadSpeedParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the post test internet download speed params
func (o *PostTestInternetDownloadSpeedParams) WithBody(body PostTestInternetDownloadSpeedBody) *PostTestInternetDownloadSpeedParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the post test internet download speed params
func (o *PostTestInternetDownloadSpeedParams) SetBody(body PostTestInternetDownloadSpeedBody) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *PostTestInternetDownloadSpeedParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
