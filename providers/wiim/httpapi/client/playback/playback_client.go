// Code generated by go-swagger; DO NOT EDIT.

package playback

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new playback API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for playback API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	GetPlayerStatus(params *GetPlayerStatusParams, opts ...ClientOption) (*GetPlayerStatusOK, error)

	SetPlayerCmd(params *SetPlayerCmdParams, opts ...ClientOption) (*SetPlayerCmdOK, error)

	SetPlayerCmdControl(params *SetPlayerCmdControlParams, opts ...ClientOption) (*SetPlayerCmdControlOK, error)

	SetPlayerCmdLoopmode(params *SetPlayerCmdLoopmodeParams, opts ...ClientOption) (*SetPlayerCmdLoopmodeOK, error)

	SetPlayerCmdMute(params *SetPlayerCmdMuteParams, opts ...ClientOption) (*SetPlayerCmdMuteOK, error)

	SetPlayerCmdPlaylist(params *SetPlayerCmdPlaylistParams, opts ...ClientOption) (*SetPlayerCmdPlaylistOK, error)

	SetPlayerCmdPlaylistHEX(params *SetPlayerCmdPlaylistHEXParams, opts ...ClientOption) (*SetPlayerCmdPlaylistHEXOK, error)

	SetPlayerCmdSeek(params *SetPlayerCmdSeekParams, opts ...ClientOption) (*SetPlayerCmdSeekOK, error)

	SetPlayerCmdSwitchMode(params *SetPlayerCmdSwitchModeParams, opts ...ClientOption) (*SetPlayerCmdSwitchModeOK, error)

	SetPlayerCmdVolume(params *SetPlayerCmdVolumeParams, opts ...ClientOption) (*SetPlayerCmdVolumeOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
GetPlayerStatus Get the playback status
*/
func (a *Client) GetPlayerStatus(params *GetPlayerStatusParams, opts ...ClientOption) (*GetPlayerStatusOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetPlayerStatusParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getPlayerStatus",
		Method:             "GET",
		PathPattern:        "/httpapi.asp?command=getPlayerStatus",
		ProducesMediaTypes: []string{"text/html"},
		ConsumesMediaTypes: []string{"text/html"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetPlayerStatusReader{formats: a.formats},
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
	success, ok := result.(*GetPlayerStatusOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getPlayerStatus: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SetPlayerCmd Play Instruction for any valid audio file or stream specified as a URL.
*/
func (a *Client) SetPlayerCmd(params *SetPlayerCmdParams, opts ...ClientOption) (*SetPlayerCmdOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPlayerCmdParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "setPlayerCmd",
		Method:             "GET",
		PathPattern:        "/httpapi.asp?command=setPlayerCmd:play:{url}",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"text/html"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &SetPlayerCmdReader{formats: a.formats},
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
	success, ok := result.(*SetPlayerCmdOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for setPlayerCmd: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SetPlayerCmdControl Control the current playback
*/
func (a *Client) SetPlayerCmdControl(params *SetPlayerCmdControlParams, opts ...ClientOption) (*SetPlayerCmdControlOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPlayerCmdControlParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "setPlayerCmdControl",
		Method:             "GET",
		PathPattern:        "/httpapi.asp?command=setPlayerCmd:{control}",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"text/html"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &SetPlayerCmdControlReader{formats: a.formats},
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
	success, ok := result.(*SetPlayerCmdControlOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for setPlayerCmdControl: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SetPlayerCmdLoopmode Set shuffle and repeat mode
*/
func (a *Client) SetPlayerCmdLoopmode(params *SetPlayerCmdLoopmodeParams, opts ...ClientOption) (*SetPlayerCmdLoopmodeOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPlayerCmdLoopmodeParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "setPlayerCmdLoopmode",
		Method:             "GET",
		PathPattern:        "/httpapi.asp?command=setPlayerCmd:loopmode:{mode}",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"text/html"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &SetPlayerCmdLoopmodeReader{formats: a.formats},
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
	success, ok := result.(*SetPlayerCmdLoopmodeOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for setPlayerCmdLoopmode: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SetPlayerCmdMute Toggle mute for the device
*/
func (a *Client) SetPlayerCmdMute(params *SetPlayerCmdMuteParams, opts ...ClientOption) (*SetPlayerCmdMuteOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPlayerCmdMuteParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "setPlayerCmdMute",
		Method:             "GET",
		PathPattern:        "/httpapi.asp?command=setPlayerCmd:mute:{mute}",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"text/html"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &SetPlayerCmdMuteReader{formats: a.formats},
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
	success, ok := result.(*SetPlayerCmdMuteOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for setPlayerCmdMute: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SetPlayerCmdPlaylist Play the URl from m3u or ASX playlist
*/
func (a *Client) SetPlayerCmdPlaylist(params *SetPlayerCmdPlaylistParams, opts ...ClientOption) (*SetPlayerCmdPlaylistOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPlayerCmdPlaylistParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "setPlayerCmdPlaylist",
		Method:             "GET",
		PathPattern:        "/httpapi.asp?command=setPlayerCmd:playlist:{url}:{index}",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"text/html"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &SetPlayerCmdPlaylistReader{formats: a.formats},
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
	success, ok := result.(*SetPlayerCmdPlaylistOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for setPlayerCmdPlaylist: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SetPlayerCmdPlaylistHEX Play the URl from m3u or ASX playlist
*/
func (a *Client) SetPlayerCmdPlaylistHEX(params *SetPlayerCmdPlaylistHEXParams, opts ...ClientOption) (*SetPlayerCmdPlaylistHEXOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPlayerCmdPlaylistHEXParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "setPlayerCmdPlaylistHEX",
		Method:             "GET",
		PathPattern:        "/httpapi.asp?command=setPlayerCmd:hex_playlist:{url}:{index}",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"text/html"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &SetPlayerCmdPlaylistHEXReader{formats: a.formats},
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
	success, ok := result.(*SetPlayerCmdPlaylistHEXOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for setPlayerCmdPlaylistHEX: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SetPlayerCmdSeek Seek with seconds for current playback, have no use when playing radio link.
*/
func (a *Client) SetPlayerCmdSeek(params *SetPlayerCmdSeekParams, opts ...ClientOption) (*SetPlayerCmdSeekOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPlayerCmdSeekParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "setPlayerCmdSeek",
		Method:             "GET",
		PathPattern:        "/httpapi.asp?command=setPlayerCmd:seek:{position}",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"text/html"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &SetPlayerCmdSeekReader{formats: a.formats},
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
	success, ok := result.(*SetPlayerCmdSeekOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for setPlayerCmdSeek: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SetPlayerCmdSwitchMode Selects the Audio Source of the Device. The available audio sources for each device will depend on the installed hardware.
*/
func (a *Client) SetPlayerCmdSwitchMode(params *SetPlayerCmdSwitchModeParams, opts ...ClientOption) (*SetPlayerCmdSwitchModeOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPlayerCmdSwitchModeParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "setPlayerCmdSwitchMode",
		Method:             "GET",
		PathPattern:        "/httpapi.asp?command=setPlayerCmd:switchmode:{mode}",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"text/html"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &SetPlayerCmdSwitchModeReader{formats: a.formats},
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
	success, ok := result.(*SetPlayerCmdSwitchModeOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for setPlayerCmdSwitchMode: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
SetPlayerCmdVolume Set system volume
*/
func (a *Client) SetPlayerCmdVolume(params *SetPlayerCmdVolumeParams, opts ...ClientOption) (*SetPlayerCmdVolumeOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewSetPlayerCmdVolumeParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "setPlayerCmdVolume",
		Method:             "GET",
		PathPattern:        "/httpapi.asp?command=setPlayerCmd:vol{volume}",
		ProducesMediaTypes: []string{"text/plain"},
		ConsumesMediaTypes: []string{"text/html"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &SetPlayerCmdVolumeReader{formats: a.formats},
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
	success, ok := result.(*SetPlayerCmdVolumeOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for setPlayerCmdVolume: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
