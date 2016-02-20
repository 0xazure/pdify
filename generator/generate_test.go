package generator

import (
	"errors"
	"os"
	"testing"

	"github.com/0xazure/pdify/fs"
	"github.com/0xazure/pdify/pdf"
)

type TestPdf struct {
	AddImageFunc func(string) error
	SupportsFunc func(string) bool
	WriteFunc    func(pdf.File) error
}

func (p *TestPdf) AddImage(path string) error {
	return p.AddImageFunc(path)
}

func (p *TestPdf) Supports(name string) bool {
	return p.SupportsFunc(name)
}
func (p *TestPdf) Write(f pdf.File) error {
	return p.WriteFunc(f)
}

type TestWalker struct {
	WalkFunc func(string, func(fs.FileInfo) bool) ([]string, error)
}

func (w *TestWalker) Walk(path string, filter func(fs.FileInfo) bool) ([]string, error) {
	return w.WalkFunc(path, filter)
}

func TestGenerator_New(t *testing.T) {
	pwd, _ := os.Getwd()
	g := New("src")

	if g.Pwd != pwd {
		t.Errorf("Expected pwd %s, got %s", pwd, g.Pwd)
	}
}

type FormatPathTest struct {
	src      string
	dst      string
	pwd      string
	expected string
}

var formatPathTests = []FormatPathTest{
	{
		src:      "album",
		dst:      "",
		pwd:      "/User/user/images",
		expected: "/User/user/images/album.pdf",
	},
	{
		src:      "album",
		dst:      "my_album",
		pwd:      "/User/user/images",
		expected: "/User/user/images/my_album.pdf",
	},
	{
		src:      "album",
		dst:      "my_album.pdf",
		pwd:      "/User/user/images",
		expected: "/User/user/images/my_album.pdf",
	},
	{
		src:      "/User/user/images/album",
		dst:      "",
		pwd:      "/User/user/images",
		expected: "/User/user/images/album.pdf",
	},
	{
		src:      "/User/user/images/album",
		dst:      "my_album",
		pwd:      "/User/user/images",
		expected: "/User/user/images/my_album.pdf",
	},
	{
		src:      "/User/user/images/album",
		dst:      "my_album.pdf",
		pwd:      "/User/user/images",
		expected: "/User/user/images/my_album.pdf",
	},
	{
		src:      "/User/user/images/album",
		dst:      "/User/user/documents/my_album.pdf",
		pwd:      "/User/user/images",
		expected: "/User/user/documents/my_album.pdf",
	},
	{
		src:      "/User/user/images/this is an album",
		dst:      "",
		pwd:      "/User/user/images",
		expected: "/User/user/images/this is an album.pdf",
	},
	{
		src:      "album",
		dst:      "../documents/comics/../pdf/album",
		pwd:      "/User/user/images",
		expected: "/User/user/documents/pdf/album.pdf",
	},
}

func TestGenerator_FormatPath(t *testing.T) {
	for _, tt := range formatPathTests {
		src := tt.src
		dst := tt.dst
		pwd := tt.pwd
		expected := tt.expected

		formattedPath := FormatPath(src, dst, pwd)

		if formattedPath != expected {
			t.Errorf("Expected path '%s', got '%s'", expected, formattedPath)
		}
	}
}

