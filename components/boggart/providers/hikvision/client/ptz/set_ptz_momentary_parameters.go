// Code generated by go-swagger; DO NOT EDIT.

package ptz

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

	models "github.com/kihamo/boggart/components/boggart/providers/hikvision/models"
)

// NewSetPtzMomentaryParams creates a new SetPtzMomentaryParams object
// with the default values initialized.
func NewSetPtzMomentaryParams() *SetPtzMomentaryParams {
	var ()
	return &SetPtzMomentaryParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewSetPtzMomentaryParamsWithTimeout creates a new SetPtzMomentaryParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewSetPtzMomentaryParamsWithTimeout(timeout time.Duration) *SetPtzMomentaryParams {
	var ()
	return &SetPtzMomentaryParams{

		timeout: timeout,
	}
}

// NewSetPtzMomentaryParamsWithContext creates a new SetPtzMomentaryParams object
// with the default values initialized, and the ability to set a context for a request
func NewSetPtzMomentaryParamsWithContext(ctx context.Context) *SetPtzMomentaryParams {
	var ()
	return &SetPtzMomentaryParams{

		Context: ctx,
	}
}

// NewSetPtzMomentaryParamsWithHTTPClient creates a new SetPtzMomentaryParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewSetPtzMomentaryParamsWithHTTPClient(client *http.Client) *SetPtzMomentaryParams {
	var ()
	return &SetPtzMomentaryParams{
		HTTPClient: client,
	}
}

/*SetPtzMomentaryParams contains all the parameters to send to the API endpoint
for the set ptz momentary operation typically these are written to a http.Request
*/
type SetPtzMomentaryParams struct {

	/*PTZData*/
	PTZData *models.PTZData
	/*Channel
	  Channel ID

	*/
	Channel uint64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the set ptz momentary params
func (o *SetPtzMomentaryParams) WithTimeout(timeout time.Duration) *SetPtzMomentaryParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set ptz momentary params
func (o *SetPtzMomentaryParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set ptz momentary params
func (o *SetPtzMomentaryParams) WithContext(ctx context.Context) *SetPtzMomentaryParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set ptz momentary params
func (o *SetPtzMomentaryParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set ptz momentary params
func (o *SetPtzMomentaryParams) WithHTTPClient(client *http.Client) *SetPtzMomentaryParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set ptz momentary params
func (o *SetPtzMomentaryParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPTZData adds the pTZData to the set ptz momentary params
func (o *SetPtzMomentaryParams) WithPTZData(pTZData *models.PTZData) *SetPtzMomentaryParams {
	o.SetPTZData(pTZData)
	return o
}

// SetPTZData adds the pTZData to the set ptz momentary params
func (o *SetPtzMomentaryParams) SetPTZData(pTZData *models.PTZData) {
	o.PTZData = pTZData
}

// WithChannel adds the channel to the set ptz momentary params
func (o *SetPtzMomentaryParams) WithChannel(channel uint64) *SetPtzMomentaryParams {
	o.SetChannel(channel)
	return o
}

// SetChannel adds the channel to the set ptz momentary params
func (o *SetPtzMomentaryParams) SetChannel(channel uint64) {
	o.Channel = channel
}

// WriteToRequest writes these params to a swagger request
func (o *SetPtzMomentaryParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.PTZData != nil {
		if err := r.SetBodyParam(o.PTZData); err != nil {
			return err
		}
	}

	// path param channel
	if err := r.SetPathParam("channel", swag.FormatUint64(o.Channel)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}