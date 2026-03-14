package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	audio "github.com/milanthrax/test-generator/internal/ai/text_to_speech"
	"github.com/milanthrax/test-generator/internal/filesystem"
	common "github.com/milanthrax/test-generator/internal/standardized_test/common/structure"
)

const (
	outputDir = "test_data"

	// Gemini TTS returns 24 kHz, 16-bit, mono linear PCM by default.
	pcmSampleRate  = 24000
	pcmNumChannels = 1
	pcmBitDepth    = 16
)

// scriptLine is the JSON input format: a speaker label and the line of text.
type scriptLine struct {
	Speaker string `json:"speaker"`
	Text    string `json:"text"`
}

// speakerRegistry maps VoiceProfile labels to their profiles.
var speakerRegistry = map[string]common.VoiceProfile{
	"Formal Announcer":   common.Announcer,
	"Smooth Narrator":    common.NarratorFemale,
	"News Anchor":        common.NewsAnchor,
	"Toddler Boy":        common.ToddlerBoy,
	"Toddler Girl":       common.ToddlerGirl,
	"Small Boy":          common.YoungChildBoy,
	"Small Girl":         common.YoungChildGirl,
	"Pre-teen Boy":       common.PreTeenBoy,
	"Pre-teen Girl":      common.PreTeenGirl,
	"College Male":       common.CollegeMale,
	"College Female":     common.CollegeFemale,
	"Excited Teenager":   common.HyperTeen,
	"Casual Young Man":   common.BroMale,
	"Professional Man":   common.BusinessMale,
	"Professional Woman": common.BusinessFemale,
	"Deep Bass Man":      common.DeepVoiceMan,
	"Motherly Voice":     common.WarmMotherly,
	"Strict Teacher":     common.StrictTeacher,
	"Middle-aged Man":    common.MatureMale,
	"Middle-aged Woman":  common.MatureFemale,
	"Gruff Man":          common.GruffContractor,
	"Grandfather":        common.Grandpa,
	"Grandmother":        common.Grandma,
	"Elderly Raspy":      common.ElderlyRaspy,
	"Whispering Voice":   common.Whisperer,
	"Fast Talker":        common.FastTalker,
	"Very Slow Speaker":  common.SlowMover,
	"Robotic Tone":       common.MonotoneRobot,
	"Panicked Voice":     common.HighPitchPanic,
	"Grumbling Voice":    common.LowGrumble,
}

func main() {
	outFlag := flag.String("out", "", "output audio file path (e.g. test_data/my-conv.m4a or .wav); defaults to test_data/conversation_audio.m4a")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: text_to_speech [flags] <conversation.json>\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	apiKey := os.Getenv("TTS_GENERATION_KEY")
	if apiKey == "" {
		log.Fatal("TTS_GENERATION_KEY environment variable is not set")
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	conversationFile := flag.Arg(0)

	// Derive the desired output format from the -out extension (default: m4a).
	desiredFormat := "m4a"
	if *outFlag != "" {
		if ext := strings.TrimPrefix(filepath.Ext(*outFlag), "."); ext != "" {
			desiredFormat = strings.ToLower(ext)
		}
	}

	// Read and parse the script JSON ([{speaker, text}, ...]).
	data, err := os.ReadFile(conversationFile)
	if err != nil {
		log.Fatalf("failed to read conversation file: %v", err)
	}
	var lines []scriptLine
	if err := json.Unmarshal(data, &lines); err != nil {
		log.Fatalf("failed to parse conversation JSON: %v", err)
	}

	conv, err := buildConversation(lines)
	if err != nil {
		log.Fatalf("failed to build conversation: %v", err)
	}

	audioFile, err := generateConversationAudio(conv, apiKey, desiredFormat)
	if err != nil {
		log.Fatalf("failed to generate conversation audio: %v", err)
	}

	// Ensure the output directory exists.
	if err := filesystem.CreateDirIfNotExists(outputDir); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// Determine the final output path.
	outputPath := *outFlag
	if outputPath == "" {
		outputPath = filepath.Join(outputDir, "conversation_audio."+string(audioFile.Format))
	} else {
		// If conversion fell back (e.g. M4A failed → WAV), correct the extension.
		wantExt := "." + string(audioFile.Format)
		if filepath.Ext(outputPath) != wantExt {
			log.Printf("note: actual format is %s; adjusting extension", audioFile.Format)
			outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + wantExt
		}
	}

	if err := os.WriteFile(outputPath, audioFile.Contents, 0644); err != nil {
		log.Fatalf("failed to write audio file: %v", err)
	}
	audioFile.FilePath = outputPath

	// Write the AudioFile metadata (without raw bytes) as JSON alongside the audio.
	meta := audioFile
	meta.Contents = nil
	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal audio metadata: %v", err)
	}
	metaPath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".json"
	if err := os.WriteFile(metaPath, metaBytes, 0644); err != nil {
		log.Fatalf("failed to write metadata file: %v", err)
	}

	log.Printf("audio saved to %s (%s, %d bytes)", outputPath, audioFile.Format, len(audioFile.Contents))
	log.Printf("metadata saved to %s", metaPath)
}

