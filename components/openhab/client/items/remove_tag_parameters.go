// Code generated by go-swagger; DO NOT EDIT.

package items

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewRemoveTagParams creates a new RemoveTagParams object
// with the default values initialized.
func NewRemoveTagParams() *RemoveTagParams {
	var ()
	return &RemoveTagParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewRemoveTagParamsWithTimeout creates a new RemoveTagParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewRemoveTagParamsWithTimeout(timeout time.Duration) *RemoveTagParams {
	var ()
	return &RemoveTagParams{

		timeout: timeout,
	}
}

// NewRemoveTagParamsWithContext creates a new RemoveTagParams object
// with the default values initialized, and the ability to set a context for a request
func NewRemoveTagParamsWithContext(ctx context.Context) *RemoveTagParams {
	var ()
	return &RemoveTagParams{

		Context: ctx,
	}
}

// NewRemoveTagParamsWithHTTPClient creates a new RemoveTagParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewRemoveTagParamsWithHTTPClient(client *http.Client) *RemoveTagParams {
	var ()
	return &RemoveTagParams{
		HTTPClient: client,
	}
}

/*RemoveTagParams contains all the parameters to send to the API endpoint
for the remove tag operation typically these are written to a http.Request
*/
type RemoveTagParams struct {

	/*Itemname
	  item name

	*/
	Itemname string
	/*Tag
	  tag

	*/
	Tag string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the remove tag params
func (o *RemoveTagParams) WithTimeout(timeout time.Duration) *RemoveTagParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the remove tag params
func (o *RemoveTagParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the remove tag params
func (o *RemoveTagParams) WithContext(ctx context.Context) *RemoveTagParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the remove tag params
func (o *RemoveTagParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the remove tag params
func (o *RemoveTagParams) WithHTTPClient(client *http.Client) *RemoveTagParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the remove tag params
func (o *RemoveTagParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithItemname adds the itemname to the remove tag params
func (o *RemoveTagParams) WithItemname(itemname string) *RemoveTagParams {
	o.SetItemname(itemname)
	return o
}

// SetItemname adds the itemname to the remove tag params
func (o *RemoveTagParams) SetItemname(itemname string) {
	o.Itemname = itemname
}

// WithTag adds the tag to the remove tag params
func (o *RemoveTagParams) WithTag(tag string) *RemoveTagParams {
	o.SetTag(tag)
	return o
}

// SetTag adds the tag to the remove tag params
func (o *RemoveTagParams) SetTag(tag string) {
	o.Tag = tag
}

// WriteToRequest writes these params to a swagger request
func (o *RemoveTagParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param itemname
	if err := r.SetPathParam("itemname", o.Itemname); err != nil {
		return err
	}

	// path param tag
	if err := r.SetPathParam("tag", o.Tag); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}