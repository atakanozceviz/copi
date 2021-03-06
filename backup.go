package copi

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

// Backup things from src to dst
func Backup(src, dst string, keep int) error {
	if err := cleanBackupDir(filepath.Base(src), dst, keep); err != nil {
		return err
	}

	now := strconv.FormatInt(time.Now().Unix(), 10)
	dst = filepath.Join(dst, now+"-"+filepath.Base(src))

	fmt.Printf("Backup: %s\n", dst)
	err := copyDir(src, dst)
	if err != nil {
		return fmt.Errorf("cannot copy directory: %v", err)
	}
	return nil
}

func cleanBackupDir(name, dst string, keep int) error {
	re := regexp.MustCompile(`^[0-9]{10}-` + name)
	backups := make([]os.FileInfo, 0, 5)

	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}
	entries, err := ioutil.ReadDir(dst)
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
		err := os.RemoveAll(filepath.Join(dst, backups[i].Name()))
		if err != nil {
			return err
		}
		count--
	}
	return nil
}
