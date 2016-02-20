package pdf

import (
	"fmt"
	"io"
	"path/filepath"

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
	Output(io.Writer) error
	RegisterImage(string, string) *gofpdf.ImageInfoType
	SetMargins(float64, float64, float64)
}

type File interface {
	Name() string
	Write([]byte) (int, error)
}

func New() (p *Pdf) {
	p = &Pdf{
		doc: gofpdf.NewCustom(&gofpdf.InitType{
			OrientationStr: "Portrait",
			UnitStr:        "pt",
		}),
		// http://stackoverflow.com/questions/10485743/contains-method-for-a-slice
		supportedExtensions: map[string]struct{}{
			".png":  {},
			".jpg":  {},
			".jpeg": {},
		},
	}
	p.doc.SetMargins(margin, margin, margin)
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

func (p *Pdf) Supports(path string) bool {
	_, contains := p.supportedExtensions[filepath.Ext(path)]
	return contains
}

func (p *Pdf) Write(f File) error {
	path, _ := filepath.Abs(f.Name())
	fmt.Printf("Writing PDF to %s\n", path)
	return p.doc.Output(f)
}
