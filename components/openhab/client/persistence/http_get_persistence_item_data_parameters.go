// Code generated by go-swagger; DO NOT EDIT.

package persistence

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewHTTPGetPersistenceItemDataParams creates a new HTTPGetPersistenceItemDataParams object
// with the default values initialized.
func NewHTTPGetPersistenceItemDataParams() *HTTPGetPersistenceItemDataParams {
	var ()
	return &HTTPGetPersistenceItemDataParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewHTTPGetPersistenceItemDataParamsWithTimeout creates a new HTTPGetPersistenceItemDataParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewHTTPGetPersistenceItemDataParamsWithTimeout(timeout time.Duration) *HTTPGetPersistenceItemDataParams {
	var ()
	return &HTTPGetPersistenceItemDataParams{

		timeout: timeout,
	}
}

// NewHTTPGetPersistenceItemDataParamsWithContext creates a new HTTPGetPersistenceItemDataParams object
// with the default values initialized, and the ability to set a context for a request
func NewHTTPGetPersistenceItemDataParamsWithContext(ctx context.Context) *HTTPGetPersistenceItemDataParams {
	var ()
	return &HTTPGetPersistenceItemDataParams{

		Context: ctx,
	}
}

// NewHTTPGetPersistenceItemDataParamsWithHTTPClient creates a new HTTPGetPersistenceItemDataParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewHTTPGetPersistenceItemDataParamsWithHTTPClient(client *http.Client) *HTTPGetPersistenceItemDataParams {
	var ()
	return &HTTPGetPersistenceItemDataParams{
		HTTPClient: client,
	}
}

/*HTTPGetPersistenceItemDataParams contains all the parameters to send to the API endpoint
for the http get persistence item data operation typically these are written to a http.Request
*/
type HTTPGetPersistenceItemDataParams struct {

	/*Boundary
	  Gets one value before and after the requested period.

	*/
	Boundary *bool
	/*Endtime
	  End time of the data to return. Will default to current time. [yyyy-MM-dd'T'HH:mm:ss.SSSZ]

	*/
	Endtime *string
	/*Itemname
	  The item name

	*/
	Itemname string
	/*Page
	  Page number of data to return. This parameter will enable paging.

	*/
	Page *int32
	/*Pagelength
	  The length of each page.

	*/
	Pagelength *int32
	/*ServiceID
	  Id of the persistence service. If not provided the default service will be used

	*/
	ServiceID *string
	/*Starttime
	  Start time of the data to return. Will default to 1 day before endtime. [yyyy-MM-dd'T'HH:mm:ss.SSSZ]

	*/
	Starttime *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) WithTimeout(timeout time.Duration) *HTTPGetPersistenceItemDataParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) WithContext(ctx context.Context) *HTTPGetPersistenceItemDataParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) WithHTTPClient(client *http.Client) *HTTPGetPersistenceItemDataParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBoundary adds the boundary to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) WithBoundary(boundary *bool) *HTTPGetPersistenceItemDataParams {
	o.SetBoundary(boundary)
	return o
}

// SetBoundary adds the boundary to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) SetBoundary(boundary *bool) {
	o.Boundary = boundary
}

// WithEndtime adds the endtime to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) WithEndtime(endtime *string) *HTTPGetPersistenceItemDataParams {
	o.SetEndtime(endtime)
	return o
}

// SetEndtime adds the endtime to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) SetEndtime(endtime *string) {
	o.Endtime = endtime
}

// WithItemname adds the itemname to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) WithItemname(itemname string) *HTTPGetPersistenceItemDataParams {
	o.SetItemname(itemname)
	return o
}

// SetItemname adds the itemname to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) SetItemname(itemname string) {
	o.Itemname = itemname
}

// WithPage adds the page to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) WithPage(page *int32) *HTTPGetPersistenceItemDataParams {
	o.SetPage(page)
	return o
}

// SetPage adds the page to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) SetPage(page *int32) {
	o.Page = page
}

// WithPagelength adds the pagelength to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) WithPagelength(pagelength *int32) *HTTPGetPersistenceItemDataParams {
	o.SetPagelength(pagelength)
	return o
}

// SetPagelength adds the pagelength to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) SetPagelength(pagelength *int32) {
	o.Pagelength = pagelength
}

// WithServiceID adds the serviceID to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) WithServiceID(serviceID *string) *HTTPGetPersistenceItemDataParams {
	o.SetServiceID(serviceID)
	return o
}

// SetServiceID adds the serviceId to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) SetServiceID(serviceID *string) {
	o.ServiceID = serviceID
}

// WithStarttime adds the starttime to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) WithStarttime(starttime *string) *HTTPGetPersistenceItemDataParams {
	o.SetStarttime(starttime)
	return o
}

// SetStarttime adds the starttime to the http get persistence item data params
func (o *HTTPGetPersistenceItemDataParams) SetStarttime(starttime *string) {
	o.Starttime = starttime
}

// WriteToRequest writes these params to a swagger request
func (o *HTTPGetPersistenceItemDataParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Boundary != nil {

		// query param boundary
		var qrBoundary bool
		if o.Boundary != nil {
			qrBoundary = *o.Boundary
		}
		qBoundary := swag.FormatBool(qrBoundary)
		if qBoundary != "" {
			if err := r.SetQueryParam("boundary", qBoundary); err != nil {
				return err
			}
		}

	}

	if o.Endtime != nil {

		// query param endtime
		var qrEndtime string
		if o.Endtime != nil {
			qrEndtime = *o.Endtime
		}
		qEndtime := qrEndtime
		if qEndtime != "" {
			if err := r.SetQueryParam("endtime", qEndtime); err != nil {
				return err
			}
		}

	}

	// path param itemname
	if err := r.SetPathParam("itemname", o.Itemname); err != nil {
		return err
	}

	if o.Page != nil {

		// query param page
		var qrPage int32
		if o.Page != nil {
			qrPage = *o.Page
		}
		qPage := swag.FormatInt32(qrPage)
		if qPage != "" {
			if err := r.SetQueryParam("page", qPage); err != nil {
				return err
			}
		}

	}

	if o.Pagelength != nil {

		// query param pagelength
		var qrPagelength int32
		if o.Pagelength != nil {
			qrPagelength = *o.Pagelength
		}
		qPagelength := swag.FormatInt32(qrPagelength)
		if qPagelength != "" {
			if err := r.SetQueryParam("pagelength", qPagelength); err != nil {
				return err
			}
		}

	}

	if o.ServiceID != nil {

		// query param serviceId
		var qrServiceID string
		if o.ServiceID != nil {
			qrServiceID = *o.ServiceID
		}
		qServiceID := qrServiceID
		if qServiceID != "" {
			if err := r.SetQueryParam("serviceId", qServiceID); err != nil {
				return err
			}
		}

	}

	if o.Starttime != nil {

		// query param starttime
		var qrStarttime string
		if o.Starttime != nil {
			qrStarttime = *o.Starttime
		}
		qStarttime := qrStarttime
		if qStarttime != "" {
			if err := r.SetQueryParam("starttime", qStarttime); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}