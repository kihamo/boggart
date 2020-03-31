// Code generated by go-swagger; DO NOT EDIT.

package voice

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

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
GetDefaultVoice gets the default voice
*/
func (a *Client) GetDefaultVoice(params *GetDefaultVoiceParams) (*GetDefaultVoiceOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetDefaultVoiceParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getDefaultVoice",
		Method:             "GET",
		PathPattern:        "/voice/defaultvoice",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetDefaultVoiceReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetDefaultVoiceOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getDefaultVoice: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetInterpreter gets a single interpreter
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
	success, ok := result.(*GetInterpreterOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getInterpreter: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
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
	success, ok := result.(*GetInterpretersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getInterpreters: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
GetVoices gets the list of all voices
*/
func (a *Client) GetVoices(params *GetVoicesParams) (*GetVoicesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetVoicesParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getVoices",
		Method:             "GET",
		PathPattern:        "/voice/voices",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{""},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetVoicesReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetVoicesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getVoices: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
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
	success, ok := result.(*InterpretOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for interpret: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
Say speaks a given text with a given voice through the given audio sink
*/
func (a *Client) Say(params *SayParams) (*SayOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSayParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "say",
		Method:             "POST",
		PathPattern:        "/voice/say",
		ProducesMediaTypes: []string{""},
		ConsumesMediaTypes: []string{"text/plain"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &SayReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*SayOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for say: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}