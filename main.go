package main

import (
	"fmt"
	"log"
	"os"

	flag "github.com/ogier/pflag"

	"github.com/0xazure/pdify/generator"
)

const (
	usage = `Usage: pdify [options] directory

Options:`
	VERSION = "0.2.0"
)

func main() {
	logger := log.New(os.Stderr, "pdify: ", 0)

	var help bool
	var output string
	var version bool
	flag.BoolVarP(&help, "help", "h", false, "display this help and exit")
	flag.StringVarP(&output, "output", "o", "", "specify output `file` name")
	flag.BoolVar(&version, "version", false, "display version information and exit")
	flag.Parse()
	input := flag.Arg(0)

	if help == true {
		fmt.Println(usage)
		flag.PrintDefaults()
		os.Exit(0)
	}

	if version == true {
		fmt.Printf("pdify %s\n", VERSION)
		os.Exit(0)
	}

	if len(input) == 0 {
		logger.Fatal("missing folder operand")
	}

	g := generator.New(input)

	err := g.Generate()
	if err != nil {
		logger.Fatal(err)
	}

	err = g.Write(output)
	if err != nil {
		logger.Fatal(err)
	}

	os.Exit(0)
}
