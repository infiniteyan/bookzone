package store

import (
	"os"
	"path/filepath"
	"strings"
)

func DeleteLocalFiles(object ...string) error {
	for _, file := range object {
		os.Remove(strings.TrimLeft(file, "/"))
	}
	return nil
}

func SaveToLocal(tmpfile, save string) (err error) {
	save = strings.TrimLeft(save, "/")
	if strings.HasPrefix(tmpfile, "./") || strings.HasPrefix(save, "./") {
		tmpfile = strings.TrimPrefix(tmpfile, "./")
		save = strings.TrimPrefix(save, "./")
	}
	if strings.ToLower(tmpfile) != strings.ToLower(save) {
		os.MkdirAll(filepath.Dir(save), os.ModePerm)
		err = os.Rename(tmpfile, save)
	}
	return
}
