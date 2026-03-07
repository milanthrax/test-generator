package html

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"strings"

	common "github.com/milanthrax/test-generator/internal/standardized_test/common/structure"
	"github.com/milanthrax/test-generator/internal/standardized_test/toeic/structure"
)

// Render writes a self-contained HTML file for the given test to w.
// lang selects which localized strings to display; falls back to English when a
// translation is missing.
func Render(t *structure.Test, lang structure.LanguageCode, w io.Writer) error {
	tmpl, err := template.New("test").Funcs(template.FuncMap{
		"loc":       loc(lang),
		"audioSrc":  audioSrc,
		"imageSrc":  imageSrc,
		"mimeAudio": mimeAudio,
		"mimeImage": mimeImage,
		"partName":  partName,
		"add":       func(a, b int) int { return a + b },
	}).Parse(testTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	if err := tmpl.Execute(w, t); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	return nil
}

// loc returns a function that picks the best available translation, falling
// back to English then any available value.
func loc(lang structure.LanguageCode) func(structure.LocalizedText) string {
	return func(lt structure.LocalizedText) string {
		if lt == nil {
			return ""
		}
		if v, ok := lt[lang]; ok && v != "" {
			return v
		}
		if v, ok := lt[structure.EN]; ok && v != "" {
			return v
		}
		for _, v := range lt {
			if v != "" {
				return v
			}
		}
		return ""
	}
}

func audioSrc(a *common.AudioFile) string {
	if a == nil {
		return ""
	}
	if len(a.Contents) > 0 {
		return "data:" + mimeAudio(a.Format) + ";base64," +
			base64.StdEncoding.EncodeToString(a.Contents)
	}
	return a.FilePath
}

func imageSrc(img *common.ImageFile) string {
	if img == nil {
		return ""
	}
	if len(img.Contents) > 0 {
		return "data:" + mimeImage(img.Format) + ";base64," +
			base64.StdEncoding.EncodeToString(img.Contents)
	}
	return img.FilePath
}

func mimeAudio(f common.AudioFormat) string {
	switch strings.ToLower(string(f)) {
	case "mp3":
		return "audio/mpeg"
	case "ogg":
		return "audio/ogg"
	case "wav":
		return "audio/wav"
	default:
		return "audio/mpeg"
	}
}

func mimeImage(f common.ImageFormat) string {
	switch strings.ToLower(string(f)) {
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	default:
		return "image/jpeg"
	}
}

func partName(n structure.PartNumber) string {
	switch n {
	case structure.Part1Photographs:
		return "Part 1 – Photographs"
	case structure.Part2QuestionResponse:
		return "Part 2 – Question-Response"
	case structure.Part3Conversations:
		return "Part 3 – Conversations"
	case structure.Part4Talks:
		return "Part 4 – Short Talks"
	case structure.Part5IncompleteSentences:
		return "Part 5 – Incomplete Sentences"
	case structure.Part6TextCompletion:
		return "Part 6 – Text Completion"
	case structure.Part7ReadingComprehension:
		return "Part 7 – Reading Comprehension"
	default:
		return fmt.Sprintf("Part %d", n)
	}
}

// testTemplate is the single-file HTML template for a TOEIC practice test.
const testTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    font-family: Georgia, "Times New Roman", serif;
    background: #f5f5f0;
    color: #1a1a1a;
    line-height: 1.7;
    padding: 2rem 1rem;
  }
  .test-header {
    max-width: 820px;
    margin: 0 auto 2.5rem;
    border-bottom: 3px solid #003580;
    padding-bottom: 1rem;
  }
  .test-header h1 { font-size: 2rem; color: #003580; }
  .test-header p  { color: #555; margin-top: .4rem; }
  .section-block {
    max-width: 820px;
    margin: 0 auto 3rem;
  }
  .section-title {
    font-size: 1.5rem;
    font-weight: bold;
    background: #003580;
    color: #fff;
    padding: .5rem 1rem;
    border-radius: 4px 4px 0 0;
  }
  .part-block {
    background: #fff;
    border: 1px solid #d0d0d0;
    border-top: none;
    padding: 1.5rem;
    margin-bottom: 1.5rem;
  }
  .part-title {
    font-size: 1.15rem;
    font-weight: bold;
    color: #003580;
    margin-bottom: .5rem;
  }
  .instructions {
    font-style: italic;
    color: #555;
    margin-bottom: 1rem;
    padding: .5rem .75rem;
    background: #eef3fb;
    border-left: 4px solid #003580;
    border-radius: 0 4px 4px 0;
  }
  /* Passage */
  .passage-block {
    background: #fafaf7;
    border: 1px solid #ccc;
    border-radius: 4px;
    padding: 1rem 1.25rem;
    margin-bottom: 1.25rem;
  }
  .passage-source {
    font-size: .8rem;
    color: #888;
    margin-bottom: .5rem;
    text-transform: uppercase;
    letter-spacing: .05em;
  }
  .passage-text { white-space: pre-wrap; }
  /* Question */
  .question-block {
    border-top: 1px solid #e0e0e0;
    padding: 1.1rem 0 .5rem;
  }
  .question-number {
    display: inline-block;
    background: #003580;
    color: #fff;
    font-size: .8rem;
    font-weight: bold;
    padding: .15rem .5rem;
    border-radius: 3px;
    margin-bottom: .5rem;
  }
  .question-prompt { font-weight: bold; margin-bottom: .75rem; }
  /* Media */
  .question-image {
    max-width: 100%;
    max-height: 340px;
    display: block;
    margin: .75rem 0;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  audio { display: block; margin: .75rem 0; width: 100%; }
  /* Passage images */
  .passage-images { display: flex; flex-wrap: wrap; gap: .75rem; margin: .75rem 0; }
  .passage-images img {
    max-height: 260px;
    max-width: 100%;
    border: 1px solid #ccc;
    border-radius: 4px;
  }
  /* Answers */
  .answers-toggle {
    margin-top: .75rem;
    border: none;
    border-radius: 4px;
    overflow: hidden;
    background: transparent;
  }
  .answers-toggle summary {
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    gap: .4rem;
    background: #e8f0fe;
    color: #003580;
    border: 1px solid #c3d3f5;
    border-radius: 4px;
    padding: .35rem .85rem;
    font-size: .88rem;
    font-weight: bold;
    list-style: none;
    user-select: none;
  }
  .answers-toggle summary::-webkit-details-marker { display: none; }
  .answers-toggle summary::before { content: "▶"; font-size: .7rem; }
  .answers-toggle[open] summary::before { content: "▼"; }
  .answers-list {
    margin-top: .75rem;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: .4rem;
  }
  .answer-item {
    display: flex;
    align-items: flex-start;
    gap: .6rem;
    padding: .45rem .7rem;
    border-radius: 4px;
    border: 1px solid #e0e0e0;
    background: #fafafa;
  }
  .answer-item.correct {
    background: #eafbea;
    border-color: #5cb85c;
  }
  .answer-choice {
    font-weight: bold;
    min-width: 1.4rem;
    color: #003580;
  }
  .answer-item.correct .answer-choice { color: #2a7a2a; }
  .correct-badge {
    margin-left: auto;
    font-size: .75rem;
    background: #5cb85c;
    color: #fff;
    padding: .1rem .45rem;
    border-radius: 3px;
    white-space: nowrap;
  }
  /* Explanation */
  .explanation-block {
    margin-top: .5rem;
    padding: .5rem .75rem;
    background: #fffbea;
    border-left: 4px solid #f0ad4e;
    border-radius: 0 4px 4px 0;
    font-size: .9rem;
    color: #555;
  }
  /* Difficulty badge */
  .difficulty {
    display: inline-block;
    font-size: .72rem;
    padding: .1rem .45rem;
    border-radius: 3px;
    margin-left: .5rem;
    font-weight: bold;
    vertical-align: middle;
  }
  .difficulty-Easy   { background: #dff0d8; color: #3c763d; }
  .difficulty-Medium { background: #fcf8e3; color: #8a6d3b; }
  .difficulty-Hard   { background: #f2dede; color: #a94442; }
  /* Tags */
  .tags { margin-top: .4rem; display: flex; flex-wrap: wrap; gap: .3rem; }
  .tag {
    font-size: .72rem;
    background: #e9ecef;
    color: #495057;
    border-radius: 3px;
    padding: .1rem .4rem;
  }
  @media print {
    .answers-toggle { display: none; }
  }
</style>
</head>
<body>

<header class="test-header">
  <h1>{{.Title}}</h1>
  {{with .Description}}<p>{{loc .}}</p>{{end}}
  <p style="font-size:.85rem;color:#888;margin-top:.5rem">Version {{.Version}} &nbsp;·&nbsp; {{.CreatedAt.Format "January 2, 2006"}}</p>
</header>

{{range .Sections}}
<div class="section-block">
  <div class="section-title">{{.Type}} Section
    {{- if gt .TimeLimitMins 0}} &nbsp;<span style="font-size:.9rem;font-weight:normal">({{.TimeLimitMins}} minutes)</span>{{end}}
  </div>

  {{range .Parts}}
  <div class="part-block">
    <div class="part-title">{{partName .Number}}</div>
    {{with .Instructions}}{{with loc .}}<div class="instructions">{{.}}</div>{{end}}{{end}}

    {{/* ── Standalone questions (Parts 1, 2, 5) ─────────────────────── */}}
    {{range .Questions}}
    {{template "question" .}}
    {{end}}

    {{/* ── Question groups (Parts 3, 4, 6, 7) ──────────────────────── */}}
    {{range .QuestionGroups}}
    <div class="passage-block">
      {{with .Passage.Source}}<div class="passage-source">{{.}}</div>{{end}}
      {{with .Passage.Audio}}
      <audio controls>
        <source src="{{audioSrc .}}" type="{{mimeAudio .Format}}">
        Your browser does not support the audio element.
      </audio>
      {{end}}
      {{with .Passage.Text}}<div class="passage-text">{{.}}</div>{{end}}
      {{with .Passage.Images}}
      <div class="passage-images">
        {{range .}}<img src="{{imageSrc .}}" alt="{{.AltText}}">{{end}}
      </div>
      {{end}}
    </div>
    {{range .Questions}}
    {{template "question" .}}
    {{end}}
    {{end}}

  </div>{{/* end part */}}
  {{end}}{{/* end range parts */}}
</div>{{/* end section */}}
{{end}}{{/* end range sections */}}

{{define "question"}}
<div class="question-block">
  <span class="question-number">Q{{.Number}}</span>
  <span class="difficulty difficulty-{{.Difficulty}}">{{.Difficulty}}</span>

  {{with .Image}}<img class="question-image" src="{{imageSrc .}}" alt="{{.AltText}}">{{end}}
  {{with .Audio}}
  <audio controls>
    <source src="{{audioSrc .}}" type="{{mimeAudio .Format}}">
    Your browser does not support the audio element.
  </audio>
  {{end}}

  {{with .Prompt}}<div class="question-prompt">{{.}}</div>{{end}}

  {{if .Tags}}
  <div class="tags">{{range .Tags}}<span class="tag">{{.}}</span>{{end}}</div>
  {{end}}

  <details class="answers-toggle">
    <summary>Show answers &amp; explanation</summary>

    <ul class="answers-list">
      {{range .Answers}}
      <li class="answer-item{{if .IsCorrect}} correct{{end}}">
        <span class="answer-choice">({{.Choice}})</span>
        {{if .Audio}}
        <audio controls style="flex:1">
          <source src="{{audioSrc .Audio}}" type="{{mimeAudio .Audio.Format}}">
        </audio>
        {{else}}
        <span style="flex:1">{{.Text}}</span>
        {{end}}
        {{if .IsCorrect}}<span class="correct-badge">✓ Correct</span>{{end}}
      </li>
      {{end}}
    </ul>

    {{/* Per-question explanation */}}
    {{with .Explanation.Text}}{{with loc .}}
    <div class="explanation-block">{{.}}</div>
    {{end}}{{end}}
  </details>
</div>
{{end}}

</body>
</html>`
