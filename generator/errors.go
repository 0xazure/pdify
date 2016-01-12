package generator

type ProcessError struct {
	Err      error
	ExitCode int
}

func (e *ProcessError) Error() string {
	return e.Err.Error()
}

func newProcessError(err error) ProcessError {
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	return ProcessError{Err: err, ExitCode: exitCode}
}
