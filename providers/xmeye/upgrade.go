package xmeye

import (
	"context"
	"fmt"
	"io"
)

func (c *Client) UpgradeStart(ctx context.Context) error {
	_, err := c.Call(ctx, CmdUpgradeRequest, nil)
	return err
}

func (c *Client) UpgradeProgress(ctx context.Context) error {
	var result interface{}

	err := c.CallWithResult(ctx, CmdUpgradeRequest, nil, &result)

	fmt.Println(result)

	return err
}

func (c *Client) UpgradeData(ctx context.Context, data io.ReadSeeker) error {
	var (
		buf     = make([]byte, 1024*10)
		packets uint16
	)

	for {
		_, err := data.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}

			break
		}

		packets++
	}

	data.Seek(0, 0)

	packet := newPacket()
	packet.messageID = 1522
	packet.totalPacket = packets

	for i := uint16(1); i <= packets; i++ {
		n, err := data.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}

			break
		}

		packet.payloadLen = n
		if i == packets {
			packet.currentPacket = 1
		}

		if err := packet.LoadPayload(buf); err != nil {
			return err
		}

		if err := c.connection.send(packet); err != nil {
			return err
		}

		response, err := c.connection.receive()
		if err != nil {
			return err
		}

		var result Response
		if err = response.payload.JSONUnmarshal(&result); err != nil {
			return err
		}

		if result.Ret != CodeUpgradeData {
			fmt.Println(response.payload)
			break
		}
	}

	return nil
}

func (c *Client) UpgradeInfo(ctx context.Context) (*OPSystemUpgrade, error) {
	var result struct {
		Response
		OPSystemUpgrade OPSystemUpgrade
	}

	err := c.CallWithResult(ctx, CmdUpgradeInfoRequest, nil, &result)

	if err != nil {
		return nil, err
	}

	return &result.OPSystemUpgrade, nil
}
