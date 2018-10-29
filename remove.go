package copi

import (
	"fmt"
	"os"
	"strings"
)

func RemoveContents(dst, stp string) error {
	config, err := parseSettings(stp)
	if err != nil {
		return err
	}

	srcContents, err := scanContents(dst)
	if err != nil {
		return err
	}

CONTENTS:
	for pth, fi := range srcContents {
		pth = strings.Replace(pth, "\\", "/", -1)
		skip := strings.TrimPrefix(pth, dst)
		for k := range config {
			if skip == strings.TrimSuffix(k, "/") || (strings.HasSuffix(k, "/") && strings.HasPrefix(skip, k)) || skip == k {
				continue CONTENTS
			}
		}

		// fmt.Printf("Remove: %s\n", pth)
		// err := os.RemoveAll(pth)
		// if err != nil {
		// 	return err
		// }
		// TODO: cannot remove folders, leaves empty folders
		if !fi.IsDir() {
			fmt.Printf("Remove: %s\n", pth)
			err := os.Remove(pth)
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

// func RemoveContents(dst, stp string) error {
// 	config, err := parseSettings(stp)
// 	if err != nil {
// 		return err
// 	}

// 	if err := filepath.Walk(dst, func(path string, info os.FileInfo, err error) error {
// 		if os.IsNotExist(err) {
// 			return nil
// 		}
// 		if err != nil {
// 			return err
// 		}

// 		if info.IsDir() && info.Name() == filepath.Base(dst) {
// 			return nil
// 		}
// 		path = strings.Replace(path, "\\", "/", -1)
// 		if info.IsDir() {
// 			for skip := range config {
// 				skip := strings.TrimPrefix(skip, dst)
// 				if strings.HasSuffix(skip, "/") && strings.TrimSuffix(skip, "/") == info.Name() || path == skip {
// 					fmt.Printf("Skip Dir: %s\n", path)
// 					return filepath.SkipDir
// 				}
// 				if strings.HasSuffix(skip, "/") && strings.HasPrefix(skip, info.Name()) {
// 					fmt.Printf("Skip Dir: %s\n", path)
// 					return nil
// 				}
// 				continue
// 			}
// 		}

// 		for skip := range config {
// 			if !strings.HasSuffix(skip, "/") && filepath.Base(skip) == info.Name() {
// 				fmt.Printf("Skip File: %s\n", path)
// 				return nil
// 			}
// 			continue
// 		}
// 		fmt.Printf("Remove: %s\n", path)
// 		return os.RemoveAll(path)
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }
