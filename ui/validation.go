package ui

import (
	"errors"
	"os"

	"fyne.io/fyne/v2"
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

// Disables start button if paths are invalid
func (u *Ui) validatePathEntries() {
	referenceErr := u.referenceEntry.Validate()
	distortedErr := u.distortedEntry.Validate()

	if referenceErr == nil && distortedErr == nil {
		fyne.Do(func() { u.startButton.Enable() })
	} else {
		fyne.Do(func() { u.startButton.Disable() })
	}
}
