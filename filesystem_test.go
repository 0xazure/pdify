package main

import (
	"testing"
)

type ExtensionMapTest struct {
	extensions        []string
	expectedMapValues []string
}

var extensionMapTests = []ExtensionMapTest{
	{
		[]string{"png", "jpg", "jpeg"},
		[]string{".png", ".jpg", ".jpeg"},
	},
	{
		[]string{".png", "jpg", ".jpeg"},
		[]string{".png", ".jpg", ".jpeg"},
	},
	{
		[]string{"", "   ", ".png"},
		[]string{".png"},
	},
	{
		[]string{"old.png", "my.jpg"},
		[]string{".old.png", ".my.jpg"},
	},
}

func TestExtensionMapBuilder(t *testing.T) {
	for _, tt := range extensionMapTests {
		extensions := tt.extensions
		expectedMapValues := tt.expectedMapValues
		extMap := ExtensionMapBuilder(extensions)

		// Check that all expected values are contained in the generated map
		for _, ext := range expectedMapValues {
			if _, ok := extMap[ext]; ok != true {
				t.Errorf("%s not contained in map.\n", ext)
			}
		}

		// Check that the map has the same number of elements as is expected
		expectedLen := len(expectedMapValues)
		extMapLen := len(extMap)
		if extMapLen != expectedLen {
			t.Errorf("Expected map length of %v, got %v\n", expectedLen, extMapLen)
		}
	}
}

type IsValidExtensionTest struct {
	ext string
	r   bool
}

var isValidExtensionTests = []IsValidExtensionTest{
	{"pg00.png", true},
	{"pg01.jpg", true},
	{"pg02.JPG", true},
	{"pg03.jpeg", true},
	{"pg04.ico", false},
	{"pg05.ICO", false},
	{"pg06.txt", false},
}

var extensionMap = map[string]struct{}{
	".png":  {},
	".jpg":  {},
	".jpeg": {},
}

func TestIsValidExtension(t *testing.T) {
	for _, tt := range isValidExtensionTests {
		extension := tt.ext
		result := tt.r

		if ok := isValidExtension(extension, extensionMap); ok != result {
			t.Errorf("%s gave result %t, expected %t\n", extension, ok, result)
		}
	}
}
