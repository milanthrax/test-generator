package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/draw"

	"github.com/milanthrax/test-generator/internal/ai/images"
	"github.com/milanthrax/test-generator/internal/filesystem"
)

const (
	outputPath = "test_data/image_generation_img.jpg"
	imgWidth   = 640
	imgHeight  = 480
)

func main() {
	apiKey := os.Getenv("IMAGE_GENERATION_API_KEY")
	if apiKey == "" {
		log.Fatal("IMAGE_GENERATION_API_KEY environment variable is not set")
	}

	if len(os.Args) < 2 {
		log.Fatal("usage: image_generation <prompt>")
	}
	prompt := strings.Join(os.Args[1:], " ")

	// List available image generation models and let the user pick one.
	fmt.Println("Fetching available image generation models...")
	models, err := images.ListModels(apiKey)
	if err != nil {
		log.Fatalf("failed to list models: %v", err)
	}
	if len(models) == 0 {
		log.Fatal("no image generation models available for this API key")
	}

	fmt.Println("Available models:")
	for i, m := range models {
		fmt.Printf("  [%d] %s\n", i+1, m)
	}

	reader := bufio.NewReader(os.Stdin)
	var selectedModel string
	for {
		fmt.Printf("Select a model (1-%d): ", len(models))
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		idx, err := strconv.Atoi(line)
		if err == nil && idx >= 1 && idx <= len(models) {
			selectedModel = models[idx-1]
			break
		}
		fmt.Printf("Invalid selection, please enter a number between 1 and %d.\n", len(models))
	}
	fmt.Printf("Using model: %s\n", selectedModel)

	generator, err := images.NewGeminiImageGenerator(apiKey, selectedModel)
	if err != nil {
		log.Fatalf("failed to create image generator: %v", err)
	}

	imgBytes, _, err := generator.GenerateImage(prompt)
	if err != nil {
		log.Fatalf("failed to generate image: %v", err)
	}

	// Decode the generated image (PNG or JPEG).
	src, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	// Resize to 640x480 using a high-quality scaler.
	dst := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	// Ensure the output directory exists.
	if err := filesystem.CreateDirIfNotExists("test_data"); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// Write the JPEG file.
	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, dst, &jpeg.Options{Quality: 90}); err != nil {
		log.Fatalf("failed to encode JPEG: %v", err)
	}

	log.Printf("image saved to %s (%dx%d)", outputPath, imgWidth, imgHeight)
}
