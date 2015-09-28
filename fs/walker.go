package fs

import (
	"os"
	"path/filepath"
)

type Walker struct{}

func (w *Walker) Walk(path string, validExts []string) (files []string, err error) {
	validExtensionList := NewExtensionList(validExts)

	walkFn := func(path string, fi os.FileInfo, err error) error {
		ext := filepath.Ext(fi.Name())
		if !fi.IsDir() && validExtensionList.Contains(ext) {
			files = append(files, path)
		}
		return nil
	}

	err = filepath.Walk(path, walkFn)

	return
}
