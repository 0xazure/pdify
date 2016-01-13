package pdf

import (
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

const (
	margin = 0
)

type Pdf struct {
	doc                 pdfDoc
	supportedExtensions map[string]struct{}
}

type pdfDoc interface {
	AddPageFormat(string, gofpdf.SizeType)
	Close()
	Err() bool
	Error() error
	Image(string, float64, float64, float64, float64, bool, string, int, string)
	OutputFileAndClose(string) error
	RegisterImage(string, string) *gofpdf.ImageInfoType
	SetMargins(float64, float64, float64)
}

func New(d pdfDoc) (p *Pdf) {
	p = new(Pdf)
	p.doc = d
	p.doc.SetMargins(margin, margin, margin)
	// http://stackoverflow.com/questions/10485743/contains-method-for-a-slice
	p.supportedExtensions = map[string]struct{}{
		".png":  {},
		".jpg":  {},
		".jpeg": {},
	}
	return
}

func (p *Pdf) AddImage(file string) error {
	imgInfo := p.doc.RegisterImage(file, "")
	w, h := imgInfo.Extent()

	p.doc.AddPageFormat("Portrait", gofpdf.SizeType{Wd: w, Ht: h})
	p.doc.Image(file, 0, 0, w, h, false, "", 0, "")

	if p.doc.Err() {
		p.doc.Close()
		return fmt.Errorf("Error adding '%s' to PDF, %v", file, p.doc.Error())
	}

	return nil
}

func (p *Pdf) SupportsExtension(extension string) (contains bool) {
	ext := formatExtension(extension)
	_, contains = p.supportedExtensions[ext]
	return
}

func (p *Pdf) Write(dest string) error {
	fmt.Printf("Writing PDF to %s\n", dest)
	return p.doc.OutputFileAndClose(dest)
}

func formatExtension(extension string) string {
	ext := strings.ToLower(extension)
	pos := strings.Index(ext, ".")

	if pos != 0 {
		if len(strings.TrimSpace(ext)) > 0 {
			return "." + ext
		}
	}
	return ext
}
