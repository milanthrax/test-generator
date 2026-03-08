package filesystem

import "os"

// CreateDirIfNotExists creates the directory (and any necessary parents)
// if it does not already exist.
func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
