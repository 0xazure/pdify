package generator

import (
	"path/filepath"
	"strings"
)

type Generator struct {
	src  string
	dest string
	Pwd  string
	Pdf  interface {
		AddImage(string) error
		Write(string) error
	}
	Walker interface {
		Walk(string, []string) ([]string, error)
	}
}

func (g *Generator) Generate(src string, dest string) ProcessError {
	g.src, _ = filepath.Abs(src)
	g.dest = destPath(g.src, dest, g.Pwd)

	files, _ := g.walk(g.src, g.validExtensions())
	for _, file := range files {
		if err := g.addImage(file); err != nil {
			return newProcessError(err)
		}
	}

	if err := g.write(); err != nil {
		return newProcessError(err)
	}

	return newProcessError(nil)
}

func (g *Generator) addImage(path string) error {
	return g.Pdf.AddImage(path)
}

func (g *Generator) validExtensions() []string {
	// github.com/jung-kurt/gofpdf only supports .png, .jpg, .jpeg images
	// See docs at https://godoc.org/github.com/jung-kurt/gofpdf
	return []string{".png", ".jpg", ".jpeg"}
}

func (g *Generator) write() error {
	return g.Pdf.Write(g.dest)
}

func (g *Generator) walk(path string, validExts []string) ([]string, error) {
	return g.Walker.Walk(path, validExts)
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
