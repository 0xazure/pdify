package main

import (
  "fmt"
  "path/filepath"

  "github.com/jung-kurt/gofpdf"
)

const (
  margin = 0
)

func Generate(files []string, output string) error {
  return generate(files, output)
}

func OutputPath(input string, output string, pwd string) (string, error) {
  return outputPath(input, output, pwd)
}

// Internal

func generate(files []string, output string) error {
  pdf := gofpdf.NewCustom(&gofpdf.InitType{
    OrientationStr: "Portrait",
    UnitStr: "pt",
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

func outputPath(input string, output string, pwd string) (string, error) {
  if output == "" {
    outputPath := filepath.Join(pwd, input) + ".pdf"
    return outputPath, nil
  } else {
    if filepath.IsAbs(output) {
      return filepath.Abs(output)
    } else {
      outputPath := filepath.Join(pwd, output)
      return filepath.Abs(outputPath)
    }
  }
}
