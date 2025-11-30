package utils

import (
	"os"
	"path/filepath"
)

func FindRoot(dir string) string {
	rootMarker := "go.mod"
	for {
		if _, err := os.Stat(filepath.Join(dir, rootMarker)); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
}
