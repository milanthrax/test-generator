package audio

type TextToSpeechGenerator interface {
	GenerateTextToSpeech(prompt string) ([]byte, string, error) // returns audio bytes, format, and error
}
