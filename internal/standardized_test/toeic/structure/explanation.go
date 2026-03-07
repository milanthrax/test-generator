package structure

// Explanation provides a correct-answer rationale in multiple languages,
// along with optional grammar points and vocabulary notes.
type Explanation struct {
	// Text is the primary explanation keyed by language.
	Text LocalizedText `json:"text"`
	// GrammarPoints highlights relevant grammar rules, each in multiple languages.
	GrammarPoints []LocalizedText `json:"grammar_points,omitempty"`
	// VocabularyNotes annotates key words or phrases, each in multiple languages.
	VocabularyNotes []LocalizedText `json:"vocabulary_notes,omitempty"`
}
