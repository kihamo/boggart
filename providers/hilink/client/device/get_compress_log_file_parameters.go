// Code generated by go-swagger; DO NOT EDIT.

package device

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
)

// NewGetCompressLogFileParams creates a new GetCompressLogFileParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetCompressLogFileParams() *GetCompressLogFileParams {
	return &GetCompressLogFileParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetCompressLogFileParamsWithTimeout creates a new GetCompressLogFileParams object
// with the ability to set a timeout on a request.
func NewGetCompressLogFileParamsWithTimeout(timeout time.Duration) *GetCompressLogFileParams {
	return &GetCompressLogFileParams{
		timeout: timeout,
	}
}

// NewGetCompressLogFileParamsWithContext creates a new GetCompressLogFileParams object
// with the ability to set a context for a request.
func NewGetCompressLogFileParamsWithContext(ctx context.Context) *GetCompressLogFileParams {
	return &GetCompressLogFileParams{
		Context: ctx,
	}
}

// NewGetCompressLogFileParamsWithHTTPClient creates a new GetCompressLogFileParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetCompressLogFileParamsWithHTTPClient(client *http.Client) *GetCompressLogFileParams {
	return &GetCompressLogFileParams{
		HTTPClient: client,
	}
}

/* GetCompressLogFileParams contains all the parameters to send to the API endpoint
   for the get compress log file operation.

   Typically these are written to a http.Request.
*/
type GetCompressLogFileParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get compress log file params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetCompressLogFileParams) WithDefaults() *GetCompressLogFileParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get compress log file params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetCompressLogFileParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get compress log file params
func (o *GetCompressLogFileParams) WithTimeout(timeout time.Duration) *GetCompressLogFileParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get compress log file params
func (o *GetCompressLogFileParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get compress log file params
func (o *GetCompressLogFileParams) WithContext(ctx context.Context) *GetCompressLogFileParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get compress log file params
func (o *GetCompressLogFileParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get compress log file params
func (o *GetCompressLogFileParams) WithHTTPClient(client *http.Client) *GetCompressLogFileParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get compress log file params
func (o *GetCompressLogFileParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetCompressLogFileParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
