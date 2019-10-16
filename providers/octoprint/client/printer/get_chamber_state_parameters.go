// Code generated by go-swagger; DO NOT EDIT.

package printer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetChamberStateParams creates a new GetChamberStateParams object
// with the default values initialized.
func NewGetChamberStateParams() *GetChamberStateParams {
	var ()
	return &GetChamberStateParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetChamberStateParamsWithTimeout creates a new GetChamberStateParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetChamberStateParamsWithTimeout(timeout time.Duration) *GetChamberStateParams {
	var ()
	return &GetChamberStateParams{

		timeout: timeout,
	}
}

// NewGetChamberStateParamsWithContext creates a new GetChamberStateParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetChamberStateParamsWithContext(ctx context.Context) *GetChamberStateParams {
	var ()
	return &GetChamberStateParams{

		Context: ctx,
	}
}

// NewGetChamberStateParamsWithHTTPClient creates a new GetChamberStateParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetChamberStateParamsWithHTTPClient(client *http.Client) *GetChamberStateParams {
	var ()
	return &GetChamberStateParams{
		HTTPClient: client,
	}
}

/*GetChamberStateParams contains all the parameters to send to the API endpoint
for the get chamber state operation typically these are written to a http.Request
*/
type GetChamberStateParams struct {

	/*History
	  The printer’s temperature history by supplying

	*/
	History *bool
	/*Limit
	  The amount of data points limited

	*/
	Limit *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get chamber state params
func (o *GetChamberStateParams) WithTimeout(timeout time.Duration) *GetChamberStateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get chamber state params
func (o *GetChamberStateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get chamber state params
func (o *GetChamberStateParams) WithContext(ctx context.Context) *GetChamberStateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get chamber state params
func (o *GetChamberStateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get chamber state params
func (o *GetChamberStateParams) WithHTTPClient(client *http.Client) *GetChamberStateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get chamber state params
func (o *GetChamberStateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithHistory adds the history to the get chamber state params
func (o *GetChamberStateParams) WithHistory(history *bool) *GetChamberStateParams {
	o.SetHistory(history)
	return o
}

// SetHistory adds the history to the get chamber state params
func (o *GetChamberStateParams) SetHistory(history *bool) {
	o.History = history
}

// WithLimit adds the limit to the get chamber state params
func (o *GetChamberStateParams) WithLimit(limit *int64) *GetChamberStateParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get chamber state params
func (o *GetChamberStateParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WriteToRequest writes these params to a swagger request
func (o *GetChamberStateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.History != nil {

		// query param history
		var qrHistory bool
		if o.History != nil {
			qrHistory = *o.History
		}
		qHistory := swag.FormatBool(qrHistory)
		if qHistory != "" {
			if err := r.SetQueryParam("history", qHistory); err != nil {
				return err
			}
		}

	}

	if o.Limit != nil {

		// query param limit
		var qrLimit int64
		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt64(qrLimit)
		if qLimit != "" {
			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}