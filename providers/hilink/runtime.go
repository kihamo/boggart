package hilink

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type Runtime struct {
	*client.Runtime

	mutex          sync.RWMutex
	authentication runtime.ClientAuthInfoWriter
	auth           func(ctx context.Context) error
}

func newRuntime(original *client.Runtime, auth func(ctx context.Context) error) *Runtime {
	rt := &Runtime{
		Runtime: original,
		auth:    auth,
	}

	rt.Jar, _ = cookiejar.New(nil)
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
		if req.GetPath() == "/api/webserver/SesTokInfo" {
			return nil
		}

		if err := r.auth(context.Background()); err != nil {
			return err
		}

		return r.callDefaultAuthentication(req, rg)
	}))
}

func (r *Runtime) SetAuthenticationLogged(token, session string) {
	u, _ := url.Parse("http://" + r.Host + r.BasePath)
	parts := strings.Split(session, "=")

	r.Jar.SetCookies(u, []*http.Cookie{
		{
			Name:     parts[0],
			Value:    parts[1],
			Path:     "/",
			HttpOnly: true,
		},
	})

	r.setDefaultAuthentication(runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) (err error) {
		return req.SetHeaderParam(headerToken, token)
	}))
}
