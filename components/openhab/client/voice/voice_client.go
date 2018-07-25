// Code generated by go-swagger; DO NOT EDIT.

package voice

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new voice API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for voice API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
GetInterpreter gets a single interpreters
*/
func (a *Client) GetInterpreter(params *GetInterpreterParams) (*GetInterpreterOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetInterpreterParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getInterpreter",
		Method:             "GET",
		PathPattern:        "/voice/interpreters/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetInterpreterReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetInterpreterOK), nil

}

/*
GetInterpreters gets the list of all interpreters
*/
func (a *Client) GetInterpreters(params *GetInterpretersParams) (*GetInterpretersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetInterpretersParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getInterpreters",
		Method:             "GET",
		PathPattern:        "/voice/interpreters",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetInterpretersReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetInterpretersOK), nil

}

/*
Interpret sends a text to a given human language interpreter
*/
func (a *Client) Interpret(params *InterpretParams) (*InterpretOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewInterpretParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "interpret",
		Method:             "POST",
		PathPattern:        "/voice/interpreters/{id}",
		ProducesMediaTypes: []string{""},
		ConsumesMediaTypes: []string{"text/plain"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &InterpretReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*InterpretOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}