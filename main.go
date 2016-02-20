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

	g := generator.New(*input)

	err := g.Generate()
	if err.ExitCode != 0 {
		fmt.Println(err.Err)
		os.Exit(err.ExitCode)
	}

	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		fmt.Println(pwdErr)
		os.Exit(1)
	}

	f, osErr := os.Create(generator.FormatPath(*input, *output, pwd))
	if osErr != nil {
		fmt.Println(osErr)
		os.Exit(1)
	}
	defer f.Close()

	err = g.Write(f)
	if err.ExitCode != 0 {
		fmt.Println(err.Err)
		os.Exit(err.ExitCode)
	}

	os.Exit(0)
}
