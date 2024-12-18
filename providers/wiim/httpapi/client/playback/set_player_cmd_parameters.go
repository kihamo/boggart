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

// NewSetPlayerCmdParams creates a new SetPlayerCmdParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSetPlayerCmdParams() *SetPlayerCmdParams {
	return &SetPlayerCmdParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSetPlayerCmdParamsWithTimeout creates a new SetPlayerCmdParams object
// with the ability to set a timeout on a request.
func NewSetPlayerCmdParamsWithTimeout(timeout time.Duration) *SetPlayerCmdParams {
	return &SetPlayerCmdParams{
		timeout: timeout,
	}
}

// NewSetPlayerCmdParamsWithContext creates a new SetPlayerCmdParams object
// with the ability to set a context for a request.
func NewSetPlayerCmdParamsWithContext(ctx context.Context) *SetPlayerCmdParams {
	return &SetPlayerCmdParams{
		Context: ctx,
	}
}

// NewSetPlayerCmdParamsWithHTTPClient creates a new SetPlayerCmdParams object
// with the ability to set a custom HTTPClient for a request.
func NewSetPlayerCmdParamsWithHTTPClient(client *http.Client) *SetPlayerCmdParams {
	return &SetPlayerCmdParams{
		HTTPClient: client,
	}
}

/*
SetPlayerCmdParams contains all the parameters to send to the API endpoint

	for the set player cmd operation.

	Typically these are written to a http.Request.
*/
type SetPlayerCmdParams struct {

	/* URL.

	   A complete URL for an audio source on the internet or addressable local device http://89.223.45.5:8000/progressive-flac example audio file http://stream.live.vc.bbcmedia.co.uk/bbc_6music example radio station file

	*/
	URL string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the set player cmd params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetPlayerCmdParams) WithDefaults() *SetPlayerCmdParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the set player cmd params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetPlayerCmdParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the set player cmd params
func (o *SetPlayerCmdParams) WithTimeout(timeout time.Duration) *SetPlayerCmdParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set player cmd params
func (o *SetPlayerCmdParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set player cmd params
func (o *SetPlayerCmdParams) WithContext(ctx context.Context) *SetPlayerCmdParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set player cmd params
func (o *SetPlayerCmdParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set player cmd params
func (o *SetPlayerCmdParams) WithHTTPClient(client *http.Client) *SetPlayerCmdParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set player cmd params
func (o *SetPlayerCmdParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithURL adds the url to the set player cmd params
func (o *SetPlayerCmdParams) WithURL(url string) *SetPlayerCmdParams {
	o.SetURL(url)
	return o
}

// SetURL adds the url to the set player cmd params
func (o *SetPlayerCmdParams) SetURL(url string) {
	o.URL = url
}

// WriteToRequest writes these params to a swagger request
func (o *SetPlayerCmdParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param url
	if err := r.SetPathParam("url", o.URL); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
