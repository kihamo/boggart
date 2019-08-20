package xmeye

import (
	"context"
	"sync/atomic"
	"time"
)

func (c *Client) Login(ctx context.Context) error {
	var result LoginResponse

	err := c.CallWithResult(ctx, CmdLoginResponse, map[string]string{
		"EncryptType": "MD5",
		"LoginType":   "DVRIP-Web",
		"PassWord":    HashPassword(c.password),
		"UserName":    c.username,
	}, &result)

	if err != nil {
		return err
	}

	atomic.StoreUint32(&c.sessionID, uint32(result.SessionID))

	c.mutex.Lock()
	if c.done != nil {
		close(c.done)
	}

	c.done = make(chan struct{}, 1)
	c.mutex.Unlock()

	if result.AliveInterval > 0 {
		go c.keepAlive(result.AliveInterval)
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
