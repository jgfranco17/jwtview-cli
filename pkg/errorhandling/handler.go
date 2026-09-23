package errorhandling

import (
	"errors"
	"io"
	"os"
)

func Exit(w io.Writer, err error) {
	if err == nil {
		return
	}

	var toolErr *ToolError
	if errors.As(err, &toolErr) {
		_, _ = w.Write([]byte(toolErr.String()))
		os.Exit(toolErr.ResolvedExitCode())
	}
	os.Exit(1)
}
