package main

import (
	"flag"
	"log"
	"os"

	"github.com/0xazure/pdify/generator"
)

const VERSION = "0.2.0"

func main() {
	logger := log.New(os.Stderr, "pdify: ", 0)

	input := flag.Arg(0)
	output := flag.String("output", "", "specify output `file` name")
	flag.Parse()

	if len(input) == 0 {
		logger.Fatal("missing argument SOURCE")
	}

	g := generator.New(input)

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
