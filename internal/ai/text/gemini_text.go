package text

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

const textModel = "gemini-2.5-flash"

// GeminiTextGenerator implements TextGenerator using the Gemini API.
type GeminiTextGenerator struct {
	client *genai.Client
}

// NewGeminiTextGenerator creates a new GeminiTextGenerator.
// apiKey is the Gemini API key.
func NewGeminiTextGenerator(apiKey string) (*GeminiTextGenerator, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("gemini text: API key must not be empty")
	}
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("gemini text: create client: %w", err)
	}
	return &GeminiTextGenerator{client: client}, nil
}

// GenerateText sends the prompt to Gemini and returns the generated text.
func (g *GeminiTextGenerator) GenerateText(prompt string) (string, error) {
	ctx := context.Background()

	resp, err := g.client.Models.GenerateContent(
		ctx,
		textModel,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("gemini text: generate content: %w", err)
	}

	var sb strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content == nil {
			continue
		}
		for _, part := range cand.Content.Parts {
			if part.Text != "" {
				sb.WriteString(part.Text)
			}
		}
	}

	result := strings.TrimSpace(sb.String())
	if result == "" {
		return "", fmt.Errorf("gemini text: empty response")
	}
	return result, nil
}
