package fs

import (
	"os"
	"path/filepath"
)

type FileInfo interface {
	Name() string
	IsDir() bool
}

type Walker struct{}

func (w *Walker) Walk(path string, filter func(FileInfo) bool) (files []string, err error) {
	walkFn := func(path string, fi os.FileInfo, err error) error {
		if filter(fi) {
			files = append(files, path)
		}
		return nil
	}

	err = filepath.Walk(path, walkFn)

	return
}
