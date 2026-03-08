package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/milanthrax/test-generator/internal/ai/text"
)

func main() {
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		log.Fatal("LLM_API_KEY environment variable is not set")
	}

	if len(os.Args) < 2 {
		log.Fatal("usage: text_generation <prompt>")
	}

	prompt := strings.Join(os.Args[1:], " ")

	generator, err := text.NewGeminiTextGenerator(apiKey)
	if err != nil {
		log.Fatalf("failed to create text generator: %v", err)
	}

	result, err := generator.GenerateText(prompt)
	if err != nil {
		log.Fatalf("failed to generate text: %v", err)
	}

	fmt.Println(result)
}
