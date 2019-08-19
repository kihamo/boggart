package xmeye

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xmeye/internal"
)

func (c *Client) Login(ctx context.Context) error {
	response := &internal.LoginResponse{}

	err := c.CallWithResult(ctx, CmdLoginResponse, map[string]string{
		"EncryptType": "MD5",
		"LoginType":   "DVRIP-Web",
		"PassWord":    HashPassword(c.password),
		"UserName":    c.username,
	}, response)

	if err != nil {
		return err
	}

	session, err := hex.DecodeString(response.SessionID[2:])
	if err != nil {
		return err
	}

	sessionID := binary.LittleEndian.Uint32([]byte{session[3], session[2], session[1], session[0]})
	atomic.StoreUint32(&c.sessionID, sessionID)

	c.mutex.Lock()
	if c.done != nil {
		close(c.done)
	}

	c.done = make(chan struct{}, 1)
	c.mutex.Unlock()

	if response.AliveInterval > 0 {
		go c.keepAlive(response.AliveInterval)
	}

	return err
}

func (c *Client) Logout(ctx context.Context) error {
	_, err := c.Call(ctx, CmdLogoutResponse, nil)

	if err != nil {
		return c.Close()
	}

	atomic.StoreUint32(&c.sessionID, 0)
	return err
}

func (c *Client) keepAlive(interval uint64) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))

	c.mutex.RLock()
	done := c.done
	c.mutex.RUnlock()

	for {
		select {
		case <-ticker.C:
			c.Cmd(context.Background(), CmdKeepAliveResponse, "KeepAlive")

		case <-done:
			ticker.Stop()
			return
		}
	}
}