// buildConversation maps a flat list of script lines to a Conversation using the voice registry.
func buildConversation(lines []scriptLine) (common.Conversation, error) {
	convLines := make([]common.Line, 0, len(lines))
	for _, l := range lines {
		profile, ok := speakerRegistry[l.Speaker]
		if !ok {
			return common.Conversation{}, fmt.Errorf("unknown speaker label %q", l.Speaker)
		}
		convLines = append(convLines, common.Line{Speaker: profile, Text: l.Text})
	}
	return common.Conversation{Lines: convLines}, nil
}

// generateConversationAudio synthesizes speech for the given Conversation via SSML and
// returns an AudioFile in the requested format ("wav", "m4a", "mp3").
func generateConversationAudio(conv common.Conversation, apiKey, targetFormat string) (common.AudioFile, error) {
	generator, err := audio.NewGeminiTTSGenerator(apiKey, "")
	if err != nil {
		return common.AudioFile{}, err
	}

	ssml := conv.ToSSML()
	log.Printf("Sending SSML to Gemini TTS:\n%s", ssml)

	audioBytes, mimeType, err := generator.GenerateTextToSpeech(ssml)
	if err != nil {
		return common.AudioFile{}, err
	}

	log.Printf("Gemini TTS returned MIME type: %q, %d bytes", mimeType, len(audioBytes))

	// Normalise the MIME type to lowercase and strip any parameters
	// (e.g. "audio/L16;rate=24000" → "audio/l16").
	mimeBase := strings.ToLower(strings.SplitN(mimeType, ";", 2)[0])
	mimeBase = strings.TrimSpace(mimeBase)

	// Treat any PCM / L16 / raw variant — or an unrecognised type — as raw
	// linear PCM and wrap it in a WAV container.
	isPCM := strings.Contains(mimeBase, "pcm") ||
		mimeBase == "audio/l16" ||
		strings.Contains(mimeBase, "raw") ||
		(mimeBase != "audio/mp3" &&
			mimeBase != "audio/mpeg" &&
			mimeBase != "audio/ogg" &&
			mimeBase != "audio/wav" &&
			mimeBase != "audio/wave" &&
			mimeBase != "audio/x-wav")
	if isPCM {
		log.Printf("Wrapping raw PCM bytes in WAV container (sample rate %d Hz, %d-bit, %d ch)",
			pcmSampleRate, pcmBitDepth, pcmNumChannels)
		audioBytes = buildWAV(audioBytes, pcmSampleRate, pcmNumChannels, pcmBitDepth)
		mimeBase = "audio/wav"
	}

	// Convert WAV → M4A only when the caller requested it (default).
	if (mimeBase == "audio/wav" || mimeBase == "audio/wave" || mimeBase == "audio/x-wav") && targetFormat == "m4a" {
		m4a, err := wavToM4A(audioBytes)
		if err != nil {
			log.Printf("warning: M4A conversion failed (%v); saving as WAV", err)
			mimeBase = "audio/wav"
		} else {
			audioBytes = m4a
			mimeBase = "audio/mp4"
		}
	}

	ext := extensionFromMIME(mimeBase)
	format := common.AudioFormat(strings.TrimPrefix(ext, "."))

	return common.AudioFile{
		Contents:   audioBytes,
		Format:     format,
		Transcript: ssml,
	}, nil
}

