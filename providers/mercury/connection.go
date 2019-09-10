package mercury

type Connection interface {
	Invoke(request []byte) (response []byte, err error)
}
