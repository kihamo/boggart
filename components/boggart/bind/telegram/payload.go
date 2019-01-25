package telegram

type FilePayload struct {
	URL  string `json:"url"`
	MIME string `json:"mime,omitempty"`
	Name string `json:"name,omitempty"`
}
