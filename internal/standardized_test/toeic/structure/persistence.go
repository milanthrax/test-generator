package structure

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Marshal serializes a Test to indented JSON.
func (t *Test) Marshal() ([]byte, error) {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal test: %w", err)
	}
	return data, nil
}

// Unmarshal deserializes JSON into an existing Test value.
func (t *Test) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, t); err != nil {
		return fmt.Errorf("unmarshal test: %w", err)
	}
	return nil
}

// UnmarshalTest deserializes JSON and returns a new *Test.
func UnmarshalTest(data []byte) (*Test, error) {
	var t Test
	if err := t.Unmarshal(data); err != nil {
		return nil, err
	}
	return &t, nil
}

// Save writes the Test as JSON to the given file path, creating or
// overwriting the file. UpdatedAt is refreshed before writing.
func (t *Test) Save(path string) error {
	t.UpdatedAt = time.Now().UTC()

	data, err := t.Marshal()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("save test to %q: %w", path, err)
	}
	return nil
}

// LoadTest reads a JSON file from path and deserializes it into a *Test.
func LoadTest(path string) (*Test, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load test from %q: %w", path, err)
	}
	return UnmarshalTest(data)
}
