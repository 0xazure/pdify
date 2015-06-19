package main

import (
  "testing"
)

type OutputPathTest struct {
  input        string
  output       string
  pwd          string
  expectedPath string
}

var outputPathTests = []OutputPathTest{
  {"album", "", "/User/user/images", "/User/user/images/album.pdf"},
  {"album", "my_album.pdf", "/User/user/images", "/User/user/images/my_album.pdf"},
  {"/User/user/images/album", "my_album.pdf", "/User/user/images", "/User/user/images/my_album.pdf"},
  {"/User/user/images/album", "/User/user/documents/my_album.pdf", "/User/user/images", "/User/user/documents/my_album.pdf"},
}

func TestOutputPath(t *testing.T) {
  for _, tt := range outputPathTests {
    input := tt.input
    output := tt.output
    pwd := tt.pwd
    expected := tt.expectedPath

    outputPath, err := OutputPath(input, output, pwd)

    if err != nil {
      t.Error(err)
    }

    if outputPath != expected {
      t.Errorf("Expected path %s, got %s\n", expected, outputPath)
    }
  }
}
