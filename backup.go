package copi

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func Backup(src, dest string, keep int) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if !filepath.IsAbs(src) {
		src = filepath.Join(wd, src)
	}
	if !filepath.IsAbs(dest) {
		dest = filepath.Join(wd, dest)
	}
	if err := cleanBackup(dest, keep); err != nil {
		return err
	}
	now := strconv.FormatInt(time.Now().Unix(), 10)
	dest = filepath.Join(dest, now+"-"+filepath.Base(src))
	return copyDir(src, dest)
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
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

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func cleanBackup(dest string, keep int) error {
	re := regexp.MustCompile(`^[0-9]{10}-.+`)
	backups := make([]os.FileInfo, 0, 5)

	entries, err := ioutil.ReadDir(dest)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() && re.MatchString(entry.Name()) {
			backups = append(backups, entry)
		}
	}
	count := len(backups)
	for i := 0; count >= keep; i++ {
		err := os.RemoveAll(filepath.Join(dest, backups[i].Name()))
		if err != nil {
			return err
		}
		count--
	}
	return nil
}