package structure

import "time"

// Test represents a complete TOEIC practice test.
type Test struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description LocalizedText `json:"description"`
	Sections    []Section     `json:"sections"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Version     string        `json:"version"`
}

func (t *Test) GetPrompt() string {
	return ""
}
