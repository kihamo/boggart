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

// NewSetPlayerCmdPlaylistHEXParams creates a new SetPlayerCmdPlaylistHEXParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSetPlayerCmdPlaylistHEXParams() *SetPlayerCmdPlaylistHEXParams {
	return &SetPlayerCmdPlaylistHEXParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSetPlayerCmdPlaylistHEXParamsWithTimeout creates a new SetPlayerCmdPlaylistHEXParams object
// with the ability to set a timeout on a request.
func NewSetPlayerCmdPlaylistHEXParamsWithTimeout(timeout time.Duration) *SetPlayerCmdPlaylistHEXParams {
	return &SetPlayerCmdPlaylistHEXParams{
		timeout: timeout,
	}
}

// NewSetPlayerCmdPlaylistHEXParamsWithContext creates a new SetPlayerCmdPlaylistHEXParams object
// with the ability to set a context for a request.
func NewSetPlayerCmdPlaylistHEXParamsWithContext(ctx context.Context) *SetPlayerCmdPlaylistHEXParams {
	return &SetPlayerCmdPlaylistHEXParams{
		Context: ctx,
	}
}

// NewSetPlayerCmdPlaylistHEXParamsWithHTTPClient creates a new SetPlayerCmdPlaylistHEXParams object
// with the ability to set a custom HTTPClient for a request.
func NewSetPlayerCmdPlaylistHEXParamsWithHTTPClient(client *http.Client) *SetPlayerCmdPlaylistHEXParams {
	return &SetPlayerCmdPlaylistHEXParams{
		HTTPClient: client,
	}
}

/*
SetPlayerCmdPlaylistHEXParams contains all the parameters to send to the API endpoint

	for the set player cmd playlist h e x operation.

	Typically these are written to a http.Request.
*/
type SetPlayerCmdPlaylistHEXParams struct {

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

// WithDefaults hydrates default values in the set player cmd playlist h e x params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetPlayerCmdPlaylistHEXParams) WithDefaults() *SetPlayerCmdPlaylistHEXParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the set player cmd playlist h e x params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetPlayerCmdPlaylistHEXParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the set player cmd playlist h e x params
func (o *SetPlayerCmdPlaylistHEXParams) WithTimeout(timeout time.Duration) *SetPlayerCmdPlaylistHEXParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set player cmd playlist h e x params
func (o *SetPlayerCmdPlaylistHEXParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set player cmd playlist h e x params
func (o *SetPlayerCmdPlaylistHEXParams) WithContext(ctx context.Context) *SetPlayerCmdPlaylistHEXParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set player cmd playlist h e x params
func (o *SetPlayerCmdPlaylistHEXParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set player cmd playlist h e x params
func (o *SetPlayerCmdPlaylistHEXParams) WithHTTPClient(client *http.Client) *SetPlayerCmdPlaylistHEXParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set player cmd playlist h e x params
func (o *SetPlayerCmdPlaylistHEXParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithIndex adds the index to the set player cmd playlist h e x params
func (o *SetPlayerCmdPlaylistHEXParams) WithIndex(index int64) *SetPlayerCmdPlaylistHEXParams {
	o.SetIndex(index)
	return o
}

// SetIndex adds the index to the set player cmd playlist h e x params
func (o *SetPlayerCmdPlaylistHEXParams) SetIndex(index int64) {
	o.Index = index
}

// WithURL adds the url to the set player cmd playlist h e x params
func (o *SetPlayerCmdPlaylistHEXParams) WithURL(url string) *SetPlayerCmdPlaylistHEXParams {
	o.SetURL(url)
	return o
}

// SetURL adds the url to the set player cmd playlist h e x params
func (o *SetPlayerCmdPlaylistHEXParams) SetURL(url string) {
	o.URL = url
}

// WriteToRequest writes these params to a swagger request
func (o *SetPlayerCmdPlaylistHEXParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
