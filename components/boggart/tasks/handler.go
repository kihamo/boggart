package tasks

import (
	"context"
	"time"

	"github.com/kihamo/shadow/components/logging"
)

type Handler interface {
	Handle(context.Context, Meta, Task) error
}

type HandlerFunc func(context.Context, Meta, Task) error

func (f HandlerFunc) Handle(ctx context.Context, meta Meta, task Task) error {
	return f(ctx, meta, task)
}

func HandlerFuncFromShortToLong(handler func(ctx context.Context) error) Handler {
	return HandlerFunc(func(ctx context.Context, _ Meta, _ Task) error {
		return handler(ctx)
	})
}

func HandlerWithTimeout(parent Handler, timeout time.Duration) Handler {
	return HandlerFunc(func(ctx context.Context, meta Meta, task Task) error {
		if parent == nil {
			return ErrParentHandlerIsNil
		}

		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()

		return parent.Handle(ctx, meta, task)
	})
}

func HandlerWithLogger(parent Handler, logger logging.Logger) Handler {
	return HandlerFunc(func(ctx context.Context, meta Meta, task Task) error {
		if parent == nil {
			return ErrParentHandlerIsNil
		}

		if logger == nil {
			return ErrLoggerIsNil
		}

		err := parent.Handle(ctx, meta, task)
		if err != nil {
			logger.Error("handle fail",
				"error", err.Error(),
				"task.id", meta.ID(),
				"task.name", task.Name(),
			)
		}

		return err
	})
}
