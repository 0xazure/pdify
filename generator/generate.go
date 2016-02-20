package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/0xazure/pdify/fs"
	"github.com/0xazure/pdify/pdf"
)

type Generator struct {
	src string
	Pwd string
	Pdf interface {
		AddImage(string) error
		Supports(string) bool
		Write(pdf.File) error
	}
	Walker interface {
		Walk(string, func(fs.FileInfo) bool) ([]string, error)
	}
}

func New(src string) *Generator {
	src, _ = filepath.Abs(src)
	pwd, _ := os.Getwd()

	return &Generator{
		src:    src,
		Pwd:    pwd,
		Pdf:    pdf.New(),
		Walker: new(fs.Walker),
	}
}

func (g *Generator) Generate() ProcessError {
	files, walkErr := g.walk(g.src, g.extFilterFunc())
	if walkErr != nil {
		return newProcessError(walkErr)
	}
	for _, file := range files {
		if err := g.addImage(file); err != nil {
			return newProcessError(err)
		}
	}
	return newProcessError(nil)
}

func (g *Generator) Write(dst string) ProcessError {
	dstPath := destPath(g.src, dst, g.Pwd)
	f, err := os.Create(dstPath)
	if err != nil {
		return newProcessError(err)
	}
	defer f.Close()

	if err := g.Pdf.Write(f); err != nil {
		return newProcessError(err)
	}

	return newProcessError(nil)
}

func (g *Generator) addImage(path string) error {
	return g.Pdf.AddImage(path)
}

func (g *Generator) extFilterFunc() func(fs.FileInfo) bool {
	return func(fi fs.FileInfo) bool {
		if !fi.IsDir() && g.Pdf.Supports(fi.Name()) {
			return true
		}
		return false
	}
}

func (g *Generator) walk(path string, filter func(fs.FileInfo) bool) ([]string, error) {
	return g.Walker.Walk(path, filter)
}

func destPath(src string, dest string, pwd string) (p string) {
	if dest == "" {
		dest = src
	}

	if filepath.IsAbs(dest) {
		p = filepath.Clean(dest)
	} else {
		p = filepath.Join(pwd, dest)
	}

	// Append `.pdf` extension if not already present
	if strings.ToLower(filepath.Ext(p)) != ".pdf" {
		p = p + ".pdf"
	}
	return
}
