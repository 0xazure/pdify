package generator

import (
	"errors"
	"os"
	"testing"

	"github.com/0xazure/pdify/fs"
)

type TestPdf struct {
	AddImageFunc func(string) error
	SupportsFunc func(string) bool
	WriteFunc    func(string) error
}

func (p *TestPdf) AddImage(path string) error {
	return p.AddImageFunc(path)
}

func (p *TestPdf) Supports(name string) bool {
	return p.SupportsFunc(name)
}
func (p *TestPdf) Write(dest string) error {
	return p.WriteFunc(dest)
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

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	p.AddImageFunc = addImageFuncErr
	w.WalkFunc = walkFuncNoErr
	err = generator.Generate()

	if err == nil {
		t.Error("Expected error return from Generate, error adding image")
	}

	p.AddImageFunc = addImageFuncNoErr
	w.WalkFunc = walkFuncErr
	err = generator.Generate()

	if err == nil {
		t.Error("Expected error return from Generate, error walking path")
	}
}

func TestGenerator_Write(t *testing.T) {
	p := &TestPdf{}

	var dest string
	writeFuncNoErr := func(d string) error {
		dest = d
		return nil
	}

	writeFuncErr := func(d string) error {
		return errors.New("Unable to write file")
	}

	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		panic(pwdErr)
	}

	generator := Generator{Pwd: pwd, Pdf: p}

	p.WriteFunc = writeFuncNoErr
	err := generator.Write("dest")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	p.WriteFunc = writeFuncErr
	err = generator.Write("dest")

	if err == nil {
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
