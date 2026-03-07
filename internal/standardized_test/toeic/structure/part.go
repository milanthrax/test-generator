package structure

import "fmt"

// PartNumber identifies which of the seven TOEIC parts a question belongs to.
type PartNumber int

const (
	Part1Photographs          PartNumber = 1 // Listening - describe a photograph
	Part2QuestionResponse     PartNumber = 2 // Listening - choose the best response
	Part3Conversations        PartNumber = 3 // Listening - short conversations
	Part4Talks                PartNumber = 4 // Listening - short talks
	Part5IncompleteSentences  PartNumber = 5 // Reading   - fill in the blank
	Part6TextCompletion       PartNumber = 6 // Reading   - paragraph with blanks
	Part7ReadingComprehension PartNumber = 7 // Reading   - single / multi-passage
)

// GetInstructions returns the student-facing directions for this part.
func (pn PartNumber) GetInstructions() string {
	switch pn {
	case Part1Photographs:
		return "Look at the photograph and listen to the four statements. " +
			"Choose the statement that best describes what you see in the photograph."
	case Part2QuestionResponse:
		return "Listen to the question and the three response options. " +
			"Choose the best response to the question."
	case Part3Conversations:
		return "Listen to the conversation and answer the three questions that follow."
	case Part4Talks:
		return "Listen to the talk and answer the three questions that follow."
	case Part5IncompleteSentences:
		return "Choose the word or phrase that best completes the sentence."
	case Part6TextCompletion:
		return "Read the text and choose the word or phrase that best completes each blank."
	case Part7ReadingComprehension:
		return "Read the passage(s) and answer the questions that follow."
	default:
		return ""
	}
}

// LLMPrompt returns a prompt to send to a language model asking it to generate
// a realistic TOEIC question for this part. The caller may append additional
// constraints (topic, vocabulary level, difficulty, etc.) before sending.
func (pn PartNumber) LLMPrompt() string {
	base := "You are an expert TOEIC test writer. " +
		"Generate exactly one realistic TOEIC %s question in valid JSON that matches " +
		"the following Go struct schema:\n\n" +
		"  Question { id, number, part, prompt, audio?, image?, answers[], " +
		"correct_answer, explanation, difficulty, tags[] }\n" +
		"  Answer   { choice (A/B/C/D), text, audio?, is_correct, explanation }\n" +
		"  Explanation { text: { \"en\": \"...\" }, grammar_points?, vocabulary_notes? }\n\n" +
		"Rules:\n" +
		"- All answer choices must be plausible; distractors should reflect common learner errors.\n" +
		"- Exactly one answer must have is_correct = true.\n" +
		"- Include an English explanation for the correct answer and for each distractor.\n" +
		"- Set difficulty to one of: Easy, Medium, Hard.\n" +
		"- Include relevant tags (e.g. grammar rule, topic, vocabulary category).\n"

	switch pn {
	case Part1Photographs:
		return fmt.Sprintf(base, "Part 1 (Photographs)") +
			"- The prompt field must be empty (no written question stem; the stimulus is a photo).\n" +
			"- Provide four answer choices (A–D) as short declarative sentences describing the scene.\n" +
			"- The image field should reference a realistic workplace or public-place photograph.\n" +
			"- Audio fields should be omitted (answers are read aloud in the real test, but for generation provide text).\n"
	case Part2QuestionResponse:
		return fmt.Sprintf(base, "Part 2 (Question-Response)") +
			"- The prompt must be a spoken question or statement (write it in the prompt field).\n" +
			"- Provide exactly three answer choices (A–C); each is a short spoken response.\n" +
			"- Distractors should use common traps: similar sounds, associated words, or off-topic replies.\n"
	case Part3Conversations:
		return fmt.Sprintf(base, "Part 3 (Conversations)") +
			"- The prompt must be a question about a short two- or three-person conversation.\n" +
			"- The conversation transcript belongs in the passage audio transcript field.\n" +
			"- Provide four answer choices (A–D).\n" +
			"- Return this question as part of a QuestionGroup of three questions sharing the same passage.\n"
	case Part4Talks:
		return fmt.Sprintf(base, "Part 4 (Short Talks)") +
			"- The prompt must be a question about a monologue (announcement, advertisement, voicemail, etc.).\n" +
			"- The monologue transcript belongs in the passage audio transcript field.\n" +
			"- Provide four answer choices (A–D).\n" +
			"- Return this question as part of a QuestionGroup of three questions sharing the same passage.\n"
	case Part5IncompleteSentences:
		return fmt.Sprintf(base, "Part 5 (Incomplete Sentences)") +
			"- The prompt must be a single sentence with exactly one blank represented by the token _____.\n" +
			"- Provide four answer choices (A–D): one correct word/phrase and three grammatically plausible distractors.\n" +
			"- Focus on grammar (verb tense, prepositions, conjunctions) or vocabulary in a business context.\n"
	case Part6TextCompletion:
		return fmt.Sprintf(base, "Part 6 (Text Completion)") +
			"- The prompt is a question about a specific blank in a short business text (email, memo, article).\n" +
			"- The full text (with the blank as _____(number)) belongs in the passage text field.\n" +
			"- Provide four answer choices (A–D).\n" +
			"- Return this question as part of a QuestionGroup of four questions sharing the same passage.\n"
	case Part7ReadingComprehension:
		return fmt.Sprintf(base, "Part 7 (Reading Comprehension)") +
			"- The prompt must be a comprehension question about the passage (main idea, detail, inference, vocabulary in context).\n" +
			"- The full passage (article, email, advertisement, form, etc.) belongs in the passage text field.\n" +
			"- Provide four answer choices (A–D).\n" +
			"- For double/triple passages, include multiple passage entries in the QuestionGroup.\n" +
			"- Return this question as part of a QuestionGroup of 2–5 questions sharing the same passage(s).\n"
	default:
		return fmt.Sprintf(base, "question")
	}
}
