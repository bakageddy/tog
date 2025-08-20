package util

import (
	"errors"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	} else if errors.Is(os.ErrNotExist, err) {
		return false
	} else {
		return true
	}
}

func PathCanonicalize(path string) string {
	if result, err := filepath.Abs(path); err != nil {
		return ""
	} else {
		return result
	}
}
