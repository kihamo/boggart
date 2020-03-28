package native_api

import (
	"context"
	"errors"
	"sync/atomic"
)

const (
	authNone uint32 = iota
	authSuccess
	authFailed
)

var (
	ErrAuthenticated      = errors.New("must authenticated")
	ErrAuthenticateFailed = errors.New("authenticated failed")
)

func (c *Client) authenticateRun() (err error) {
	if err = c.connectionCheck(); err != nil {
		return err
	}

	if atomic.LoadUint32(&c.authenticated) == authNone {
		err = c.write(context.Background(), &ConnectRequest{Password: c.password})

		if err == nil {
			atomic.StoreUint32(&c.authenticated, authSuccess)
		} else {
			atomic.StoreUint32(&c.authenticated, authFailed)
		}
	}

	return err
}

func (c *Client) authenticateClose() {
	atomic.StoreUint32(&c.authenticated, authNone)
}

func (c *Client) authenticateCheck() (err error) {
	switch atomic.LoadUint32(&c.authenticated) {
	case authNone:
		err = ErrAuthenticated

	case authFailed:
		err = ErrAuthenticateFailed
	}

	return err
}
