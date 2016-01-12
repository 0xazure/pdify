package fs

import (
	"strings"
)

type ExtensionList struct {
	m map[string]struct{}
}

func NewExtensionList(exts []string) ExtensionList {
	l := ExtensionList{}
	cExts := cleanExts(exts)
	l.m = buildMap(cExts)
	return l
}

func (l *ExtensionList) Exts() (exts []string) {
	for k := range l.m {
		exts = append(exts, k)
	}
	return
}

func (l *ExtensionList) Contains(ext string) bool {
	_, contains := l.m[ext]
	return contains
}

// http://stackoverflow.com/questions/10485743/contains-method-for-a-slice
func buildMap(exts []string) (m map[string]struct{}) {
	m = make(map[string]struct{}, len(exts))
	for _, ext := range exts {
		m[ext] = struct{}{}
	}
	return
}

func cleanExts(exts []string) (cExts []string) {
	for _, ext := range exts {
		ext := strings.ToLower(ext)
		pos := strings.Index(ext, ".")

		if pos == 0 {
			cExts = append(cExts, ext)
		} else {
			if len(strings.TrimSpace(ext)) > 0 {
				cExts = append(cExts, "."+ext)
			}
		}
	}
	return
}
