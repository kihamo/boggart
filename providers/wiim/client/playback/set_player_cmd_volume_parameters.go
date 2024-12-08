// Code generated by go-swagger; DO NOT EDIT.

package playback

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

// NewSetPlayerCmdVolumeParams creates a new SetPlayerCmdVolumeParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSetPlayerCmdVolumeParams() *SetPlayerCmdVolumeParams {
	return &SetPlayerCmdVolumeParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSetPlayerCmdVolumeParamsWithTimeout creates a new SetPlayerCmdVolumeParams object
// with the ability to set a timeout on a request.
func NewSetPlayerCmdVolumeParamsWithTimeout(timeout time.Duration) *SetPlayerCmdVolumeParams {
	return &SetPlayerCmdVolumeParams{
		timeout: timeout,
	}
}

// NewSetPlayerCmdVolumeParamsWithContext creates a new SetPlayerCmdVolumeParams object
// with the ability to set a context for a request.
func NewSetPlayerCmdVolumeParamsWithContext(ctx context.Context) *SetPlayerCmdVolumeParams {
	return &SetPlayerCmdVolumeParams{
		Context: ctx,
	}
}

// NewSetPlayerCmdVolumeParamsWithHTTPClient creates a new SetPlayerCmdVolumeParams object
// with the ability to set a custom HTTPClient for a request.
func NewSetPlayerCmdVolumeParamsWithHTTPClient(client *http.Client) *SetPlayerCmdVolumeParams {
	return &SetPlayerCmdVolumeParams{
		HTTPClient: client,
	}
}

/*
SetPlayerCmdVolumeParams contains all the parameters to send to the API endpoint

	for the set player cmd volume operation.

	Typically these are written to a http.Request.
*/
type SetPlayerCmdVolumeParams struct {

	/* Volume.

	   Adjust volume for current device :vol: direct volue, value range is 0-100 --: Decrease by 2 %2b%2b: Increase by 2

	*/
	Volume string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the set player cmd volume params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetPlayerCmdVolumeParams) WithDefaults() *SetPlayerCmdVolumeParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the set player cmd volume params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetPlayerCmdVolumeParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the set player cmd volume params
func (o *SetPlayerCmdVolumeParams) WithTimeout(timeout time.Duration) *SetPlayerCmdVolumeParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set player cmd volume params
func (o *SetPlayerCmdVolumeParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set player cmd volume params
func (o *SetPlayerCmdVolumeParams) WithContext(ctx context.Context) *SetPlayerCmdVolumeParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set player cmd volume params
func (o *SetPlayerCmdVolumeParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set player cmd volume params
func (o *SetPlayerCmdVolumeParams) WithHTTPClient(client *http.Client) *SetPlayerCmdVolumeParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set player cmd volume params
func (o *SetPlayerCmdVolumeParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithVolume adds the volume to the set player cmd volume params
func (o *SetPlayerCmdVolumeParams) WithVolume(volume string) *SetPlayerCmdVolumeParams {
	o.SetVolume(volume)
	return o
}

// SetVolume adds the volume to the set player cmd volume params
func (o *SetPlayerCmdVolumeParams) SetVolume(volume string) {
	o.Volume = volume
}

// WriteToRequest writes these params to a swagger request
func (o *SetPlayerCmdVolumeParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param volume
	if err := r.SetPathParam("volume", o.Volume); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
