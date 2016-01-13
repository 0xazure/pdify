package generator

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"

	"github.com/0xazure/pdify/fs"
	"github.com/0xazure/pdify/pdf"
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
		Walk(string, func(fs.FileInfo) bool) ([]string, error)
	}
}

func New(src string) *Generator {
	src, _ = filepath.Abs(src)
	pwd, _ := os.Getwd()
	gofpdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "Portrait",
		UnitStr:        "pt",
	})

	return &Generator{
		src:    src,
		Pwd:    pwd,
		Pdf:    pdf.New(gofpdf),
		Walker: new(fs.Walker),
	}
}

func (g *Generator) Generate() ProcessError {
	files, _ := g.walk(g.src, g.extFilterFunc())
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

func (g *Generator) extFilterFunc() func(fs.FileInfo) bool {
	validExtensionList := fs.NewExtensionList(g.validExtensions())

	return func(fi fs.FileInfo) bool {
		ext := filepath.Ext(fi.Name())
		if !fi.IsDir() && validExtensionList.Contains(ext) {
			return true
		}
		return false
	}
}

func (g *Generator) validExtensions() []string {
	// github.com/jung-kurt/gofpdf only supports .png, .jpg, .jpeg images
	// See docs at https://godoc.org/github.com/jung-kurt/gofpdf
	return []string{".png", ".jpg", ".jpeg"}
}

func (g *Generator) write() error {
	return g.Pdf.Write(g.dest)
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
