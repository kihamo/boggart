package swagger

import (
	"sync"

	"github.com/go-openapi/runtime/client"
	l "github.com/go-openapi/runtime/logger"
)

var (
	lockDebug  sync.Mutex
	lockLogger sync.Mutex
)

func SetDebug(r *client.Runtime, debug bool) {
	lockDebug.Lock()
	defer lockDebug.Unlock()

	r.SetDebug(debug)
}

func SetLogger(r *client.Runtime, logger l.Logger) {
	lockLogger.Lock()
	defer lockLogger.Unlock()

	r.SetLogger(logger)
}
