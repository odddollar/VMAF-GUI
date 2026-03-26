package ui

import (
	"errors"
	"os"
)

// Custom entry validator to ensure path entered exists
func validateFileExists(path string) error {
	// Check empty path
	if path == "" {
		return errors.New("path can't be empty")
	}

	// Check path exists
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("file doesn't exist")
		}
		return err
	}

	return nil
}
