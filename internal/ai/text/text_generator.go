package text

type TextGenerator interface {
	GenerateText(prompt string) (string, error)
}
