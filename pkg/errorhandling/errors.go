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
	outputMessage.WriteString(fmt.Sprintf("\nFailed with exit code %d\n", e.ExitCode))
	if e.HelpText != "" {
		outputMessage.WriteString(fmt.Sprintf("\nAdditional details: %s", e.HelpText))
	}

	return outputMessage.String()
}
