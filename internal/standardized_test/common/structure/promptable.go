package common

type PromptableComponent interface {
	// GetLLMPrompt returns the text prompt for this component, if applicable.
	GetLLMPrompt() string
}
