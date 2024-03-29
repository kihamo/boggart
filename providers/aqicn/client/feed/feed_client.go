// Code generated by go-swagger; DO NOT EDIT.

package feed

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new feed API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for feed API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
GetByCity get by city API
*/
func (a *Client) GetByCity(params *GetByCityParams, authInfo runtime.ClientAuthInfoWriter) (*GetByCityOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetByCityParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getByCity",
		Method:             "GET",
		PathPattern:        "/feed/{city}/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetByCityReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetByCityOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetByCityDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
GetByIP get by IP API
*/
func (a *Client) GetByIP(params *GetByIPParams, authInfo runtime.ClientAuthInfoWriter) (*GetByIPOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetByIPParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getByIP",
		Method:             "GET",
		PathPattern:        "/feed/here/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetByIPReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetByIPOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetByIPDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
GetByLatLng get by lat lng API
*/
func (a *Client) GetByLatLng(params *GetByLatLngParams, authInfo runtime.ClientAuthInfoWriter) (*GetByLatLngOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetByLatLngParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getByLatLng",
		Method:             "GET",
		PathPattern:        "/feed/geo:{lat};{lng}/",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetByLatLngReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetByLatLngOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetByLatLngDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
