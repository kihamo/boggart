package voice

type Speaker interface {
	Speech(text string) error
	SpeechWithOptions(text string, volume int64, speed float64) error
}
