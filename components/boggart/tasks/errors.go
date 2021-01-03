package tasks

import (
	"errors"
)

var (
	ErrTaskIsEmpty        = errors.New("task is empty")
	ErrHandlerIsEmpty     = errors.New("handler is empty")
	ErrScheduleIsEmpty    = errors.New("schedule is empty")
	ErrAlreadyRunning     = errors.New("already running")
	ErrTaskNotFound       = errors.New("task not found")
	ErrParentHandlerIsNil = errors.New("parent handler is nil")
	ErrLoggerIsNil        = errors.New("logger is nil")
)
