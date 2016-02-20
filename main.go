package main

import (
	"log"
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
	logger := log.New(os.Stderr, "pdify: ", 0)

	kingpin.Version(VERSION)
	kingpin.Parse()

	g := generator.New(*input)

	err := g.Generate()
	if err != nil {
		logger.Fatal(err)
	}

	err = g.Write(*output)
	if err != nil {
		logger.Fatal(err)
	}

	os.Exit(0)
}
