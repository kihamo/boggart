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

// NewSetPlayerCmdPlaylistParams creates a new SetPlayerCmdPlaylistParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSetPlayerCmdPlaylistParams() *SetPlayerCmdPlaylistParams {
	return &SetPlayerCmdPlaylistParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSetPlayerCmdPlaylistParamsWithTimeout creates a new SetPlayerCmdPlaylistParams object
// with the ability to set a timeout on a request.
func NewSetPlayerCmdPlaylistParamsWithTimeout(timeout time.Duration) *SetPlayerCmdPlaylistParams {
	return &SetPlayerCmdPlaylistParams{
		timeout: timeout,
	}
}

// NewSetPlayerCmdPlaylistParamsWithContext creates a new SetPlayerCmdPlaylistParams object
// with the ability to set a context for a request.
func NewSetPlayerCmdPlaylistParamsWithContext(ctx context.Context) *SetPlayerCmdPlaylistParams {
	return &SetPlayerCmdPlaylistParams{
		Context: ctx,
	}
}

// NewSetPlayerCmdPlaylistParamsWithHTTPClient creates a new SetPlayerCmdPlaylistParams object
// with the ability to set a custom HTTPClient for a request.
func NewSetPlayerCmdPlaylistParamsWithHTTPClient(client *http.Client) *SetPlayerCmdPlaylistParams {
	return &SetPlayerCmdPlaylistParams{
		HTTPClient: client,
	}
}

/*
SetPlayerCmdPlaylistParams contains all the parameters to send to the API endpoint

	for the set player cmd playlist operation.

	Typically these are written to a http.Request.
*/
type SetPlayerCmdPlaylistParams struct {

	/* Index.

	   Is the start index
	*/
	Index int64

	/* URL.

	   Is the m3u or ASX playlist link and should be hexed
	*/
	URL string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the set player cmd playlist params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetPlayerCmdPlaylistParams) WithDefaults() *SetPlayerCmdPlaylistParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the set player cmd playlist params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetPlayerCmdPlaylistParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the set player cmd playlist params
func (o *SetPlayerCmdPlaylistParams) WithTimeout(timeout time.Duration) *SetPlayerCmdPlaylistParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set player cmd playlist params
func (o *SetPlayerCmdPlaylistParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set player cmd playlist params
func (o *SetPlayerCmdPlaylistParams) WithContext(ctx context.Context) *SetPlayerCmdPlaylistParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set player cmd playlist params
func (o *SetPlayerCmdPlaylistParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set player cmd playlist params
func (o *SetPlayerCmdPlaylistParams) WithHTTPClient(client *http.Client) *SetPlayerCmdPlaylistParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set player cmd playlist params
func (o *SetPlayerCmdPlaylistParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithIndex adds the index to the set player cmd playlist params
func (o *SetPlayerCmdPlaylistParams) WithIndex(index int64) *SetPlayerCmdPlaylistParams {
	o.SetIndex(index)
	return o
}

// SetIndex adds the index to the set player cmd playlist params
func (o *SetPlayerCmdPlaylistParams) SetIndex(index int64) {
	o.Index = index
}

// WithURL adds the url to the set player cmd playlist params
func (o *SetPlayerCmdPlaylistParams) WithURL(url string) *SetPlayerCmdPlaylistParams {
	o.SetURL(url)
	return o
}

// SetURL adds the url to the set player cmd playlist params
func (o *SetPlayerCmdPlaylistParams) SetURL(url string) {
	o.URL = url
}

// WriteToRequest writes these params to a swagger request
func (o *SetPlayerCmdPlaylistParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param index
	if err := r.SetPathParam("index", swag.FormatInt64(o.Index)); err != nil {
		return err
	}

	// path param url
	if err := r.SetPathParam("url", o.URL); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
