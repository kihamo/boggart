package boggart

type Door interface {
	IsOpen() bool
	IsClose() bool
}
