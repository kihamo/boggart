// Code generated by go-swagger; DO NOT EDIT.

package forecast

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

// NewGetForecastByGeographicCoordinatesParams creates a new GetForecastByGeographicCoordinatesParams object
// with the default values initialized.
func NewGetForecastByGeographicCoordinatesParams() *GetForecastByGeographicCoordinatesParams {
	var ()
	return &GetForecastByGeographicCoordinatesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetForecastByGeographicCoordinatesParamsWithTimeout creates a new GetForecastByGeographicCoordinatesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetForecastByGeographicCoordinatesParamsWithTimeout(timeout time.Duration) *GetForecastByGeographicCoordinatesParams {
	var ()
	return &GetForecastByGeographicCoordinatesParams{

		timeout: timeout,
	}
}

// NewGetForecastByGeographicCoordinatesParamsWithContext creates a new GetForecastByGeographicCoordinatesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetForecastByGeographicCoordinatesParamsWithContext(ctx context.Context) *GetForecastByGeographicCoordinatesParams {
	var ()
	return &GetForecastByGeographicCoordinatesParams{

		Context: ctx,
	}
}

// NewGetForecastByGeographicCoordinatesParamsWithHTTPClient creates a new GetForecastByGeographicCoordinatesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetForecastByGeographicCoordinatesParamsWithHTTPClient(client *http.Client) *GetForecastByGeographicCoordinatesParams {
	var ()
	return &GetForecastByGeographicCoordinatesParams{
		HTTPClient: client,
	}
}

/*GetForecastByGeographicCoordinatesParams contains all the parameters to send to the API endpoint
for the get forecast by geographic coordinates operation typically these are written to a http.Request
*/
type GetForecastByGeographicCoordinatesParams struct {

	/*Cnt
	  To limit number of listed cities please setup 'cnt' parameter that specifies the number of lines returned

	*/
	Cnt *uint64
	/*Lang
	  Multilingual support

	*/
	Lang *string
	/*Lat
	  Coordinates of the location of your interest

	*/
	Lat float64
	/*Lon
	  Coordinates of the location of your interest

	*/
	Lon float64
	/*Units
	  Standard, metric, and imperial units are available

	*/
	Units *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) WithTimeout(timeout time.Duration) *GetForecastByGeographicCoordinatesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) WithContext(ctx context.Context) *GetForecastByGeographicCoordinatesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) WithHTTPClient(client *http.Client) *GetForecastByGeographicCoordinatesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCnt adds the cnt to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) WithCnt(cnt *uint64) *GetForecastByGeographicCoordinatesParams {
	o.SetCnt(cnt)
	return o
}

// SetCnt adds the cnt to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) SetCnt(cnt *uint64) {
	o.Cnt = cnt
}

// WithLang adds the lang to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) WithLang(lang *string) *GetForecastByGeographicCoordinatesParams {
	o.SetLang(lang)
	return o
}

// SetLang adds the lang to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) SetLang(lang *string) {
	o.Lang = lang
}

// WithLat adds the lat to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) WithLat(lat float64) *GetForecastByGeographicCoordinatesParams {
	o.SetLat(lat)
	return o
}

// SetLat adds the lat to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) SetLat(lat float64) {
	o.Lat = lat
}

// WithLon adds the lon to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) WithLon(lon float64) *GetForecastByGeographicCoordinatesParams {
	o.SetLon(lon)
	return o
}

// SetLon adds the lon to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) SetLon(lon float64) {
	o.Lon = lon
}

// WithUnits adds the units to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) WithUnits(units *string) *GetForecastByGeographicCoordinatesParams {
	o.SetUnits(units)
	return o
}

// SetUnits adds the units to the get forecast by geographic coordinates params
func (o *GetForecastByGeographicCoordinatesParams) SetUnits(units *string) {
	o.Units = units
}

// WriteToRequest writes these params to a swagger request
func (o *GetForecastByGeographicCoordinatesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Cnt != nil {

		// query param cnt
		var qrCnt uint64
		if o.Cnt != nil {
			qrCnt = *o.Cnt
		}
		qCnt := swag.FormatUint64(qrCnt)
		if qCnt != "" {
			if err := r.SetQueryParam("cnt", qCnt); err != nil {
				return err
			}
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