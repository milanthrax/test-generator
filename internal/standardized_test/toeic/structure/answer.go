package structure

import common "github.com/milanthrax/test-generator/internal/standardized_test/common/structure"

// AnswerChoice is one of the labelled options presented to the test taker.
type AnswerChoice string

const (
	ChoiceA AnswerChoice = "A"
	ChoiceB AnswerChoice = "B"
	ChoiceC AnswerChoice = "C"
	ChoiceD AnswerChoice = "D"
)

// Answer represents a single answer option within a Question.
// For Part 1 and Part 2 the answer text may be empty because the options are
// delivered via audio; in that case Audio should be set.
type Answer struct {
	Choice      AnswerChoice      `json:"choice"`
	Text        string            `json:"text"`            // written answer text (Parts 2-7)
	Audio       *common.AudioFile `json:"audio,omitempty"` // spoken answer option (Parts 1-2)
	IsCorrect   bool              `json:"is_correct"`
	Explanation Explanation       `json:"explanation"` // why this choice is right or wrong
}
