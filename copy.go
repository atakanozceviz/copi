package copi

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Copy(src, dest, stp string) error {
	contents, err := contentsToCopy(src, stp)
	if err != nil {
		return err
	}
	// Create folders
	for in, fi := range contents {
		if !fi.IsDir() {
			continue
		}
		pth := strings.TrimPrefix(in, src)
		dir := filepath.Join(dest, pth)
		fmt.Printf("Create: %s\n", dir)
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	wg := &sync.WaitGroup{}
	jobs := make(chan *Job)
	workerCount := len(contents)
	if workerCount > 500 {
		workerCount = 500
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(jobs, wg)
	}

	// Copy files
	for in, fi := range contents {
		if fi.IsDir() {
			continue
		}
		pth := strings.TrimPrefix(in, src)
		out := filepath.Join(dest, pth)
		jobs <- &Job{
			Src: in,
			Dst: out,
		}
		//if err := fcopy(in, out); err != nil {
		//	return err
		//}
	}
	close(jobs)
	wg.Wait()
	return nil
}

func contentsToCopy(src, stp string) (map[string]os.FileInfo, error) {
	config, err := parseSettings(stp)
	if err != nil {
		return nil, err
	}

	srcContents, err := scanContents(src)
	if err != nil {
		return nil, err
	}

	contentsToCopy := make(map[string]os.FileInfo)
CONTENTS:
	for pth, fi := range srcContents {
		pth = strings.Replace(pth, "\\", "/", -1)
		skip := strings.TrimPrefix(pth, src)
		for k := range config {
			if (strings.HasSuffix(k, "/") && strings.HasPrefix(skip, k)) || skip == k {
				if !fi.IsDir() {
					fmt.Printf("Skip: %s\n", skip)
				}
				continue CONTENTS
			}
		}
		contentsToCopy[pth] = fi
	}
	return contentsToCopy, nil
}

func scanContents(dir string) (map[string]os.FileInfo, error) {
	contents := make(map[string]os.FileInfo)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		contents[path] = info
		return nil
	}); err != nil {
		return nil, err
	}
	return contents, nil
}

func fcopy(src, dest string) (err error) {
	// Open src file
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dest, si.Mode())
	if err != nil {
		return
	}

	return
}
