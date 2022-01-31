// Code generated by go-swagger; DO NOT EDIT.

package plugin

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new plugin API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for plugin API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	DisplayLayerProgress(params *DisplayLayerProgressParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DisplayLayerProgressOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  DisplayLayerProgress receives the layer height and other values
*/
func (a *Client) DisplayLayerProgress(params *DisplayLayerProgressParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*DisplayLayerProgressOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDisplayLayerProgressParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "displayLayerProgress",
		Method:             "GET",
		PathPattern:        "/plugin/DisplayLayerProgress/values",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DisplayLayerProgressReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DisplayLayerProgressOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for displayLayerProgress: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
