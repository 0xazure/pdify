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

type FileSystem struct {
	WalkFn func(string, filepath.WalkFunc) error
}

// Wanted to get rid of the `fs` package, but now we have a different fs namespace
func NewFileSystem() *FileSystem {
	walkFn := func(root string, walkFn filepath.WalkFunc) error {
		return filepath.Walk(root, walkFn)
	}
	return &FileSystem{
		WalkFn: walkFn,
	}
}
func (fs *FileSystem) Walk(root string, walkFn filepath.WalkFunc) error {
	return fs.WalkFn(root, walkFn)
}

type Walker struct {
	Fs FileSystem
}

func newWalker() *Walker {
	return &Walker{
		Fs: FileSystem{},
	}
}

// Move this into g.Walk() and eliminate Walker type entierly?
func (w *Walker) Walk(path string, filter func(FileInfo) bool) ([]string, error) {
	var files []string

	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filter(fi) {
			files = append(files, path)
		}
		return nil
	}

	err := w.Fs.Walk(path, walkFn)
	return files, err
}

type Generator struct {
	src  string
	dest string
	Pwd  string
	Pdf  interface {
		AddImage(string) error
		Supports(string) bool
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
		Walker: newWalker(),
	}
}

func (g *Generator) Generate() error {
	files, walkErr := g.walk(g.src, g.extFilterFunc())
	if walkErr != nil {
		return walkErr
	}
	for _, file := range files {
		if err := g.addImage(file); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) Write(dest string) error {
	g.dest = destPath(g.src, dest, g.Pwd)

	if err := g.write(); err != nil {
		return err
	}

	return nil
}

func (g *Generator) addImage(path string) error {
	return g.Pdf.AddImage(path)
}

func (g *Generator) extFilterFunc() func(FileInfo) bool {
	return func(fi FileInfo) bool {
		if !fi.IsDir() && g.Pdf.Supports(fi.Name()) {
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
