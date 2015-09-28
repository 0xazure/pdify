package generator

import (
	"errors"
	"testing"
)

type ProcessErrorTest struct {
	err  error
	code int
}

var processErrorTests = []ProcessErrorTest{
	{errors.New("Invalid"), 1},
	{nil, 0},
}

func TestGenerator_newProcessError(t *testing.T) {
	for _, tt := range processErrorTests {
		err := tt.err
		actualErr := newProcessError(err)

		if err != actualErr.Err {
			t.Errorf("Expected error %v, got %v", err, actualErr.Err)
		}

		if err != nil && actualErr.Err != nil {
			expectedMsg := err.Error()
			actualMsg := actualErr.Error()
			if expectedMsg != actualMsg {
				t.Errorf("Expected error message '%s', got '%s'", expectedMsg, actualMsg)
			}
		}

		expectedCode := tt.code
		actualCode := actualErr.ExitCode
		if expectedCode != actualCode {
			t.Errorf("Expected exit code %d, got %d", expectedCode, actualCode)
		}
	}
}
