package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"

	"github.com/0xazure/pdify/generator"
)

const VERSION = "0.2.0"

var (
	input  = kingpin.Arg("input", "Source folder.").Required().String()
	output = kingpin.Flag("output", "Output file name.").Short('o').String()
)

func main() {
	kingpin.Version(VERSION)
	kingpin.Parse()

	g := generator.New()
	err := g.Generate(*input, *output)

	if err.ExitCode != 0 {
		fmt.Println(err.Err)
	}
	os.Exit(err.ExitCode)
}
