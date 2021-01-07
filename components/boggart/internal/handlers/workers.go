package handlers

import (
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/shadow/components/dashboard"
)

type WorkersHandler struct {
	dashboard.Handler

	manager *tasks.Manager
}

func NewWorkersHandler(manager *tasks.Manager) *WorkersHandler {
	return &WorkersHandler{
		manager: manager,
	}
}

func (h *WorkersHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	workerID := r.URL().Query().Get("id")
	action := r.URL().Query().Get("action")

	if workerID != "" {
		switch action {
		case "run":
			err := h.manager.Handle(r.Context(), workerID)
			if errors.Is(err, tasks.ErrTaskNotFound) {
				h.NotFound(w, r)
				return
			}

			if err != nil {
				_ = w.SendJSON(boggart.NewResponseJSON().FailedError(err))
				return
			}

			_ = w.SendJSON(boggart.NewResponseJSON().Success(""))
			return

		case "unregister":
			if !r.IsPost() {
				h.MethodNotAllowed(w, r)
				return
			}

			h.manager.Unregister(workerID)

			_ = w.SendJSON(boggart.NewResponseJSON().Success("Worker #" + workerID + " unregister success"))
			return

		case "recalculate":
			if !r.IsPost() {
				h.MethodNotAllowed(w, r)
				return
			}

			err := h.manager.Recalculate(workerID)
			if !errors.Is(err, tasks.ErrTaskNotFound) {
				if err != nil {
					_ = w.SendJSON(boggart.NewResponseJSON().FailedError(err))
					return
				}

				_ = w.SendJSON(boggart.NewResponseJSON().Success("Worker #" + workerID + " recalculate success"))
				return
			}
		}

		h.NotFound(w, r)
		return
	}

	info := h.manager.Info()

	h.Render(r.Context(), "workers", map[string]interface{}{
		"workers": info,
	})
}