func TestGenerator_Generate(t *testing.T) {
	p := &TestPdf{}

	var imageCount int
	addImageFuncNoErr := func(p string) error {
		imageCount++
		return nil
	}

	addImageFuncErr := func(p string) error {
		return errors.New("Unable to add image")
	}

	files := []string{
		"hello",
		"how",
		"are",
		"you",
	}

	w := &TestWalker{}

	walkFuncNoErr := func(p string, f func(fs.FileInfo) bool) ([]string, error) {
		return files, nil
	}

	walkFuncErr := func(p string, f func(fs.FileInfo) bool) ([]string, error) {
		return nil, errors.New("Problem walking path")
	}

	generator := Generator{Pdf: p, Walker: w}

	p.AddImageFunc = addImageFuncNoErr
	w.WalkFunc = walkFuncNoErr
	err := generator.Generate()

	if imageCount != len(files) {
		t.Errorf("Expected %d images, got %d", len(files), imageCount)
	}

	if err.Err != nil {
		t.Errorf("Expected no error, got %v", err.Err)
	}

	p.AddImageFunc = addImageFuncErr
	w.WalkFunc = walkFuncNoErr
	err = generator.Generate()

	if err.Err == nil {
		t.Error("Expected error return from Generate, error adding image")
	}

	p.AddImageFunc = addImageFuncNoErr
	w.WalkFunc = walkFuncErr
	err = generator.Generate()

	if err.Err == nil {
		t.Error("Expected error return from Generate, error walking path")
	}
}

type TestFile struct {
	NameFunc  func() string
	WriteFunc func([]byte) (int, error)
}

func (f *TestFile) Name() string {
	return f.NameFunc()
}

func (f *TestFile) Write(b []byte) (int, error) {
	return f.WriteFunc(p)
}

func TestGenerator_Write(t *testing.T) {
	p := &TestPdf{}

	writeFuncNoErr := func(f pdf.File) error {
		return nil
	}

	writeFuncErr := func(f pdf.File) error {
		return errors.New("Error writing PDF")
	}

	generator := Generator{Pdf: p}

	f := &TestFile{
		NameFunc: func() string {
			return "test.pdf"
		},
		WriteFunc: func(b []byte) (int, error) {
			return 0, nil
		},
	}

	p.WriteFunc = writeFuncNoErr
	err := generator.Write(f)
	if err.Err != nil {
		t.Errorf("Expected no error return from Write, got %v", err.Err)
	}

	p.WriteFunc = writeFuncErr
	err = generator.Write(f)
	if err.Err == nil {
		t.Error("Expected error return from Write, error writing PDF")
	}
}

type TestFileInfo struct {
	NameFunc  func() string
	IsDirFunc func() bool
}

func (i *TestFileInfo) Name() string {
	return i.NameFunc()
}

func (i *TestFileInfo) IsDir() bool {
	return i.IsDirFunc()
}

type ExtFilterFuncTest struct {
	nameFunc  func() string
	isDirFunc func() bool
	expected  bool
}

func TestGenerator_extFilterFunc(t *testing.T) {
	validFileExt := "image.jpg"
	nameFuncValidExt := func() string {
		return validFileExt
	}

	invalidFileExt := "image.xyz"
	nameFuncInvalidExt := func() string {
		return invalidFileExt
	}

	isDirFuncIsDir := func() bool {
		return true
	}

	isDirFuncIsNotDir := func() bool {
		return false
	}

	extFilterFuncTests := []ExtFilterFuncTest{
		{
			nameFunc:  nameFuncValidExt,
			isDirFunc: isDirFuncIsNotDir,
			expected:  true,
		},
		{
			nameFunc:  nameFuncValidExt,
			isDirFunc: isDirFuncIsDir,
			expected:  false,
		},
		{
			nameFunc:  nameFuncInvalidExt,
			isDirFunc: isDirFuncIsNotDir,
			expected:  false,
		},
		{
			nameFunc:  nameFuncInvalidExt,
			isDirFunc: isDirFuncIsDir,
			expected:  false,
		},
	}

	g := New("src")
	extFilterFunc := g.extFilterFunc()
	fi := &TestFileInfo{}

	for _, tt := range extFilterFuncTests {
		fi.NameFunc = tt.nameFunc
		fi.IsDirFunc = tt.isDirFunc
		expected := tt.expected

		added := extFilterFunc(fi)

		if added != expected {
			t.Errorf("Expected %s (IsDir: %t) to be added: %t", fi.Name(), fi.IsDir(), expected)
		}
	}
}
