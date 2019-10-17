// Code generated by go-swagger; DO NOT EDIT.

package system

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new system API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for system API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
ExecuteCommand executes a registered system command
*/
func (a *Client) ExecuteCommand(params *ExecuteCommandParams, authInfo runtime.ClientAuthInfoWriter) (*ExecuteCommandNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewExecuteCommandParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "executeCommand",
		Method:             "POST",
		PathPattern:        "/system/commands/{source}/{action}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ExecuteCommandReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ExecuteCommandNoContent)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for executeCommand: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetCommands lists all registered system commands
*/
func (a *Client) GetCommands(params *GetCommandsParams, authInfo runtime.ClientAuthInfoWriter) (*GetCommandsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetCommandsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getCommands",
		Method:             "GET",
		PathPattern:        "/system/commands",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetCommandsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetCommandsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getCommands: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetCommandsBySource lists all registered system commands for a source
*/
func (a *Client) GetCommandsBySource(params *GetCommandsBySourceParams, authInfo runtime.ClientAuthInfoWriter) (*GetCommandsBySourceOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetCommandsBySourceParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getCommandsBySource",
		Method:             "GET",
		PathPattern:        "/system/commands/{source}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetCommandsBySourceReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetCommandsBySourceOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getCommandsBySource: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
