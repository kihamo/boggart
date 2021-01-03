package tasks

const (
	StatusNew Status = iota
	StatusRunning
	StatusWaiting
	StatusClosed
)

type Status uint32

func (i Status) Is(status Status) bool {
	return i == status
}

func (i Status) IsNew() bool {
	return i.Is(StatusNew)
}

func (i Status) IsRunning() bool {
	return i.Is(StatusRunning)
}

func (i Status) IsWaiting() bool {
	return i.Is(StatusWaiting)
}

func (i Status) IsClosed() bool {
	return i.Is(StatusClosed)
}

func (i Status) IsAllowExecute() bool {
	return i.IsNew() || i.IsWaiting()
}

func (i Status) Uint32() uint32 {
	return uint32(i)
}
