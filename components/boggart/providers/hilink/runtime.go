package hilink

import (
	"sync"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type Runtime struct {
	*client.Runtime

	mutex          sync.RWMutex
	authentication runtime.ClientAuthInfoWriter
}

func newRuntime(original *client.Runtime) *Runtime {
	rt := &Runtime{
		Runtime: original,
	}

	rt.DefaultAuthentication = runtime.ClientAuthInfoWriterFunc(rt.CallDefaultAuthentication)

	return rt
}

func (r *Runtime) CallDefaultAuthentication(request runtime.ClientRequest, registry strfmt.Registry) error {
	r.mutex.RLock()
	auth := r.authentication
	r.mutex.RUnlock()

	if auth == nil {
		return nil
	}

	return auth.AuthenticateRequest(request, registry)
}

func (r *Runtime) SetDefaultAuthentication(auth runtime.ClientAuthInfoWriter) {
	r.mutex.Lock()
	r.authentication = auth
	r.mutex.Unlock()
}
