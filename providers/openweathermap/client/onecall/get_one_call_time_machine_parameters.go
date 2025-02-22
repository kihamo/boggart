// Code generated by go-swagger; DO NOT EDIT.

package onecall

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

// NewGetOneCallTimeMachineParams creates a new GetOneCallTimeMachineParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetOneCallTimeMachineParams() *GetOneCallTimeMachineParams {
	return &GetOneCallTimeMachineParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetOneCallTimeMachineParamsWithTimeout creates a new GetOneCallTimeMachineParams object
// with the ability to set a timeout on a request.
func NewGetOneCallTimeMachineParamsWithTimeout(timeout time.Duration) *GetOneCallTimeMachineParams {
	return &GetOneCallTimeMachineParams{
		timeout: timeout,
	}
}

// NewGetOneCallTimeMachineParamsWithContext creates a new GetOneCallTimeMachineParams object
// with the ability to set a context for a request.
func NewGetOneCallTimeMachineParamsWithContext(ctx context.Context) *GetOneCallTimeMachineParams {
	return &GetOneCallTimeMachineParams{
		Context: ctx,
	}
}

// NewGetOneCallTimeMachineParamsWithHTTPClient creates a new GetOneCallTimeMachineParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetOneCallTimeMachineParamsWithHTTPClient(client *http.Client) *GetOneCallTimeMachineParams {
	return &GetOneCallTimeMachineParams{
		HTTPClient: client,
	}
}

/*
GetOneCallTimeMachineParams contains all the parameters to send to the API endpoint

	for the get one call time machine operation.

	Typically these are written to a http.Request.
*/
type GetOneCallTimeMachineParams struct {

	/* Dt.

	   Exclude some parts of the weather data from the API response. It should be a comma-delimited list (without spaces)

	   Format: uint64
	*/
	Dt uint64

	/* Lang.

	   Multilingual support
	*/
	Lang *string

	/* Lat.

	   Coordinates of the location of your interest
	*/
	Lat float64

	/* Lon.

	   Coordinates of the location of your interest
	*/
	Lon float64

	/* Units.

	   Standard, metric, and imperial units are available
	*/
	Units *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get one call time machine params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetOneCallTimeMachineParams) WithDefaults() *GetOneCallTimeMachineParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get one call time machine params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetOneCallTimeMachineParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get one call time machine params
func (o *GetOneCallTimeMachineParams) WithTimeout(timeout time.Duration) *GetOneCallTimeMachineParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get one call time machine params
func (o *GetOneCallTimeMachineParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get one call time machine params
func (o *GetOneCallTimeMachineParams) WithContext(ctx context.Context) *GetOneCallTimeMachineParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get one call time machine params
func (o *GetOneCallTimeMachineParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get one call time machine params
func (o *GetOneCallTimeMachineParams) WithHTTPClient(client *http.Client) *GetOneCallTimeMachineParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get one call time machine params
func (o *GetOneCallTimeMachineParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDt adds the dt to the get one call time machine params
func (o *GetOneCallTimeMachineParams) WithDt(dt uint64) *GetOneCallTimeMachineParams {
	o.SetDt(dt)
	return o
}

// SetDt adds the dt to the get one call time machine params
func (o *GetOneCallTimeMachineParams) SetDt(dt uint64) {
	o.Dt = dt
}

// WithLang adds the lang to the get one call time machine params
func (o *GetOneCallTimeMachineParams) WithLang(lang *string) *GetOneCallTimeMachineParams {
	o.SetLang(lang)
	return o
}

// SetLang adds the lang to the get one call time machine params
func (o *GetOneCallTimeMachineParams) SetLang(lang *string) {
	o.Lang = lang
}

// WithLat adds the lat to the get one call time machine params
func (o *GetOneCallTimeMachineParams) WithLat(lat float64) *GetOneCallTimeMachineParams {
	o.SetLat(lat)
	return o
}

// SetLat adds the lat to the get one call time machine params
func (o *GetOneCallTimeMachineParams) SetLat(lat float64) {
	o.Lat = lat
}

// WithLon adds the lon to the get one call time machine params
func (o *GetOneCallTimeMachineParams) WithLon(lon float64) *GetOneCallTimeMachineParams {
	o.SetLon(lon)
	return o
}

// SetLon adds the lon to the get one call time machine params
func (o *GetOneCallTimeMachineParams) SetLon(lon float64) {
	o.Lon = lon
}

// WithUnits adds the units to the get one call time machine params
func (o *GetOneCallTimeMachineParams) WithUnits(units *string) *GetOneCallTimeMachineParams {
	o.SetUnits(units)
	return o
}

// SetUnits adds the units to the get one call time machine params
func (o *GetOneCallTimeMachineParams) SetUnits(units *string) {
	o.Units = units
}

// WriteToRequest writes these params to a swagger request
func (o *GetOneCallTimeMachineParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param dt
	qrDt := o.Dt
	qDt := swag.FormatUint64(qrDt)
	if qDt != "" {

		if err := r.SetQueryParam("dt", qDt); err != nil {
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

	// query param lat
	qrLat := o.Lat
	qLat := swag.FormatFloat64(qrLat)
	if qLat != "" {

		if err := r.SetQueryParam("lat", qLat); err != nil {
			return err
		}
	}

	// query param lon
	qrLon := o.Lon
	qLon := swag.FormatFloat64(qrLon)
	if qLon != "" {

		if err := r.SetQueryParam("lon", qLon); err != nil {
			return err
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
