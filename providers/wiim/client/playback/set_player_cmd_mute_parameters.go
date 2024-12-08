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
	"github.com/go-openapi/swag"
)

// NewSetPlayerCmdMuteParams creates a new SetPlayerCmdMuteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSetPlayerCmdMuteParams() *SetPlayerCmdMuteParams {
	return &SetPlayerCmdMuteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSetPlayerCmdMuteParamsWithTimeout creates a new SetPlayerCmdMuteParams object
// with the ability to set a timeout on a request.
func NewSetPlayerCmdMuteParamsWithTimeout(timeout time.Duration) *SetPlayerCmdMuteParams {
	return &SetPlayerCmdMuteParams{
		timeout: timeout,
	}
}

// NewSetPlayerCmdMuteParamsWithContext creates a new SetPlayerCmdMuteParams object
// with the ability to set a context for a request.
func NewSetPlayerCmdMuteParamsWithContext(ctx context.Context) *SetPlayerCmdMuteParams {
	return &SetPlayerCmdMuteParams{
		Context: ctx,
	}
}

// NewSetPlayerCmdMuteParamsWithHTTPClient creates a new SetPlayerCmdMuteParams object
// with the ability to set a custom HTTPClient for a request.
func NewSetPlayerCmdMuteParamsWithHTTPClient(client *http.Client) *SetPlayerCmdMuteParams {
	return &SetPlayerCmdMuteParams{
		HTTPClient: client,
	}
}

/*
SetPlayerCmdMuteParams contains all the parameters to send to the API endpoint

	for the set player cmd mute operation.

	Typically these are written to a http.Request.
*/
type SetPlayerCmdMuteParams struct {

	/* Mute.

	   Set the mute mode 0: Not muted 1: Muted

	*/
	Mute int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the set player cmd mute params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetPlayerCmdMuteParams) WithDefaults() *SetPlayerCmdMuteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the set player cmd mute params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetPlayerCmdMuteParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the set player cmd mute params
func (o *SetPlayerCmdMuteParams) WithTimeout(timeout time.Duration) *SetPlayerCmdMuteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set player cmd mute params
func (o *SetPlayerCmdMuteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set player cmd mute params
func (o *SetPlayerCmdMuteParams) WithContext(ctx context.Context) *SetPlayerCmdMuteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set player cmd mute params
func (o *SetPlayerCmdMuteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set player cmd mute params
func (o *SetPlayerCmdMuteParams) WithHTTPClient(client *http.Client) *SetPlayerCmdMuteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set player cmd mute params
func (o *SetPlayerCmdMuteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithMute adds the mute to the set player cmd mute params
func (o *SetPlayerCmdMuteParams) WithMute(mute int64) *SetPlayerCmdMuteParams {
	o.SetMute(mute)
	return o
}

// SetMute adds the mute to the set player cmd mute params
func (o *SetPlayerCmdMuteParams) SetMute(mute int64) {
	o.Mute = mute
}

// WriteToRequest writes these params to a swagger request
func (o *SetPlayerCmdMuteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param mute
	if err := r.SetPathParam("mute", swag.FormatInt64(o.Mute)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
