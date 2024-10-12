// Code generated by go-swagger; DO NOT EDIT.

package weather

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

// NewGetCurrentByCityIDParams creates a new GetCurrentByCityIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetCurrentByCityIDParams() *GetCurrentByCityIDParams {
	return &GetCurrentByCityIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetCurrentByCityIDParamsWithTimeout creates a new GetCurrentByCityIDParams object
// with the ability to set a timeout on a request.
func NewGetCurrentByCityIDParamsWithTimeout(timeout time.Duration) *GetCurrentByCityIDParams {
	return &GetCurrentByCityIDParams{
		timeout: timeout,
	}
}

// NewGetCurrentByCityIDParamsWithContext creates a new GetCurrentByCityIDParams object
// with the ability to set a context for a request.
func NewGetCurrentByCityIDParamsWithContext(ctx context.Context) *GetCurrentByCityIDParams {
	return &GetCurrentByCityIDParams{
		Context: ctx,
	}
}

// NewGetCurrentByCityIDParamsWithHTTPClient creates a new GetCurrentByCityIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetCurrentByCityIDParamsWithHTTPClient(client *http.Client) *GetCurrentByCityIDParams {
	return &GetCurrentByCityIDParams{
		HTTPClient: client,
	}
}

/*
GetCurrentByCityIDParams contains all the parameters to send to the API endpoint

	for the get current by city ID operation.

	Typically these are written to a http.Request.
*/
type GetCurrentByCityIDParams struct {

	/* ID.

	   City ID

	   Format: uint64
	*/
	ID uint64

	/* Lang.

	   Multilingual support
	*/
	Lang *string

	/* Units.

	   Standard, metric, and imperial units are available
	*/
	Units *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get current by city ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetCurrentByCityIDParams) WithDefaults() *GetCurrentByCityIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get current by city ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetCurrentByCityIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get current by city ID params
func (o *GetCurrentByCityIDParams) WithTimeout(timeout time.Duration) *GetCurrentByCityIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get current by city ID params
func (o *GetCurrentByCityIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get current by city ID params
func (o *GetCurrentByCityIDParams) WithContext(ctx context.Context) *GetCurrentByCityIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get current by city ID params
func (o *GetCurrentByCityIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get current by city ID params
func (o *GetCurrentByCityIDParams) WithHTTPClient(client *http.Client) *GetCurrentByCityIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get current by city ID params
func (o *GetCurrentByCityIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get current by city ID params
func (o *GetCurrentByCityIDParams) WithID(id uint64) *GetCurrentByCityIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get current by city ID params
func (o *GetCurrentByCityIDParams) SetID(id uint64) {
	o.ID = id
}

// WithLang adds the lang to the get current by city ID params
func (o *GetCurrentByCityIDParams) WithLang(lang *string) *GetCurrentByCityIDParams {
	o.SetLang(lang)
	return o
}

// SetLang adds the lang to the get current by city ID params
func (o *GetCurrentByCityIDParams) SetLang(lang *string) {
	o.Lang = lang
}

// WithUnits adds the units to the get current by city ID params
func (o *GetCurrentByCityIDParams) WithUnits(units *string) *GetCurrentByCityIDParams {
	o.SetUnits(units)
	return o
}

// SetUnits adds the units to the get current by city ID params
func (o *GetCurrentByCityIDParams) SetUnits(units *string) {
	o.Units = units
}

// WriteToRequest writes these params to a swagger request
func (o *GetCurrentByCityIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param id
	qrID := o.ID
	qID := swag.FormatUint64(qrID)
	if qID != "" {

		if err := r.SetQueryParam("id", qID); err != nil {
			return err
		}
	}

	if o.Lang != nil {

		// query param lang
		var qrLang string

		if o.Lang != nil {
			qrLang = *o.Lang
		}
		qLang := qrLang
		if qLang != "" {

			if err := r.SetQueryParam("lang", qLang); err != nil {
				return err
			}
		}
	}

	if o.Units != nil {

		// query param units
		var qrUnits string

		if o.Units != nil {
			qrUnits = *o.Units
		}
		qUnits := qrUnits
		if qUnits != "" {

			if err := r.SetQueryParam("units", qUnits); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
