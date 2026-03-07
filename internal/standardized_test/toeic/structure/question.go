package structure

import common "github.com/milanthrax/test-generator/internal/standardized_test/common/structure"

// DifficultyLevel is a coarse difficulty rating for a question.
type DifficultyLevel string

const (
	Easy   DifficultyLevel = "Easy"
	Medium DifficultyLevel = "Medium"
	Hard   DifficultyLevel = "Hard"
)

// Question represents one discrete TOEIC question.
type Question struct {
	ID            string            `json:"id"`
	Number        int               `json:"number"` // position within the test (1-200)
	Part          PartNumber        `json:"part"`
	Prompt        string            `json:"prompt"`          // the question stem shown / read to the test taker
	Audio         *common.AudioFile `json:"audio,omitempty"` // stimulus audio (Parts 1-4)
	Image         *common.ImageFile `json:"image,omitempty"` // stimulus photograph (Part 1)
	Answers       []Answer          `json:"answers"`         // typically four options (A-D); Part 2 uses three (A-C)
	CorrectAnswer AnswerChoice      `json:"correct_answer"`
	Explanation   Explanation       `json:"explanation"` // overall rationale for the correct answer
	Difficulty    DifficultyLevel   `json:"difficulty"`
	Tags          []string          `json:"tags,omitempty"` // e.g. ["grammar", "prepositions", "business-vocabulary"]
}

// Passage is a shared reading or listening stimulus used by a QuestionGroup.
type Passage struct {
	ID     string             `json:"id"`
	Text   string             `json:"text,omitempty"`   // written text (Parts 6-7); may be empty for audio-only passages
	Audio  *common.AudioFile  `json:"audio,omitempty"`  // spoken stimulus (Parts 3-4)
	Images []common.ImageFile `json:"images,omitempty"` // charts, advertisements, forms, etc. (Part 7)
	Source string             `json:"source,omitempty"` // optional provenance note (e.g. "email", "advertisement")
}

// QuestionGroup bundles a shared Passage with the Questions that refer to it.
// Used for Parts 3, 4, 6, and 7.
type QuestionGroup struct {
	ID        string     `json:"id"`
	Passage   Passage    `json:"passage"`
	Questions []Question `json:"questions"`
}
