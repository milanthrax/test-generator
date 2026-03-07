package images

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

const imageModel = "imagen-3.0-generate-002"

// GeminiImageGenerator implements ImageGenerator using the Gemini / Imagen API.
type GeminiImageGenerator struct {
	client *genai.Client
}

// NewGeminiImageGenerator creates a new GeminiImageGenerator.
// apiKey is the Gemini API key.
func NewGeminiImageGenerator(apiKey string) (*GeminiImageGenerator, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("gemini image: API key must not be empty")
	}
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("gemini image: create client: %w", err)
	}
	return &GeminiImageGenerator{client: client}, nil
}

// GenerateImage generates an image for the given text prompt.
// It returns the raw image bytes and the MIME type (e.g. "image/png" or "image/jpeg").
func (g *GeminiImageGenerator) GenerateImage(prompt string) ([]byte, string, error) {
	ctx := context.Background()

	resp, err := g.client.Models.GenerateImages(
		ctx,
		imageModel,
		prompt,
		&genai.GenerateImagesConfig{
			NumberOfImages: 1,
		},
	)
	if err != nil {
		return nil, "", fmt.Errorf("gemini image: generate images: %w", err)
	}

	if len(resp.GeneratedImages) == 0 || resp.GeneratedImages[0].Image == nil {
		return nil, "", fmt.Errorf("gemini image: no image data in response")
	}

	img := resp.GeneratedImages[0].Image
	mimeType := img.MIMEType
	if mimeType == "" {
		mimeType = "image/png"
	}
	return img.ImageBytes, mimeType, nil
}
