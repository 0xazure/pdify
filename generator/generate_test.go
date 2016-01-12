package generator

import (
	"errors"
	"os"
	"testing"
)

type TestPdf struct {
	AddImageFunc        func(string) error
	ValidExtensionsFunc func() []string
	WriteFunc           func(string) error
}

func (p *TestPdf) AddImage(path string) error {
	return p.AddImageFunc(path)
}

func (p *TestPdf) ValidExtensions() []string {
	return p.ValidExtensionsFunc()
}

func (p *TestPdf) Write(dest string) error {
	return p.WriteFunc(dest)
}

type TestWalker struct {
	WalkFunc func(string, []string) ([]string, error)
}

func (w *TestWalker) Walk(path string, validExts []string) ([]string, error) {
	return w.WalkFunc(path, validExts)
}

func TestGenerator_Generate(t *testing.T) {
	p := &TestPdf{}

	p.ValidExtensionsFunc = func() []string {
		return []string{".png", ".jpg", ".jpeg"}
	}

	var imageCount int
	addImageFuncNoErr := func(p string) error {
		imageCount++
		return nil
	}

	addImageFuncErr := func(p string) error {
		return errors.New("Unable to add image")
	}

	var dest string
	writeFuncNoErr := func(d string) error {
		dest = d
		return nil
	}

	writeFuncErr := func(d string) error {
		return errors.New("Unable to write file")
	}

	files := []string{
		"hello",
		"how",
		"are",
		"you",
	}

	w := &TestWalker{}
	w.WalkFunc = func(p string, e []string) ([]string, error) {
		return files, nil
	}

	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		panic(pwdErr)
	}

	generator := Generator{Pwd: pwd, Pdf: p, Walker: w}

	p.AddImageFunc = addImageFuncNoErr
	p.WriteFunc = writeFuncNoErr
	err := generator.Generate("src", "dest")

	if imageCount != len(files) {
		t.Errorf("Expected %d images, got %d", len(files), imageCount)
	}

	if err.Err != nil {
		t.Errorf("Expected no error, got %v", err.Err)
	}

	p.AddImageFunc = addImageFuncErr
	p.WriteFunc = writeFuncNoErr
	err = generator.Generate("src", "dest")

	if err.Err == nil {
		t.Error("Expected error return from Generate, error adding image")
	}

	p.AddImageFunc = addImageFuncNoErr
	p.WriteFunc = writeFuncErr
	err = generator.Generate("src", "dest")

	if err.Err == nil {
		t.Error("Expected error return from Generate, error writing file")
	}
}

type DestPathTest struct {
	src      string
	dest     string
	pwd      string
	expected string
}

var destPathTests = []DestPathTest{
	{
		src:      "album",
		dest:     "",
		pwd:      "/User/user/images",
		expected: "/User/user/images/album.pdf",
	},
	{
		src:      "album",
		dest:     "my_album",
		pwd:      "/User/user/images",
		expected: "/User/user/images/my_album.pdf",
	},
	{
		src:      "album",
		dest:     "my_album.pdf",
		pwd:      "/User/user/images",
		expected: "/User/user/images/my_album.pdf",
	},
	{
		src:      "/User/user/images/album",
		dest:     "",
		pwd:      "/User/user/images",
		expected: "/User/user/images/album.pdf",
	},
	{
		src:      "/User/user/images/album",
		dest:     "my_album",
		pwd:      "/User/user/images",
		expected: "/User/user/images/my_album.pdf",
	},
	{
		src:      "/User/user/images/album",
		dest:     "my_album.pdf",
		pwd:      "/User/user/images",
		expected: "/User/user/images/my_album.pdf",
	},
	{
		src:      "/User/user/images/album",
		dest:     "/User/user/documents/my_album.pdf",
		pwd:      "/User/user/images",
		expected: "/User/user/documents/my_album.pdf",
	},
	{
		src:      "/User/user/images/this is an album",
		dest:     "",
		pwd:      "/User/user/images",
		expected: "/User/user/images/this is an album.pdf",
	},
	{
		src:      "album",
		dest:     "../documents/comics/../pdf/album",
		pwd:      "/User/user/images",
		expected: "/User/user/documents/pdf/album.pdf",
	},
}

func TestGenerator_destPath(t *testing.T) {
	for _, tt := range destPathTests {
		src := tt.src
		dest := tt.dest
		pwd := tt.pwd
		expected := tt.expected

		destPath := destPath(src, dest, pwd)

		if destPath != expected {
			t.Errorf("Expected path '%s', got '%s'", expected, destPath)
		}
	}
}

func TestGenerator_validExtensions(t *testing.T) {
	// github.com/jung-kurt/gofpdf only supports .png, .jpg, .jpeg images
	// See docs at https://godoc.org/github.com/jung-kurt/gofpdf
	supportedExts := []string{".png", ".jpg", ".jpeg"}

	g := Generator{}
	exts := g.validExtensions()

	// Check that the map has the same number of elements as is expected
	if len(supportedExts) != len(exts) {
		t.Errorf("Expected array length of %d, got %d", len(supportedExts), len(exts))
	}

	// Check that all expected values are contained in the array
	for _, supported := range supportedExts {
		found := false
		for _, ext := range exts {
			if supported == ext {
				found = true
			}
		}
		if !found {
			t.Errorf("'%s' should be a valid extension", supported)
		}
	}
}
