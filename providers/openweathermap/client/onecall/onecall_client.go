// Code generated by go-swagger; DO NOT EDIT.

package onecall

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new onecall API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for onecall API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientService is the interface for Client methods
type ClientService interface {
	GetOneCall(params *GetOneCallParams, authInfo runtime.ClientAuthInfoWriter) (*GetOneCallOK, error)

	GetOneCallTimeMachine(params *GetOneCallTimeMachineParams, authInfo runtime.ClientAuthInfoWriter) (*GetOneCallTimeMachineOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  GetOneCall get one call API
*/
func (a *Client) GetOneCall(params *GetOneCallParams, authInfo runtime.ClientAuthInfoWriter) (*GetOneCallOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetOneCallParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getOneCall",
		Method:             "GET",
		PathPattern:        "/data/2.5/onecall?lat={lat}&lon={lon}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetOneCallReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetOneCallOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetOneCallDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  GetOneCallTimeMachine get one call time machine API
*/
func (a *Client) GetOneCallTimeMachine(params *GetOneCallTimeMachineParams, authInfo runtime.ClientAuthInfoWriter) (*GetOneCallTimeMachineOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetOneCallTimeMachineParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getOneCallTimeMachine",
		Method:             "GET",
		PathPattern:        "/data/2.5/onecall/timemachine?lat={lat}&lon={lon}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetOneCallTimeMachineReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetOneCallTimeMachineOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetOneCallTimeMachineDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
