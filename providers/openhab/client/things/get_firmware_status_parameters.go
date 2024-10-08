// Code generated by go-swagger; DO NOT EDIT.

package things

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

// NewGetFirmwareStatusParams creates a new GetFirmwareStatusParams object
// with the default values initialized.
func NewGetFirmwareStatusParams() *GetFirmwareStatusParams {
	var ()
	return &GetFirmwareStatusParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetFirmwareStatusParamsWithTimeout creates a new GetFirmwareStatusParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetFirmwareStatusParamsWithTimeout(timeout time.Duration) *GetFirmwareStatusParams {
	var ()
	return &GetFirmwareStatusParams{

		timeout: timeout,
	}
}

// NewGetFirmwareStatusParamsWithContext creates a new GetFirmwareStatusParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetFirmwareStatusParamsWithContext(ctx context.Context) *GetFirmwareStatusParams {
	var ()
	return &GetFirmwareStatusParams{

		Context: ctx,
	}
}

// NewGetFirmwareStatusParamsWithHTTPClient creates a new GetFirmwareStatusParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetFirmwareStatusParamsWithHTTPClient(client *http.Client) *GetFirmwareStatusParams {
	var ()
	return &GetFirmwareStatusParams{
		HTTPClient: client,
	}
}

/*GetFirmwareStatusParams contains all the parameters to send to the API endpoint
for the get firmware status operation typically these are written to a http.Request
*/
type GetFirmwareStatusParams struct {

	/*AcceptLanguage*/
	AcceptLanguage *string
	/*ThingUID
	  thing

	*/
	ThingUID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get firmware status params
func (o *GetFirmwareStatusParams) WithTimeout(timeout time.Duration) *GetFirmwareStatusParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get firmware status params
func (o *GetFirmwareStatusParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get firmware status params
func (o *GetFirmwareStatusParams) WithContext(ctx context.Context) *GetFirmwareStatusParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get firmware status params
func (o *GetFirmwareStatusParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get firmware status params
func (o *GetFirmwareStatusParams) WithHTTPClient(client *http.Client) *GetFirmwareStatusParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get firmware status params
func (o *GetFirmwareStatusParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAcceptLanguage adds the acceptLanguage to the get firmware status params
func (o *GetFirmwareStatusParams) WithAcceptLanguage(acceptLanguage *string) *GetFirmwareStatusParams {
	o.SetAcceptLanguage(acceptLanguage)
	return o
}

// SetAcceptLanguage adds the acceptLanguage to the get firmware status params
func (o *GetFirmwareStatusParams) SetAcceptLanguage(acceptLanguage *string) {
	o.AcceptLanguage = acceptLanguage
}

// WithThingUID adds the thingUID to the get firmware status params
func (o *GetFirmwareStatusParams) WithThingUID(thingUID string) *GetFirmwareStatusParams {
	o.SetThingUID(thingUID)
	return o
}

// SetThingUID adds the thingUid to the get firmware status params
func (o *GetFirmwareStatusParams) SetThingUID(thingUID string) {
	o.ThingUID = thingUID
}

// WriteToRequest writes these params to a swagger request
func (o *GetFirmwareStatusParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param thingUID
	if err := r.SetPathParam("thingUID", o.ThingUID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
