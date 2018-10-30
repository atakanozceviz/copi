package copi

import (
	"fmt"
	"os"
	"path/filepath"
)

// Transform files specified in the list
func Transform(dst string, list map[string]string) error {
	for old, new := range list {
		if !filepath.IsAbs(old) {
			old = filepath.Clean(filepath.Join(dst, old))
		}
		if !filepath.IsAbs(new) {
			new = filepath.Clean(filepath.Join(dst, new))
		}
		fmt.Printf("Move: %s to %s\n", old, new)
		err := os.Rename(old, new)
		if err != nil {
			return err
		}
	}
	return nil
}
