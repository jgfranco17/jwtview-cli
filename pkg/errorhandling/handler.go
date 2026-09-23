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

	var toolErr ToolError
	if errors.Is(err, &toolErr) {
		_, _ = w.Write([]byte(toolErr.String()))
		os.Exit(toolErr.ExitCode)
	}
	os.Exit(1)
}
