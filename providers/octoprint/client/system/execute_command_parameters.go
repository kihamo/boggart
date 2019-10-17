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

	strfmt "github.com/go-openapi/strfmt"
)

// NewExecuteCommandParams creates a new ExecuteCommandParams object
// with the default values initialized.
func NewExecuteCommandParams() *ExecuteCommandParams {
	var ()
	return &ExecuteCommandParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewExecuteCommandParamsWithTimeout creates a new ExecuteCommandParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewExecuteCommandParamsWithTimeout(timeout time.Duration) *ExecuteCommandParams {
	var ()
	return &ExecuteCommandParams{

		timeout: timeout,
	}
}

// NewExecuteCommandParamsWithContext creates a new ExecuteCommandParams object
// with the default values initialized, and the ability to set a context for a request
func NewExecuteCommandParamsWithContext(ctx context.Context) *ExecuteCommandParams {
	var ()
	return &ExecuteCommandParams{

		Context: ctx,
	}
}

// NewExecuteCommandParamsWithHTTPClient creates a new ExecuteCommandParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewExecuteCommandParamsWithHTTPClient(client *http.Client) *ExecuteCommandParams {
	var ()
	return &ExecuteCommandParams{
		HTTPClient: client,
	}
}

/*ExecuteCommandParams contains all the parameters to send to the API endpoint
for the execute command operation typically these are written to a http.Request
*/
type ExecuteCommandParams struct {

	/*Action
	  The identifier of the command

	*/
	Action string
	/*Source
	  The source for which to list commands

	*/
	Source string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the execute command params
func (o *ExecuteCommandParams) WithTimeout(timeout time.Duration) *ExecuteCommandParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the execute command params
func (o *ExecuteCommandParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the execute command params
func (o *ExecuteCommandParams) WithContext(ctx context.Context) *ExecuteCommandParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the execute command params
func (o *ExecuteCommandParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the execute command params
func (o *ExecuteCommandParams) WithHTTPClient(client *http.Client) *ExecuteCommandParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the execute command params
func (o *ExecuteCommandParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAction adds the action to the execute command params
func (o *ExecuteCommandParams) WithAction(action string) *ExecuteCommandParams {
	o.SetAction(action)
	return o
}

// SetAction adds the action to the execute command params
func (o *ExecuteCommandParams) SetAction(action string) {
	o.Action = action
}

// WithSource adds the source to the execute command params
func (o *ExecuteCommandParams) WithSource(source string) *ExecuteCommandParams {
	o.SetSource(source)
	return o
}

// SetSource adds the source to the execute command params
func (o *ExecuteCommandParams) SetSource(source string) {
	o.Source = source
}

// WriteToRequest writes these params to a swagger request
func (o *ExecuteCommandParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param action
	if err := r.SetPathParam("action", o.Action); err != nil {
		return err
	}

	// path param source
	if err := r.SetPathParam("source", o.Source); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
