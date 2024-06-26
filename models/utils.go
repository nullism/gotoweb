package models

import (
	"fmt"
	"os"
	"path/filepath"
)

// getConfigPath returns the absolute path to config.yaml
func getConfigPath(relCount int) (string, error) {
	relative := "."
	for range relCount {
		relative = filepath.Join(relative, "..")
	}

	p := filepath.Join(relative, "config.yaml")
	_, err := os.Stat(p)
	if err != nil {
		if relCount > 5 {
			return "", fmt.Errorf("could not find config.yaml in self or parent directories")
		}
		return getConfigPath(relCount + 1)
	}
	fp, err := filepath.Abs(p)
	if err != nil {
		return "", fmt.Errorf("could not get absolute path: %w", err)
	}
	return fp, nil
}
