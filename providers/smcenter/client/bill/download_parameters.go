// Code generated by go-swagger; DO NOT EDIT.

package bill

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
	"github.com/go-openapi/swag"
)

// NewDownloadParams creates a new DownloadParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDownloadParams() *DownloadParams {
	return &DownloadParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDownloadParamsWithTimeout creates a new DownloadParams object
// with the ability to set a timeout on a request.
func NewDownloadParamsWithTimeout(timeout time.Duration) *DownloadParams {
	return &DownloadParams{
		timeout: timeout,
	}
}

// NewDownloadParamsWithContext creates a new DownloadParams object
// with the ability to set a context for a request.
func NewDownloadParamsWithContext(ctx context.Context) *DownloadParams {
	return &DownloadParams{
		Context: ctx,
	}
}

// NewDownloadParamsWithHTTPClient creates a new DownloadParams object
// with the ability to set a custom HTTPClient for a request.
func NewDownloadParamsWithHTTPClient(client *http.Client) *DownloadParams {
	return &DownloadParams{
		HTTPClient: client,
	}
}

/* DownloadParams contains all the parameters to send to the API endpoint
   for the download operation.

   Typically these are written to a http.Request.
*/
type DownloadParams struct {

	/* ID.

	   Bill ID

	   Format: uint64
	*/
	ID uint64

	/* InJpg.

	   JPEG format or not

	   Format: uint64
	*/
	InJpg *uint64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the download params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DownloadParams) WithDefaults() *DownloadParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the download params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DownloadParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the download params
func (o *DownloadParams) WithTimeout(timeout time.Duration) *DownloadParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the download params
func (o *DownloadParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the download params
func (o *DownloadParams) WithContext(ctx context.Context) *DownloadParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the download params
func (o *DownloadParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the download params
func (o *DownloadParams) WithHTTPClient(client *http.Client) *DownloadParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the download params
func (o *DownloadParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the download params
func (o *DownloadParams) WithID(id uint64) *DownloadParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the download params
func (o *DownloadParams) SetID(id uint64) {
	o.ID = id
}

// WithInJpg adds the inJpg to the download params
func (o *DownloadParams) WithInJpg(inJpg *uint64) *DownloadParams {
	o.SetInJpg(inJpg)
	return o
}

// SetInJpg adds the inJpg to the download params
func (o *DownloadParams) SetInJpg(inJpg *uint64) {
	o.InJpg = inJpg
}

// WriteToRequest writes these params to a swagger request
func (o *DownloadParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatUint64(o.ID)); err != nil {
		return err
	}

	if o.InJpg != nil {

		// query param inJpg
		var qrInJpg uint64

		if o.InJpg != nil {
			qrInJpg = *o.InJpg
		}
		qInJpg := swag.FormatUint64(qrInJpg)
		if qInJpg != "" {

			if err := r.SetQueryParam("inJpg", qInJpg); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
