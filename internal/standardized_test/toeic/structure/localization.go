package structure

// LanguageCode is an ISO 639-1 language identifier (e.g. "en", "ja", "ko").
type LanguageCode string

const (
	EN LanguageCode = "en"
	JA LanguageCode = "ja"
	KO LanguageCode = "ko"
	ZH LanguageCode = "zh"
	ES LanguageCode = "es"
	FR LanguageCode = "fr"
	PT LanguageCode = "pt"
	VI LanguageCode = "vi"
	TH LanguageCode = "th"
	ID LanguageCode = "id"
)

// LocalizedText maps a language code to its translated string.
type LocalizedText map[LanguageCode]string
