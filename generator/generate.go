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

func (g *Generator) Write(f pdf.File) ProcessError {
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

func FormatPath(src string, dst string, pwd string) string {
	if dst == "" {
		dst = src
	}

	var p string
	if filepath.IsAbs(dst) {
		p = filepath.Clean(dst)
	} else {
		p = filepath.Join(pwd, dst)
	}

	// Append `.pdf` extension if not already present
	if strings.ToLower(filepath.Ext(p)) != ".pdf" {
		p = p + ".pdf"
	}

	return p
}
