package images

type ImageGenerator interface {
	GenerateImage(prompt string) ([]byte, string, error) // returns image bytes, format, and error
}
