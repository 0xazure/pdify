package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

const (
	margin = 0
)

func Generate(files []string, output string) error {
	return generate(files, output)
}

func OutputPath(input string, output string, pwd string) string {
	return outputPath(input, output, pwd)
}

// Internal

func generate(files []string, output string) error {
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "Portrait",
		UnitStr:        "pt",
	})
	pdf.SetMargins(margin, margin, margin)

	for _, file := range files {
		imgInfo := pdf.RegisterImage(file, "")
		w, h := imgInfo.Extent()

		if pdf.Err() {
			pdf.Close()
			return fmt.Errorf("Error adding %s to PDF, %s", file, pdf.Error())
		}

		pdf.AddPageFormat("Portrait", gofpdf.SizeType{Wd: w, Ht: h})
		pdf.Image(file, 0, 0, w, h, false, "", 0, "")
	}

	fmt.Printf("Writing PDF to %s\n", output)
	return pdf.OutputFileAndClose(output)
}

func outputPath(input string, output string, pwd string) string {
	if output == "" {
		output = input
	}

	var opath string

	if filepath.IsAbs(output) {
		opath = filepath.Clean(output)
	} else {
		opath = filepath.Join(pwd, output)
	}

	if strings.ToLower(filepath.Ext(opath)) != ".pdf" {
		opath = opath + ".pdf"
	}

	return opath
}
