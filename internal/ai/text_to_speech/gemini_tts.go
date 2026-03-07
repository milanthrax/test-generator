package audio

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

const ttsModel = "gemini-2.5-flash-preview-tts"

// GeminiTTSGenerator implements TextToSpeechGenerator using the Gemini API.
type GeminiTTSGenerator struct {
	client *genai.Client
	voice  string // Gemini prebuilt voice name, e.g. "Kore", "Charon", "Fenrir"
}

// NewGeminiTTSGenerator creates a new GeminiTTSGenerator.
// apiKey is the Gemini API key.
// voice is an optional prebuilt voice name (defaults to "Kore" when empty).
func NewGeminiTTSGenerator(apiKey, voice string) (*GeminiTTSGenerator, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("gemini TTS: API key must not be empty")
	}
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("gemini TTS: create client: %w", err)
	}
	if voice == "" {
		voice = "Kore"
	}
	return &GeminiTTSGenerator{client: client, voice: voice}, nil
}

// GenerateTextToSpeech synthesizes speech for the given text prompt.
// It returns the raw audio bytes and the MIME type (e.g. "audio/pcm" or "audio/mp3").
func (g *GeminiTTSGenerator) GenerateTextToSpeech(prompt string) ([]byte, string, error) {
	ctx := context.Background()

	resp, err := g.client.Models.GenerateContent(
		ctx,
		ttsModel,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			ResponseModalities: []string{"AUDIO"},
			SpeechConfig: &genai.SpeechConfig{
				VoiceConfig: &genai.VoiceConfig{
					PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
						VoiceName: g.voice,
					},
				},
			},
		},
	)
	if err != nil {
		return nil, "", fmt.Errorf("gemini TTS: generate content: %w", err)
	}

	for _, cand := range resp.Candidates {
		if cand.Content == nil {
			continue
		}
		for _, part := range cand.Content.Parts {
			if part.InlineData != nil && strings.HasPrefix(part.InlineData.MIMEType, "audio/") {
				return part.InlineData.Data, part.InlineData.MIMEType, nil
			}
		}
	}

	return nil, "", fmt.Errorf("gemini TTS: no audio data in response")
}
