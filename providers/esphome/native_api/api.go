package native_api

import (
	"context"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

func (c *Client) Hello(ctx context.Context) (*HelloResponse, error) {
	out, err := c.invoke(ctx, &HelloRequest{ClientInfo: c.ClientID()})
	if err != nil {
		return nil, err
	}

	return out.(*HelloResponse), nil
}

func (c *Client) Connect(ctx context.Context) (*ConnectResponse, error) {
	out, err := c.invoke(ctx, &ConnectRequest{Password: c.password})
	if err != nil {
		return nil, err
	}

	return out.(*ConnectResponse), nil
}

func (c *Client) Disconnect(ctx context.Context) (*DisconnectResponse, error) {
	out, err := c.invoke(ctx, &DisconnectRequest{})
	if err != nil {
		return nil, err
	}

	return out.(*DisconnectResponse), nil
}

func (c *Client) Ping(ctx context.Context) (*PingResponse, error) {
	out, err := c.invoke(ctx, &PingRequest{})
	if err != nil {
		return nil, err
	}
	return out.(*PingResponse), nil
}

func (c *Client) DeviceInfo(ctx context.Context) (*DeviceInfoResponse, error) {
	out, err := c.invoke(ctx, &DeviceInfoRequest{})
	if err != nil {
		return nil, err
	}
	return out.(*DeviceInfoResponse), nil
}

func (c *Client) ListEntities(ctx context.Context) ([]proto.Message, error) {
	if err := c.checkAuthenticate(ctx); err != nil {
		return nil, err
	}

	messages := make([]proto.Message, 0)

	message, err := c.invoke(ctx, &ListEntitiesRequest{})
	if err != nil {
		return nil, err
	}

	messages = append(messages, message)

	for {
		message, err = c.connection.Read(ctx)
		if err != nil {
			return messages, err
		}

		if strings.Compare(proto.MessageName(message), "native_api.ListEntitiesDoneResponse") == 0 {
			return messages, nil
		}
	}

	return messages, nil
}

func (c *Client) SubscribeStates(ctx context.Context) (<-chan proto.Message, <-chan error) {
	client, err := c.Clone()
	if err == nil {
		client.connectionPing = false

		err = client.checkAuthenticate(ctx)
	}

	if err != nil {
		chMessage := make(chan proto.Message)
		chError := make(chan error)

		go func() {
			chError <- err

			close(chMessage)
			close(chError)
		}()

		return chMessage, chError
	}

	subscribe := NewSubscribe(client, ctx, &SubscribeStatesRequest{}, time.Second)
	return subscribe.NextMessage(), subscribe.NextError()
}

func (c *Client) SubscribeLogs(ctx context.Context, logLevel LogLevel, dumpConfig bool) (<-chan proto.Message, <-chan error) {
	client, err := c.Clone()
	if err == nil {
		client.connectionPing = false

		err = client.checkAuthenticate(ctx)
	}

	if err != nil {
		chMessage := make(chan proto.Message)
		chError := make(chan error)

		go func() {
			chError <- err

			close(chMessage)
			close(chError)
		}()

		return chMessage, chError
	}

	subscribe := NewSubscribe(client, ctx, &SubscribeLogsRequest{Level: logLevel, DumpConfig: dumpConfig}, time.Second)
	return subscribe.NextMessage(), subscribe.NextError()
}

func (c *Client) SubscribeHomeassistantServices(ctx context.Context) (*HomeassistantServiceResponse, error) {
	if err := c.checkAuthenticate(ctx); err != nil {
		return nil, err
	}

	out, err := c.invoke(ctx, &SubscribeHomeassistantServicesRequest{})
	if err != nil {
		return nil, err
	}
	return out.(*HomeassistantServiceResponse), nil
}

func (c *Client) SubscribeHomeAssistantStates(ctx context.Context) (*SubscribeHomeAssistantStateResponse, error) {
	if err := c.checkAuthenticate(ctx); err != nil {
		return nil, err
	}

	out, err := c.invoke(ctx, &SubscribeHomeAssistantStatesRequest{})
	if err != nil {
		return nil, err
	}
	return out.(*SubscribeHomeAssistantStateResponse), nil
}

func (c *Client) GetTime(ctx context.Context) (*GetTimeResponse, error) {
	out, err := c.invoke(ctx, &GetTimeRequest{})
	if err != nil {
		return nil, err
	}
	return out.(*GetTimeResponse), nil
}

func (c *Client) ExecuteService(ctx context.Context, key uint32, args []*ExecuteServiceArgument) error {
	if err := c.checkAuthenticate(ctx); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, &ExecuteServiceRequest{Key: key, Args: args})
}

func (c *Client) CoverCommand(ctx context.Context, in *CoverCommandRequest) error {
	if err := c.checkAuthenticate(ctx); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, in)
}

func (c *Client) FanCommand(ctx context.Context, in *FanCommandRequest) error {
	if err := c.checkAuthenticate(ctx); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, in)
}

func (c *Client) LightCommand(ctx context.Context, in *LightCommandRequest) error {
	if err := c.checkAuthenticate(ctx); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, in)
}

func (c *Client) SwitchCommand(ctx context.Context, in *SwitchCommandRequest) error {
	if err := c.checkAuthenticate(ctx); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, in)
}

func (c *Client) CameraImage(ctx context.Context, in *CameraImageRequest) (*CameraImageResponse, error) {
	if err := c.checkAuthenticate(ctx); err != nil {
		return nil, err
	}

	out, err := c.invoke(ctx, in)
	if err != nil {
		return nil, err
	}
	return out.(*CameraImageResponse), nil
}

func (c *Client) ClimateCommand(ctx context.Context, in *ClimateCommandRequest) error {
	if err := c.checkAuthenticate(ctx); err != nil {
		return err
	}

	return c.invokeNoDelay(ctx, in)
}
