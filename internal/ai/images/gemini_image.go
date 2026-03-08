package images

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// GeminiImageGenerator implements ImageGenerator using the Gemini / Imagen API.
type GeminiImageGenerator struct {
	client *genai.Client
	model  string
}

func newGeminiClient(apiKey string) (*genai.Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("gemini image: API key must not be empty")
	}
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:      apiKey,
		Backend:     genai.BackendGeminiAPI,
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1beta"},
	})
	if err != nil {
		return nil, fmt.Errorf("gemini image: create client: %w", err)
	}
	return client, nil
}

// NewGeminiImageGenerator creates a new GeminiImageGenerator with the given model.
func NewGeminiImageGenerator(apiKey, model string) (*GeminiImageGenerator, error) {
	client, err := newGeminiClient(apiKey)
	if err != nil {
		return nil, err
	}
	return &GeminiImageGenerator{client: client, model: model}, nil
}

// ListModels returns the names of all models available for this API key.
func ListModels(apiKey string) ([]string, error) {
	client, err := newGeminiClient(apiKey)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	var names []string
	var pageToken string
	for {
		page, err := client.Models.List(ctx, &genai.ListModelsConfig{PageToken: pageToken})
		if err != nil {
			return nil, fmt.Errorf("gemini image: list models: %w", err)
		}
		for _, model := range page.Items {
			names = append(names, model.Name)
		}
		if page.NextPageToken == "" {
			break
		}
		pageToken = page.NextPageToken
	}
	return names, nil
}

// isImagenModel returns true for Imagen models that use the GenerateImages endpoint.
func isImagenModel(model string) bool {
	return strings.Contains(model, "imagen")
}

// GenerateImage generates an image for the given text prompt.
// It returns the raw image bytes and the MIME type (e.g. "image/png" or "image/jpeg").
// Imagen models use the GenerateImages API; Gemini image models use GenerateContent.
func (g *GeminiImageGenerator) GenerateImage(prompt string) ([]byte, string, error) {
	if isImagenModel(g.model) {
		return g.generateWithImagen(prompt)
	}
	return g.generateWithGemini(prompt)
}

func (g *GeminiImageGenerator) generateWithImagen(prompt string) ([]byte, string, error) {
	ctx := context.Background()
	resp, err := g.client.Models.GenerateImages(
		ctx,
		g.model,
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

func (g *GeminiImageGenerator) generateWithGemini(prompt string) ([]byte, string, error) {
	ctx := context.Background()
	resp, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseModalities: []string{"IMAGE", "TEXT"},
		},
	)
	if err != nil {
		return nil, "", fmt.Errorf("gemini image: generate content: %w", err)
	}
	for _, cand := range resp.Candidates {
		if cand.Content == nil {
			continue
		}
		for _, part := range cand.Content.Parts {
			if part.InlineData != nil && len(part.InlineData.Data) > 0 {
				mimeType := part.InlineData.MIMEType
				if mimeType == "" {
					mimeType = "image/png"
				}
				return part.InlineData.Data, mimeType, nil
			}
		}
	}
	return nil, "", fmt.Errorf("gemini image: no image data in response")
}
