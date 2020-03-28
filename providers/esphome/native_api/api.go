package nativeapi

import (
	"context"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

func (c *Client) Hello(ctx context.Context) (*HelloResponse, error) {
	out, err := c.invoke(ctx, &HelloRequest{ClientInfo: c.ClientID()}, "native_api.HelloResponse")
	if err != nil {
		return nil, err
	}

	return out.(*HelloResponse), nil
}

func (c *Client) Connect(ctx context.Context, password string) (*ConnectResponse, error) {
	out, err := c.invoke(ctx, &ConnectRequest{Password: password}, "native_api.ConnectResponse")
	if err != nil {
		return nil, err
	}

	return out.(*ConnectResponse), nil
}

func (c *Client) Disconnect(ctx context.Context) (*DisconnectResponse, error) {
	out, err := c.invoke(ctx, &DisconnectRequest{}, "native_api.DisconnectResponse")
	if err != nil {
		return nil, err
	}

	return out.(*DisconnectResponse), nil
}

func (c *Client) Ping(ctx context.Context) (*PingResponse, error) {
	out, err := c.invoke(ctx, &PingRequest{}, "native_api.PingResponse")
	if err != nil {
		return nil, err
	}
	return out.(*PingResponse), nil
}

func (c *Client) DeviceInfo(ctx context.Context) (*DeviceInfoResponse, error) {
	out, err := c.invoke(ctx, &DeviceInfoRequest{}, "native_api.DeviceInfoResponse")
	if err != nil {
		return nil, err
	}
	return out.(*DeviceInfoResponse), nil
}

func (c *Client) ListEntities(ctx context.Context) (_ []proto.Message, err error) {
	chMessages := make(chan proto.Message, 1)
	chDone := make(chan error, 1)

	var (
		canceled uint32
		cancel   context.CancelFunc
	)

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	err = c.invokeHandler(ctx, &ListEntitiesRequest{}, func(message proto.Message, err error) bool {
		if atomic.LoadUint32(&canceled) != 0 {
			return true
		}

		if err != nil {
			chDone <- err
			return true
		}

		switch proto.MessageName(message) {
		case "native_api.ListEntitiesBinarySensorResponse",
			"native_api.ListEntitiesCoverResponse",
			"native_api.ListEntitiesFanResponse",
			"native_api.ListEntitiesLightResponse",
			"native_api.ListEntitiesSensorResponse",
			"native_api.ListEntitiesSwitchResponse",
			"native_api.ListEntitiesTextSensorResponse",
			"native_api.ListEntitiesServicesResponse",
			"native_api.ListEntitiesCameraResponse",
			"native_api.ListEntitiesClimateResponse":
			chMessages <- message

		case "native_api.ListEntitiesDoneResponse":
			close(chDone)
			return true
		}

		return false
	})

	if err != nil {
		return nil, err
	}

	messages := make([]proto.Message, 0)

	for {
		select {
		case <-ctx.Done():
			atomic.StoreUint32(&canceled, 1)
			return nil, ctx.Err()

		case err := <-chDone:
			if err != nil {
				return nil, err
			}

			return messages, nil

		case message := <-chMessages:
			messages = append(messages, message)
		}
	}
}

func (c *Client) SubscribeStates(ctx context.Context) (_ <-chan proto.Message, err error) {
	return c.subscribe(ctx,
		&SubscribeStatesRequest{},
		func(message proto.Message) bool {
			switch proto.MessageName(message) {
			case "native_api.BinarySensorStateResponse",
				"native_api.CoverStateResponse",
				"native_api.FanStateResponse",
				"native_api.LightStateResponse",
				"native_api.SensorStateResponse",
				"native_api.SwitchStateResponse",
				"native_api.TextSensorStateResponse",
				"native_api.SubscribeHomeAssistantStateResponse",
				"native_api.HomeAssistantStateResponse",
				"native_api.ClimateStateResponse":
				return true
			}

			return false
		})
}

func (c *Client) SubscribeLogs(ctx context.Context, logLevel LogLevel, dumpConfig bool) (_ <-chan proto.Message, err error) {
	return c.subscribe(ctx,
		&SubscribeLogsRequest{Level: logLevel, DumpConfig: dumpConfig},
		func(message proto.Message) bool {
			return proto.MessageName(message) == "native_api.SubscribeLogsResponse"
		})
}

func (c *Client) SubscribeHomeassistantServices(ctx context.Context) (*HomeassistantServiceResponse, error) {
	if err := c.authenticateCheck(); err != nil {
		return nil, err
	}

	out, err := c.invoke(ctx, &SubscribeHomeassistantServicesRequest{}, "native_api.HomeassistantServiceResponse")
	if err != nil {
		return nil, err
	}
	return out.(*HomeassistantServiceResponse), nil
}

func (c *Client) SubscribeHomeAssistantStates(ctx context.Context) (*SubscribeHomeAssistantStateResponse, error) {
	if err := c.authenticateCheck(); err != nil {
		return nil, err
	}

	out, err := c.invoke(ctx, &SubscribeHomeAssistantStatesRequest{}, "native_api.SubscribeHomeAssistantStateResponse")
	if err != nil {
		return nil, err
	}
	return out.(*SubscribeHomeAssistantStateResponse), nil
}

func (c *Client) GetTime(ctx context.Context) (*GetTimeResponse, error) {
	out, err := c.invoke(ctx, &GetTimeRequest{}, "native_api.GetTimeResponse")
	if err != nil {
		return nil, err
	}
	return out.(*GetTimeResponse), nil
}

func (c *Client) ExecuteService(ctx context.Context, key uint32, args []*ExecuteServiceArgument) error {
	if err := c.authenticateCheck(); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, &ExecuteServiceRequest{Key: key, Args: args})
}

// nolint:interfacer
func (c *Client) CoverCommand(ctx context.Context, in *CoverCommandRequest) error {
	if err := c.authenticateCheck(); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, in)
}

// nolint:interfacer
func (c *Client) FanCommand(ctx context.Context, in *FanCommandRequest) error {
	if err := c.authenticateCheck(); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, in)
}

// nolint:interfacer
func (c *Client) LightCommand(ctx context.Context, in *LightCommandRequest) error {
	if err := c.authenticateCheck(); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, in)
}

// nolint:interfacer
func (c *Client) SwitchCommand(ctx context.Context, in *SwitchCommandRequest) error {
	if err := c.authenticateCheck(); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, in)
}

// nolint:interfacer
func (c *Client) CameraImage(ctx context.Context, in *CameraImageRequest) (*CameraImageResponse, error) {
	if err := c.authenticateCheck(); err != nil {
		return nil, err
	}

	out, err := c.invoke(ctx, in, "native_api.CameraImageResponse")
	if err != nil {
		return nil, err
	}
	return out.(*CameraImageResponse), nil
}

// nolint:interfacer
func (c *Client) ClimateCommand(ctx context.Context, in *ClimateCommandRequest) error {
	if err := c.authenticateCheck(); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, in)
}
