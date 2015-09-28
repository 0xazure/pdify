package pdf

import (
	"testing"

	"github.com/jung-kurt/gofpdf"
)

var testPdfDoc = gofpdf.NewCustom(&gofpdf.InitType{
	OrientationStr: "Portrait",
	UnitStr:        "pt",
})

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

func (p *TestPdf) Image(name string, x, y, w, h float64, flow bool, format string, link int, linkUrl string) {
	p.ImageFunc(name, x, y, w, h, flow, format, link, linkUrl)
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
	p := New(testPdfDoc)

	if p.doc == nil {
		t.Errorf("Expected PDF to have an internal representation, got %v", p.doc)
	}
}

func TestPdf_Write(t *testing.T) {
	d := &TestPdf{}

	d.OutputFileAndCloseFunc = func(p string) error {
		return nil
	}

	d.SetMarginsFunc = func(l, t, r float64) {}

	p := New(d)

	var expectedErr error = nil
	actualErr := p.Write("testdata")
	if expectedErr != actualErr {
		t.Errorf("Expected error %v, got %v", expectedErr, actualErr)
	}
}
