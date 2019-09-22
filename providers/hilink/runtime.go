package hilink

import (
	"context"
	"sync"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type Runtime struct {
	*client.Runtime

	mutex          sync.RWMutex
	authentication runtime.ClientAuthInfoWriter

	token   string
	session string
	auth    func(ctx context.Context) error
}

func newRuntime(original *client.Runtime, auth func(ctx context.Context) error) *Runtime {
	rt := &Runtime{
		Runtime: original,
		auth:    auth,
	}

	rt.DefaultAuthentication = runtime.ClientAuthInfoWriterFunc(rt.callDefaultAuthentication)
	rt.SetAuthenticationAnonymous()

	return rt
}

func (r *Runtime) callDefaultAuthentication(request runtime.ClientRequest, registry strfmt.Registry) error {
	r.mutex.RLock()
	auth := r.authentication
	r.mutex.RUnlock()

	if auth == nil {
		return nil
	}

	return auth.AuthenticateRequest(request, registry)
}

func (r *Runtime) setDefaultAuthentication(auth runtime.ClientAuthInfoWriter) {
	r.mutex.Lock()
	r.authentication = auth
	r.mutex.Unlock()
}

func (r *Runtime) SetAuthenticationAnonymous() {
	r.setDefaultAuthentication(runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, rg strfmt.Registry) error {
		if req.GetPath() == "/webserver/SesTokInfo" {
			return nil
		}

		if err := r.auth(context.Background()); err != nil {
			return err
		}

		return r.callDefaultAuthentication(req, rg)
	}))
}

func (r *Runtime) SetAuthenticationLogged(token, session string) {
	r.mutex.RLock()
	t, s := r.token, r.session
	r.mutex.RUnlock()

	if t == token && s == session {
		return
	}

	r.setDefaultAuthentication(runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) (err error) {
		err = req.SetHeaderParam(headerToken, token)
		if err == nil {
			err = req.SetHeaderParam(headerSession, session)
		}

		return err
	}))
}
