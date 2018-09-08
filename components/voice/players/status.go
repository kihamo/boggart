package players

type Status int64

const (
	StatusUnknown Status = iota
	StatusStopped
	StatusPlaying
	StatusPause
)

func (i Status) Int64() int64 {
	if i < 0 || i >= Status(len(_StatusIndex)-1) {
		return -1
	}

	return int64(i)
}