// wavToM4A encodes WAV bytes to AAC inside an MPEG-4 (M4A) container via ffmpeg.
// It prefers the macOS AudioToolbox encoder (aac_at) which uses the OS-native
// AAC implementation — no extra codec library installation required.
// Falls back to ffmpeg's built-in AAC encoder on non-macOS systems.
func wavToM4A(wavBytes []byte) ([]byte, error) {
	ffmpegPath, err := resolveBinary("ffmpeg",
		"/opt/homebrew/bin/ffmpeg", // Apple Silicon Homebrew
		"/usr/local/bin/ffmpeg",    // Intel Mac / Linux Homebrew
		"/usr/bin/ffmpeg",
	)
	if err != nil {
		return nil, err
	}

	// Try AudioToolbox first (macOS native, zero extra installs), then the
	// built-in ffmpeg AAC encoder as a cross-platform fallback.
	for _, codec := range []string{"aac_at", "aac"} {
		var out, errBuf bytes.Buffer
		cmd := exec.Command(
			ffmpegPath,
			"-hide_banner", "-loglevel", "error",
			"-i", "pipe:0",
			"-c:a", codec,
			"-b:a", "128k",
			"-f", "mp4",
			"-movflags", "frag_keyframe+empty_moov", // required for pipe output
			"pipe:1",
		)
		cmd.Stdin = bytes.NewReader(wavBytes)
		cmd.Stdout = &out
		cmd.Stderr = &errBuf

		if err := cmd.Run(); err != nil {
			log.Printf("codec %s failed (%s), trying next", codec, strings.TrimSpace(errBuf.String()))
			continue
		}
		log.Printf("encoded M4A using codec %s", codec)
		return out.Bytes(), nil
	}
	return nil, fmt.Errorf("all AAC codecs failed; make sure ffmpeg is installed (brew install ffmpeg)")
}

// resolveBinary looks up a binary by name on PATH, then tries each fallback path.
func resolveBinary(name string, fallbacks ...string) (string, error) {
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}
	for _, p := range fallbacks {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("%q not found on PATH or in %v", name, fallbacks)
}

// buildWAV wraps raw linear PCM bytes in a RIFF/WAV container.
func buildWAV(pcm []byte, sampleRate, numChannels, bitDepth int) []byte {
	byteRate := sampleRate * numChannels * bitDepth / 8
	blockAlign := numChannels * bitDepth / 8
	dataSize := len(pcm)
	riffSize := 36 + dataSize

	buf := make([]byte, 44+dataSize)
	copy(buf[0:4], []byte("RIFF"))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(riffSize))
	copy(buf[8:12], []byte("WAVE"))
	copy(buf[12:16], []byte("fmt "))
	binary.LittleEndian.PutUint32(buf[16:20], 16) // PCM chunk size
	binary.LittleEndian.PutUint16(buf[20:22], 1)  // PCM format
	binary.LittleEndian.PutUint16(buf[22:24], uint16(numChannels))
	binary.LittleEndian.PutUint32(buf[24:28], uint32(sampleRate))
	binary.LittleEndian.PutUint32(buf[28:32], uint32(byteRate))
	binary.LittleEndian.PutUint16(buf[32:34], uint16(blockAlign))
	binary.LittleEndian.PutUint16(buf[34:36], uint16(bitDepth))
	copy(buf[36:40], []byte("data"))
	binary.LittleEndian.PutUint32(buf[40:44], uint32(dataSize))
	copy(buf[44:], pcm)
	return buf
}

// extensionFromMIME returns a file extension (with leading dot) for the given audio MIME type.
func extensionFromMIME(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "audio/pcm"),
		mimeType == "audio/l16",
		mimeType == "audio/wav",
		mimeType == "audio/wave",
		mimeType == "audio/x-wav":
		return ".wav"
	case mimeType == "audio/mp3", mimeType == "audio/mpeg":
		return ".mp3"
	case mimeType == "audio/mp4", mimeType == "audio/m4a", mimeType == "video/mp4":
		return ".m4a"
	}
	// Fallback: try stdlib mime.
	exts, err := mime.ExtensionsByType(mimeType)
	if err == nil && len(exts) > 0 {
		return exts[0]
	}
	return ".wav" // safe fallback — we always produce at least a WAV
}
