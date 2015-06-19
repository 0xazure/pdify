package main

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/alecthomas/kingpin"
)

const VERSION = "0.1.0"

var (
  input = kingpin.Arg("input", "Source folder.").Required().String()

  output = kingpin.Flag("output", "Output file name.").Short('o').String()
)

func Run() int {
  return run()
}

func main() {
  os.Exit(Run())
}

// Internal

func run() int {
  kingpin.Version(VERSION)
  kingpin.Parse()

  inputPath, err := filepath.Abs(*input)
  if err != nil {
    fmt.Printf("Error: Unable to read %s: %s", *input, err)
    return 1
  }

  pwd, _ := os.Getwd()
  outputPath, err := OutputPath(*input, *output, pwd)
  if err != nil {
    fmt.Printf("Error: Invalid output %s: %s", *output, err)
    return 1
  }

  files := Analyze(inputPath)

  if pdfErr := Generate(files, outputPath); err != nil {
    fmt.Printf("Error: unable to generate PDF, %s", pdfErr)
    return 1
  }

  return 0
}
