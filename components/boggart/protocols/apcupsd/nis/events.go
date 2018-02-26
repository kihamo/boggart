package nis

type EventsReader struct {
	reader
}

func NewEventsReader(address string) *EventsReader {
	return &EventsReader{
		reader{
			address: address,
			command: "events",
		},
	}
}
