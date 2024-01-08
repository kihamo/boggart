// Code generated by go-swagger; DO NOT EDIT.

package system

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

	"github.com/kihamo/boggart/providers/hikvision/models"
)

// NewSetNtpServerParams creates a new SetNtpServerParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewSetNtpServerParams() *SetNtpServerParams {
	return &SetNtpServerParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewSetNtpServerParamsWithTimeout creates a new SetNtpServerParams object
// with the ability to set a timeout on a request.
func NewSetNtpServerParamsWithTimeout(timeout time.Duration) *SetNtpServerParams {
	return &SetNtpServerParams{
		timeout: timeout,
	}
}

// NewSetNtpServerParamsWithContext creates a new SetNtpServerParams object
// with the ability to set a context for a request.
func NewSetNtpServerParamsWithContext(ctx context.Context) *SetNtpServerParams {
	return &SetNtpServerParams{
		Context: ctx,
	}
}

// NewSetNtpServerParamsWithHTTPClient creates a new SetNtpServerParams object
// with the ability to set a custom HTTPClient for a request.
func NewSetNtpServerParamsWithHTTPClient(client *http.Client) *SetNtpServerParams {
	return &SetNtpServerParams{
		HTTPClient: client,
	}
}

/*
SetNtpServerParams contains all the parameters to send to the API endpoint

	for the set ntp server operation.

	Typically these are written to a http.Request.
*/
type SetNtpServerParams struct {

	// NTPServer.
	NTPServer *models.NTPServer

	/* ID.

	   NTP server ID

	   Format: uint64
	*/
	ID uint64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the set ntp server params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetNtpServerParams) WithDefaults() *SetNtpServerParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the set ntp server params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *SetNtpServerParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the set ntp server params
func (o *SetNtpServerParams) WithTimeout(timeout time.Duration) *SetNtpServerParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the set ntp server params
func (o *SetNtpServerParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the set ntp server params
func (o *SetNtpServerParams) WithContext(ctx context.Context) *SetNtpServerParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the set ntp server params
func (o *SetNtpServerParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the set ntp server params
func (o *SetNtpServerParams) WithHTTPClient(client *http.Client) *SetNtpServerParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the set ntp server params
func (o *SetNtpServerParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithNTPServer adds the nTPServer to the set ntp server params
func (o *SetNtpServerParams) WithNTPServer(nTPServer *models.NTPServer) *SetNtpServerParams {
	o.SetNTPServer(nTPServer)
	return o
}

// SetNTPServer adds the nTPServer to the set ntp server params
func (o *SetNtpServerParams) SetNTPServer(nTPServer *models.NTPServer) {
	o.NTPServer = nTPServer
}

// WithID adds the id to the set ntp server params
func (o *SetNtpServerParams) WithID(id uint64) *SetNtpServerParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the set ntp server params
func (o *SetNtpServerParams) SetID(id uint64) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *SetNtpServerParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.NTPServer != nil {
		if err := r.SetBodyParam(o.NTPServer); err != nil {
			return err
		}
	}

	// path param id
	if err := r.SetPathParam("id", swag.FormatUint64(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
