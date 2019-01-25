package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/shadow/components/dashboard"
)

type FSHandler struct {
	dashboard.Handler

	Component storage.Component
}

func (h *FSHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	namespace := strings.ToLower(r.URL().Query().Get(":namespace"))
	if namespace == "" {
		h.NotFound(w, r)
		return
	}

	path, err := h.Component.NamespacePath(namespace)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	prefix := storage.RouteFileStoragePrefix + namespace + "/"
	fileHandler := http.StripPrefix(prefix, http.FileServer(http.Dir(path)))

	fileName := strings.TrimPrefix(r.URL().Path, prefix)
	fileName = strings.TrimLeft(fileName, storage.Separator)
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

	if mime, err := storage.MimeTypeFromData(file); err == nil {
		w.Header().Set("Content-Type", mime.String())
		w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	}

	if r.IsGet() {
		fileHandler.ServeHTTP(w, r.Original())
	}
}
