# Test Generator

Contains different generators for standardized tests using Google Gemini AI.

---

## Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| [Go](https://go.dev/dl/) | 1.21+ | `brew install go` |
| [ffmpeg](https://ffmpeg.org) | any recent | `brew install ffmpeg` — only needed for M4A encoding; uses macOS AudioToolbox (no extra codec installs) |
| Gemini API key | — | [Get one here](https://aistudio.google.com/app/apikey) |

> **Why ffmpeg?** Gemini TTS returns raw PCM audio. The app wraps it in a WAV container
> entirely in Go (no external tools), then hands it to ffmpeg for AAC/M4A encoding using
> the **macOS AudioToolbox** codec (`aac_at`) — a codec that ships with every Mac.
> No additional codec libraries need to be downloaded or compiled.

---

## Environment setup

```bash
# 1. Install Go dependencies
go mod download

# 2. Install ffmpeg (macOS)
brew install ffmpeg

# 3. Export your Gemini API key
export TTS_GENERATION_KEY="your-api-key-here"

# Optional: add to your shell profile so you don't have to re-export every session
echo 'export TTS_GENERATION_KEY="your-api-key-here"' >> ~/.zshrc
```

---

## Text-to-Speech generator

Converts a JSON conversation file into an M4A audio file using Gemini TTS.

### Input format

A JSON array of `{speaker, text}` objects. The `speaker` value must be one of the
predefined voice labels (see [Voice library](#voice-library) below).

```json
[
  { "speaker": "Strict Teacher", "text": "Good morning, class." },
  { "speaker": "Small Boy",      "text": "Good morning, Miss!" }
]
```

### Run

```bash
go run ./cmd/common/text_to_speech/... -out <output.m4a> <conversation.json>
```

### Sample conversations

```bash
go run ./cmd/common/text_to_speech/... -out test_data/conv-001-new-student.m4a    data/tts/sample-conversation-001.json
go run ./cmd/common/text_to_speech/... -out test_data/conv-002-homework.m4a        data/tts/sample-conversation-002.json
go run ./cmd/common/text_to_speech/... -out test_data/conv-003-sports.m4a          data/tts/sample-conversation-003.json
go run ./cmd/common/text_to_speech/... -out test_data/conv-004-movie-night.m4a     data/tts/sample-conversation-004.json
go run ./cmd/common/text_to_speech/... -out test_data/conv-005-grandparents.m4a    data/tts/sample-conversation-005.json
go run ./cmd/common/text_to_speech/... -out test_data/conv-006-standup.m4a         data/tts/sample-conversation-006.json
go run ./cmd/common/text_to_speech/... -out test_data/conv-007-dinner.m4a          data/tts/sample-conversation-007.json
go run ./cmd/common/text_to_speech/... -out test_data/conv-008-news.m4a            data/tts/sample-conversation-008.json
```

### Voice library

| Label | Character |
|-------|-----------|
| `Formal Announcer` | Formal male announcer |
| `Smooth Narrator` | Calm female narrator |
| `News Anchor` | Male news anchor |
| `Toddler Boy` | Very young boy |
| `Toddler Girl` | Very young girl |
| `Small Boy` | Young child (boy) |
| `Small Girl` | Young child (girl) |
| `Pre-teen Boy` | Pre-adolescent boy |
| `Pre-teen Girl` | Pre-adolescent girl |
| `College Male` | Young adult male |
| `College Female` | Young adult female |
| `Excited Teenager` | Energetic teenage girl |
| `Casual Young Man` | Relaxed young man |
| `Professional Man` | Adult business male |
| `Professional Woman` | Adult business female |
| `Deep Bass Man` | Deep-voiced adult male |
| `Motherly Voice` | Warm maternal female |
| `Strict Teacher` | Authoritative female teacher |
| `Middle-aged Man` | Middle-aged male |
| `Middle-aged Woman` | Middle-aged female |
| `Gruff Man` | Gruff middle-aged male |
| `Grandfather` | Elderly male |
| `Grandmother` | Elderly female |
| `Elderly Raspy` | Raspy elderly male |
| `Whispering Voice` | Soft whispering female |
| `Fast Talker` | Rapid-speech male |
| `Very Slow Speaker` | Slow-speech male |
| `Robotic Tone` | Monotone robotic female |
| `Panicked Voice` | High-pitched panicked female |
| `Grumbling Voice` | Low grumbling male |
