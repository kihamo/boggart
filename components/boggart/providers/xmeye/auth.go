package xmeye

import (
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xmeye/internal"
)

func (c *Client) Login() error {
	response := &internal.LoginResponse{}

	err := c.CallWithResult(CmdLoginResponse, map[string]string{
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

func (c *Client) Logout() error {
	_, err := c.Call(CmdLogoutResponse, nil)

	if err != nil {
		return c.Close()
	}

	atomic.StoreUint32(&c.sessionID, 0)
	return err
}

func (c *Client) keepAlive(interval uint64) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))

	for {
		select {
		case <-ticker.C:
			c.Cmd(CmdKeepAliveResponse, "KeepAlive")

		case <-c.done:
			ticker.Stop()
			return
		}
	}
}
