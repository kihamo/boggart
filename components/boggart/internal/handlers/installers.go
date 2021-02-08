package handlers

import (
	"bytes"
	"io"
	"path/filepath"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/shadow/components/dashboard"
)

type InstallerHandler struct {
	dashboard.Handler

	component boggart.Component
}

func NewInstallerHandler(b boggart.Component) *InstallerHandler {
	return &InstallerHandler{
		component: b,
	}
}

func (h *InstallerHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	var bindItem boggart.BindItem

	id := q.Get(":id")
	systemID := q.Get(":system")

	if id == "" || systemID == "" {
		h.NotFound(w, r)
		return
	}

	bindItem = h.component.Bind(id)
	if bindItem == nil {
		h.NotFound(w, r)
		return
	}

	install, ok := bindItem.Bind().(installer.HasInstaller)
	if !ok {
		h.NotFound(w, r)
		return
	}

	var system installer.System

	for _, s := range install.InstallersSupport() {
		if s == installer.System(systemID) {
			system = s
		}
	}

	if system == "" {
		h.NotFound(w, r)
		return
	}

	steps, err := install.InstallerSteps(r.Context(), system)
	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	}

	if len(steps) == 0 && err == nil {
		h.NotFound(w, r)
		return
	}

	// merge steps by file
	stepsMerged := make(map[string]int, len(steps))

	for i := len(steps) - 1; i >= 0; i-- {
		if steps[i].Content == "" {
			steps = append(steps[:i], steps[i+1:]...)
			continue
		}

		if steps[i].FilePath == "" {
			continue
		}

		if existStepIndex, ok := stepsMerged[steps[i].FilePath]; !ok {
			// оставляем только последний вариант
			stepsMerged[steps[i].FilePath] = len(steps) - i
		} else {
			// все остальные склеиваем
			existStep := &steps[len(steps)-existStepIndex]

			if existStep.Description != "" {
				existStep.Description = "\n" + existStep.Description
			}
			existStep.Description = steps[i].Description + existStep.Description

			if existStep.Content != "" {
				existStep.Content = "\n" + existStep.Content
			}
			existStep.Content = steps[i].Content + existStep.Content

			steps = append(steps[:i], steps[i+1:]...)
		}
	}

	if filePath := q.Get("file"); filePath != "" {
		index := -1
		if s := q.Get("step"); s != "" {
			if i, err := strconv.Atoi(s); err == nil {
				index = i
			}
		}

		buf := bytes.NewBuffer(nil)

		for i, step := range steps {
			if index > -1 && i != index || step.FilePath != filePath {
				continue
			}

			if buf.Len() > 0 {
				buf.WriteString("\n\n")
			}

			buf.WriteString(step.Content)

			if index > -1 && i == index {
				break
			}
		}

		w.Header().Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filepath.Base(filePath)+"\"")

		_, _ = io.Copy(w, buf)
		return
	}

	h.Render(r.Context(), "installer", map[string]interface{}{
		"system": systemID,
		"steps":  steps,
	})
}
