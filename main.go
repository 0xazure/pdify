package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/jung-kurt/gofpdf"

	. "github.com/0xazure/pdify/fs"
	. "github.com/0xazure/pdify/generator"
	"github.com/0xazure/pdify/pdf"
)

const VERSION = "0.1.1"

var (
	input  = kingpin.Arg("input", "Source folder.").Required().String()
	output = kingpin.Flag("output", "Output file name.").Short('o').String()
)

func main() {
	kingpin.Version(VERSION)
	kingpin.Parse()

	p := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "Portrait",
		UnitStr:        "pt",
	})

	w := new(Walker)

	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		panic(pwdErr)
	}

	g := Generator{Pwd: pwd, Pdf: pdf.New(p), Walker: w}
	err := g.Generate(*input, *output)

	if err.ExitCode != 0 {
		fmt.Println(err.Err)
	}
	os.Exit(err.ExitCode)
}
