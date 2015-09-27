package main

import (
	"os"
	"path/filepath"
	"strings"
)

// github.com/jung-kurt/gofpdf only supports .png, .jpg, .jpeg images
// See docs at https://godoc.org/github.com/jung-kurt/gofpdf
var imageExtensions = []string{".png", ".jpg", ".jpeg"}
var extMap = ExtensionMapBuilder(imageExtensions)

func Analyze(path string) []string {
	return analyze(path)
}

func ExtensionMapBuilder(extensions []string) map[string]struct{} {
	return extensionMapBuilder(extensions)
}

// Internal

func analyze(path string) []string {
	var files []string

	walkFn := func(path string, fi os.FileInfo, err error) error {
		isValid := isValidExtension(fi.Name(), extMap)
		if !fi.IsDir() && isValid {
			files = append(files, path)
		}
		return nil
	}

	filepath.Walk(path, walkFn)

	return files
}

// http://stackoverflow.com/questions/10485743/contains-method-for-a-slice
func extensionMapBuilder(extensions []string) map[string]struct{} {
	sanitizedExt := sanitizeExtensions(extensions)

	extMap := make(map[string]struct{}, len(sanitizedExt))

	for _, ext := range sanitizedExt {
		extMap[ext] = struct{}{}
	}

	return extMap
}

func isValidExtension(file string, extMap map[string]struct{}) bool {
	ext := strings.ToLower(filepath.Ext(file))
	_, isValid := extMap[ext]

	return isValid
}

func sanitizeExtensions(extensions []string) []string {
	var sanitized []string

	for _, ext := range extensions {
		pos := strings.Index(ext, ".")
		if pos == 0 {
			sanitized = append(sanitized, ext)
		} else {
			if len(strings.TrimSpace(ext)) > 0 {
				sanitized = append(sanitized, "."+ext)
			}
		}
	}

	return sanitized
}
