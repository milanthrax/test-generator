package common

type AudioFormat string

const (
	MP3 AudioFormat = "mp3"
	WAV AudioFormat = "wav"
	OGG AudioFormat = "ogg"
)

// AudioFile references a spoken audio clip used in listening questions.
type AudioFile struct {
	Contents     []byte      `json:"contents,omitempty"` // base64-encoded when serialized
	Format       AudioFormat `json:"format"`
	FilePath     string      `json:"file_path,omitempty"`
	DurationSecs int         `json:"duration_secs"`
	Transcript   string      `json:"transcript,omitempty"` // optional verbatim transcript of the audio
}
