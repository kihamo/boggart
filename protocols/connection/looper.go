package connection

type Looper interface {
	Conn

	Loop() (responses <-chan []byte, errors <-chan error, kill chan<- struct{}, done <-chan struct{})
}

type looper struct {
	Conn
}

func NewLooper(conn Conn) Looper {
	if i, ok := conn.(Looper); ok {
		return i
	}

	return &looper{
		Conn: conn,
	}
}

func (l *looper) Loop() (<-chan []byte, <-chan error, chan<- struct{}, <-chan struct{}) {
	response := make(chan []byte)
	errors := make(chan error)
	kill := make(chan struct{}, 1)
	done := make(chan struct{}, 1)

	go func() {
		buf := make([]byte, bufferSize)

		for {
			select {
			case <-kill:
				close(response)
				close(errors)
				close(done)

				return

			default:
				n, err := l.Conn.Read(buf)

				if n > 0 {
					response <- buf[:n]
				}

				if err != nil {
					errors <- err
				}
			}
		}
	}()

	return response, errors, kill, done
}
