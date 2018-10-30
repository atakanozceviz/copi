package copi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RemoveContentsExcept specified in the list of files and folders to ignore
func RemoveContentsExcept(dst string, list map[string]struct{}) error {
	err := filepath.Walk(dst, func(path string, info os.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.IsDir() && filepath.Base(dst) == info.Name() {
			return nil
		}

		upath := strings.Replace(path, "\\", "/", -1)
		skip := strings.Replace(upath, dst, "", -1)

		if info.IsDir() {
			if _, ok := list[skip+"/"]; ok {
				return filepath.SkipDir
			}
		}
		if _, ok := list[skip]; ok {
			return nil
		}
		for k := range list {
			if strings.HasPrefix(k, skip+"/") {
				return nil
			}
		}

		fmt.Printf("Remove: %s\n", path)
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
