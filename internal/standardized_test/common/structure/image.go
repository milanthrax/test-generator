package common

type ImageFormat string

const (
	JPEG ImageFormat = "jpeg"
	PNG  ImageFormat = "png"
	GIF  ImageFormat = "gif"
)

// ImageFile references a photograph or diagram used in a question.
type ImageFile struct {
	Contents []byte      `json:"contents,omitempty"` // base64-encoded when serialized
	Format   ImageFormat `json:"format"`
	FilePath string      `json:"file_path,omitempty"`
	AltText  string      `json:"alt_text,omitempty"` // accessibility description
}
