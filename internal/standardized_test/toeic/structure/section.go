package structure

// SectionType distinguishes the two major sections of the TOEIC.
type SectionType string

const (
	ListeningSection SectionType = "Listening"
	ReadingSection   SectionType = "Reading"
)

// Part represents one of the seven numbered parts of the TOEIC,
// nested within a Section.
type Part struct {
	Number         PartNumber      `json:"number"`
	SectionType    SectionType     `json:"section_type"`
	Instructions   LocalizedText   `json:"instructions"`              // directions shown to the test taker, in multiple languages
	Questions      []Question      `json:"questions,omitempty"`       // standalone questions (Parts 1, 2, 5)
	QuestionGroups []QuestionGroup `json:"question_groups,omitempty"` // grouped questions sharing a passage (Parts 3, 4, 6, 7)
}

// Section holds all Parts that belong to either the Listening or Reading section.
type Section struct {
	Type          SectionType `json:"type"`
	Parts         []Part      `json:"parts"`
	TimeLimitMins int         `json:"time_limit_mins"` // 45 for Listening, 75 for Reading
}
