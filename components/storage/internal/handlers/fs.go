package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/shadow/components/dashboard"
)

const (
	routePrefix = "/" + storage.ComponentName + "/"
)

type FSHandler struct {
	dashboard.Handler

	separator string
}

func NewFSHandler() *FSHandler {
	return &FSHandler{
		separator: string(os.PathSeparator),
	}
}

func (h *FSHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	namespace := strings.ToLower(r.URL().Query().Get(":namespace"))
	if namespace == "" {
		h.NotFound(w, r)
		return
	}

	namespaces := dashboard.ConfigFromContext(r.Context()).String(storage.ConfigFileNameSpaces)
	if namespaces == "" {
		h.NotFound(w, r)
		return
	}

	var (
		fileHandler http.Handler
		path        string
		prefix      string
		err         error
	)

	for _, n := range strings.Split(namespaces, ",") {
		parts := strings.Split(n, ":")
		if len(parts) < 2 {
			continue
		}

		ns := strings.ToLower(strings.TrimSpace(parts[0]))
		if ns != namespace {
			continue
		}

		path = strings.TrimRight(strings.TrimSpace(parts[1]), h.separator)
		path = filepath.FromSlash(path)
		path, err = filepath.Abs(path)

		if err != nil {
			continue
		}

		prefix = routePrefix + namespace + "/"
		fileHandler = http.StripPrefix(prefix, http.FileServer(http.Dir(path)))
	}

	if fileHandler == nil {
		h.NotFound(w, r)
		return
	}

	fileName := strings.TrimPrefix(r.URL().Path, prefix)
	fileName = strings.TrimLeft(fileName, h.separator)
	fileName = filepath.FromSlash(fileName)
	fileName = filepath.Join(path, strings.TrimLeft(fileName, "/"))

	fileName, err = filepath.Abs(fileName)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		h.NotFound(w, r)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		h.NotFound(w, r)
		return
	}

	if stat.IsDir() {
		h.NotFound(w, r)
		return
	}

	fileHandler.ServeHTTP(w, r.Original())
}
