package pdf

import (
	"testing"

	"github.com/jung-kurt/gofpdf"
)

type TestPdf struct {
	AddPageFormatFunc      func(string, gofpdf.SizeType)
	CloseFunc              func()
	ErrFunc                func() bool
	ErrorFunc              func() error
	ImageFunc              func(string, float64, float64, float64, float64, bool, string, int, string)
	OutputFileAndCloseFunc func(string) error
	RegisterImageFunc      func(string, string) *gofpdf.ImageInfoType
	SetMarginsFunc         func(float64, float64, float64)
}

func (p *TestPdf) AddPageFormat(orientation string, size gofpdf.SizeType) {
	p.AddPageFormatFunc(orientation, size)
}

func (p *TestPdf) Close() {
	p.CloseFunc()
}

func (p *TestPdf) Err() bool {
	return p.ErrFunc()
}

func (p *TestPdf) Error() error {
	return p.ErrorFunc()
}

func (p *TestPdf) Image(name string, x, y, w, h float64, flow bool, format string, link int, linkURL string) {
	p.ImageFunc(name, x, y, w, h, flow, format, link, linkURL)
}

func (p *TestPdf) OutputFileAndClose(path string) error {
	return p.OutputFileAndCloseFunc(path)
}

func (p *TestPdf) RegisterImage(path string, format string) *gofpdf.ImageInfoType {
	return p.RegisterImageFunc(path, format)
}

func (p *TestPdf) SetMargins(left, top, right float64) {
	p.SetMarginsFunc(left, top, right)
}

func TestPdf_New(t *testing.T) {
	p := New()

	if p.doc == nil {
		t.Errorf("Expected PDF to have an internal representation, got %v", p.doc)
	}
}

type SupportedFileTest struct {
	path      string
	supported bool
}

// github.com/jung-kurt/gofpdf only supports .png, .jpg, .jpeg images
// See docs at https://godoc.org/github.com/jung-kurt/gofpdf
var supportedFileTests = []SupportedFileTest{
	{"test.png", true},
	{"test.jpg", true},
	{"test.jpeg", true},
	{"test.jpg.tar", false},
	{"", false},
	{"test", false},
	{"test.exe", false},
}

func TestPdf_Supports(t *testing.T) {
	p := New()

	for _, tt := range supportedFileTests {
		path := tt.path
		expected := tt.supported

		support := p.Supports(path)
		if support != expected {
			t.Errorf("Expected support for %s to be %t, got %t", path, expected, support)
		}
	}
}

func TestPdf_Write(t *testing.T) {
	d := &TestPdf{}

	d.OutputFileAndCloseFunc = func(p string) error {
		return nil
	}

	d.SetMarginsFunc = func(l, t, r float64) {}

	p := &Pdf{
		doc: d,
	}

	var expectedErr error
	actualErr := p.Write("testdata")
	if expectedErr != actualErr {
		t.Errorf("Expected error %v, got %v", expectedErr, actualErr)
	}
}
