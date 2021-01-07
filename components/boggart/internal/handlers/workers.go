package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/shadow/components/dashboard"
)

// easyjson:json
type workersHandlerItem struct {
	ID              string     `json:"id"`
	Name            string     `json:"name,omitempty"`
	Status          string     `json:"status"`
	AttemptsSuccess uint64     `json:"attempts_success"`
	AttemptsFails   uint64     `json:"attempts_fails"`
	FirstRunAt      *time.Time `json:"first_run_at,omitempty"`
	LastRunAt       *time.Time `json:"last_run_at,omitempty"`
	NextRunAt       *time.Time `json:"next_run_at,omitempty"`
}

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
	switch r.URL().Query().Get("action") {
	case "run":
		workerID := r.URL().Query().Get("id")
		if workerID == "" {
			break
		}

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

		workerID := r.URL().Query().Get("id")
		if workerID == "" {
			break
		}

		h.manager.Unregister(workerID)

		_ = w.SendJSON(boggart.NewResponseJSON().Success("Worker #" + workerID + " unregister success"))
		return

	case "recalculate":
		if !r.IsPost() {
			h.MethodNotAllowed(w, r)
			return
		}

		workerID := r.URL().Query().Get("id")
		if workerID == "" {
			break
		}

		var nextRunBefore, nextRunAfter *time.Time

		meta, err := h.manager.Meta(workerID)
		if err == nil {
			nextRunBefore = meta.NextRunAt()

			err = h.manager.Recalculate(workerID)
		}

		if err == nil {
			meta, err = h.manager.Meta(workerID)
		}

		if err == nil {
			var durationBetween, durationNow time.Duration

			nextRunAfter = meta.NextRunAt()

			if nextRunAfter != nil {
				durationNow = time.Until(*nextRunAfter)

				if nextRunBefore != nil {
					durationBetween = nextRunAfter.Sub(*nextRunBefore)
				}
			}

			_ = w.SendJSON(boggart.NewResponseJSON().Success(
				fmt.Sprintf(
					"Worker #"+workerID+" recalculate success from %v to %v (from previous %v, from now %v)",
					nextRunBefore, nextRunAfter, durationBetween, durationNow,
				),
			))
			return
		}

		if !errors.Is(err, tasks.ErrTaskNotFound) {
			_ = w.SendJSON(boggart.NewResponseJSON().FailedError(err))
			return
		}

	case "list":
		if r.IsAjax() {
			items := h.manager.Info()
			data := make([]workersHandlerItem, len(items))

			for i, item := range items {
				data[i].ID = item.Meta.ID()
				data[i].Name = item.Task.Name()
				data[i].Status = item.Meta.Status().String()
				data[i].AttemptsSuccess = item.Meta.Success()
				data[i].AttemptsFails = item.Meta.Fails()
				data[i].FirstRunAt = item.Meta.FirstRunAt()
				data[i].LastRunAt = item.Meta.LastRunAt()
				data[i].NextRunAt = item.Meta.NextRunAt()
			}

			reply := boggart.NewResponseDataTable()
			reply.Data = data
			reply.Total = len(data)
			reply.Filtered = reply.Total

			w.SendJSON(reply)
			return
		}

	default:
		h.Render(r.Context(), "workers", nil)
		return
	}

	h.NotFound(w, r)
}
