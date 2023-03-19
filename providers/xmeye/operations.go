package xmeye

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
	"strconv"
	"time"
)

func (c *Client) OPTime(ctx context.Context) (*time.Time, error) {
	var result struct {
		Response
		OPTimeQuery string
	}

	err := c.CmdWithResult(ctx, CmdTimeRequest, "OPTimeQuery", &result)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(TimeLayout, result.OPTimeQuery)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *Client) OPTimeSetting(ctx context.Context, t time.Time) error {
	_, err := c.Call(ctx, CmdSysManagerRequest, map[string]interface{}{
		"Name":          "OPTimeSetting",
		"SessionID":     c.connection.SessionIDAsString(),
		"OPTimeSetting": t.Format(TimeLayout),
	})

	return err
}

func (c *Client) LogExport(ctx context.Context) (io.Reader, error) {
	packet, err := c.Call(ctx, CmdLogExportRequest, nil)

	return packet.payload, err
}

func (c *Client) Reboot(ctx context.Context) error {
	_, err := c.Call(ctx, CmdSysManagerRequest, map[string]interface{}{
		"Name":      "OPMachine",
		"SessionID": c.connection.SessionIDAsString(),
		"OPMachine": map[string]string{
			"Action": "Reboot",
		},
	})

	return err
}

func (c *Client) Snapshot(ctx context.Context, channel uint64) (io.Reader, error) {
	packet, err := c.Call(ctx, CmdNetSnapshotRequest, map[string]interface{}{
		"Name":      "OPSNAP",
		"SessionID": c.connection.SessionIDAsString(),
		"OPSNAP": map[string]uint64{
			"Channel": channel,
		},
	})

	return packet.payload, err
}

func (c *Client) SnapshotRTSP(ctx context.Context, channel uint64) (io.Reader, error) {
	cmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-loglevel", "error",
		"-rtsp_transport", "tcp",
		"-i", c.rtspURL.String()+"&channel="+strconv.FormatUint(channel, 10),
		"-f", "singlejpeg",
		"-vframes", "1",
		"-pix_fmt", "yuvj420p",
		"pipe:1",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			sc := bufio.NewScanner(&stderr)
			for sc.Scan() {
				return nil, errors.New(sc.Text())
			}
		}
	}

	return &stdout, err
}
