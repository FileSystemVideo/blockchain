package util

import (
	"errors"
	"fmt"
	"github.com/otiai10/copy"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)


func DirDelete(paths string) error {
	dir, err := ioutil.ReadDir(paths)
	if err != nil {
		return err
	}
	for _, d := range dir {
		err = os.RemoveAll(path.Join([]string{paths, d.Name()}...))
		if err != nil {
			return err
		}
	}
	return nil
}


func Backup(sourcePath, backupPath, name string) error {
	if sourcePath == backupPath {
		return errors.New("The backup directory cannot be the destination directory")
	}
	st := time.Now()
	stStr := fmt.Sprintf("%d-%d-%d", st.Year(), st.Month(), st.Day())
	dst := filepath.Join(backupPath, fmt.Sprintf(name+"-backup-%s", stStr))
	if err := copy.Copy(sourcePath, dst); err != nil {
		return fmt.Errorf("error while taking data backup: %w", err)
	}
	return nil
}
func RenameBackup(sourcePath, name string) error {
	st := time.Now()
	stStr := fmt.Sprintf("%d-%d-%d", st.Year(), st.Month(), st.Day())
	oldPath := filepath.Join(sourcePath, name)
	dst := filepath.Join(sourcePath, fmt.Sprintf(name+"-backup-%s", stStr))
	if err := os.Rename(oldPath, dst); err != nil {
		return fmt.Errorf("error while taking data backup: %w", err)
	}
	return nil
}


func FilePermissionChange(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	_, err = os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		if os.IsPermission(err) {
			err = os.Chmod(path, 0777)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
