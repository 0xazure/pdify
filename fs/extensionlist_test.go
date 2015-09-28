package fs

import (
	"testing"
)

type ExtensionListTest struct {
	exts     []string
	expected []string
}

var extListTests = []ExtensionListTest{
	{
		[]string{"png", "jpg", "jpeg"},
		[]string{".png", ".jpg", ".jpeg"},
	},
	{
		[]string{".PNG", ".png", ".JPG"},
		[]string{".png", ".jpg"},
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

func TestNewExtensionList(t *testing.T) {
	for _, tt := range extListTests {
		exts := tt.exts
		expected := tt.expected

		list := NewExtensionList(exts)

		// Check that the list has the same number of elements as is expected
		expectedLen := len(expected)
		listLen := len(list.Exts())
		if expectedLen != listLen {
			t.Errorf("Expected map length of %d, got %d", expectedLen, listLen)
		}
	}
}

func TestExtensionList_Contains(t *testing.T) {
	for _, tt := range extListTests {
		exts := tt.exts
		expected := tt.expected

		list := NewExtensionList(exts)

		// Check that all expected values are comtained in the list
		for _, ext := range expected {
			if contained := list.Contains(ext); !contained {
				t.Errorf("'%s' is not contained in the list", ext)
			}
		}
	}
}
