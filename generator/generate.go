package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/0xazure/pdify/pdf"
)

type FileInfo interface {
	Name() string
	IsDir() bool
}

type Walker struct{}

func (w *Walker) Walk(path string, filter func(FileInfo) bool) (files []string, err error) {
	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filter(fi) {
			files = append(files, path)
		}
		return nil
	}

	err = filepath.Walk(path, walkFn)

	return
}

type Generator struct {
	src  string
	dest string
	Pwd  string
	Pdf  interface {
		AddImage(string) error
		SupportsExtension(string) bool
		Write(string) error
	}
	Walker interface {
		Walk(string, func(FileInfo) bool) ([]string, error)
	}
}

func New(src string) *Generator {
	src, _ = filepath.Abs(src)
	pwd, _ := os.Getwd()

	return &Generator{
		src:    src,
		Pwd:    pwd,
		Pdf:    pdf.New(),
		Walker: new(Walker),
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

func (g *Generator) Write(dest string) ProcessError {
	g.dest = destPath(g.src, dest, g.Pwd)

	if err := g.write(); err != nil {
		return newProcessError(err)
	}

	return newProcessError(nil)
}

func (g *Generator) addImage(path string) error {
	return g.Pdf.AddImage(path)
}

func (g *Generator) extFilterFunc() func(FileInfo) bool {
	return func(fi FileInfo) bool {
		ext := filepath.Ext(fi.Name())
		if !fi.IsDir() && g.Pdf.SupportsExtension(ext) {
			return true
		}
		return false
	}
}

func (g *Generator) write() error {
	return g.Pdf.Write(g.dest)
}

func (g *Generator) walk(path string, filter func(FileInfo) bool) ([]string, error) {
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
