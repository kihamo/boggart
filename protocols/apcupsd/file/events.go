package file

type EventsReader struct {
	reader
}

func NewEventsReader(path string) *EventsReader {
	return &EventsReader{
		reader{
			path: path,
		},
	}
}
