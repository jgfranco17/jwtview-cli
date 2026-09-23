package errorhandling

import (
	"fmt"
	"strings"
)

const (
	ExitCodeSuccess         int = 0
	ExitCodeGeneralFailure  int = 1
	ExitCodeInvalidInput    int = 2
	ExitCodeValidationError int = 3
)

// ToolError is a custom error type that includes an exit code for CLI applications.
type ToolError struct {
	BaseError error
	ExitCode  int
	HelpText  string
}

func (e *ToolError) Error() string {
	return e.BaseError.Error()
}

func (e *ToolError) String() string {
	outputMessage := strings.Builder{}

	outputMessage.WriteString("[ERROR] Tool encountered an error:\n")
	outputMessage.WriteString(fmt.Sprintf("%s\n", e.BaseError.Error()))
	outputMessage.WriteString(fmt.Sprintf("\nFailed with exit code %d\n", e.ResolvedExitCode()))
	if e.HelpText != "" {
		outputMessage.WriteString(fmt.Sprintf("\nAdditional details: %s", e.HelpText))
	}

	return outputMessage.String()
}

// ResolvedExitCode returns the error's exit code, defaulting to a general
// failure since an unset ExitCode must never be reported as success.
func (e *ToolError) ResolvedExitCode() int {
	if e.ExitCode == ExitCodeSuccess {
		return ExitCodeGeneralFailure
	}
	return e.ExitCode
}
