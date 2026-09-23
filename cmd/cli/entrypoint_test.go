// Package cli_test exercises the wired-together root command (flags,
// parsing, and logging) via CommandRoot directly, without a subprocess.
package cli_test

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strings"
	"testing"

	"github.com/jgfranco17/jwtview-cli/cmd/cli"
	"github.com/jgfranco17/jwtview-cli/pkg/errorhandling"
	"github.com/jgfranco17/jwtview-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type cliResult struct {
	stdout   string
	stderr   string
	exitCode int
}

func TestCLIHelp(t *testing.T) {
	testutils.SkipUnlessIntegrationTestsEnabled(t)

	result := runCLI(t, "--help")

	assert.Equal(t, 0, result.exitCode)
	assert.Contains(t, result.stdout, "JwtView is a command-line tool")
}

func TestCLIMissingToken(t *testing.T) {
	testutils.SkipUnlessIntegrationTestsEnabled(t)

	result := runCLI(t)

	assert.Equal(t, 1, result.exitCode)
	assert.Contains(t, result.stderr, "accepts 1 arg(s), received 0")
}

func TestCLIInvalidToken(t *testing.T) {
	testutils.SkipUnlessIntegrationTestsEnabled(t)

	result := runCLI(t, "not-a-jwt")

	assert.Equal(t, 1, result.exitCode)
	assert.Contains(t, result.stderr, "parse failure")
}

func TestCLIMutuallyExclusiveFlags(t *testing.T) {
	testutils.SkipUnlessIntegrationTestsEnabled(t)

	token := fixtureToken(`{"sub":"user-123"}`)
	result := runCLI(t, "--pretty", "--compact", token)

	assert.Equal(t, 1, result.exitCode)
	assert.Contains(t, result.stderr, "none of the others can be")
}

func TestCLIDecodeToken(t *testing.T) {
	testutils.SkipUnlessIntegrationTestsEnabled(t)

	token := fixtureToken(`{"sub":"user-123"}`)

	tests := []struct {
		name      string
		args      []string
		multiline bool
	}{
		{name: "default", args: []string{token}, multiline: true},
		{name: "pretty", args: []string{"--pretty", token}, multiline: true},
		{name: "compact", args: []string{"--compact", token}, multiline: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := runCLI(t, test.args...)

			require.Equal(t, 0, result.exitCode)
			assert.Contains(t, result.stdout, `"sub"`)
			assert.Contains(t, result.stdout, "user-123")
			body := strings.TrimSpace(result.stdout)
			assert.Equal(t, test.multiline, strings.Contains(body, "\n"))
		})
	}
}

func TestCLIVerboseLogging(t *testing.T) {
	testutils.SkipUnlessIntegrationTestsEnabled(t)

	token := fixtureToken(`{"sub":"user-123"}`)
	result := runCLI(t, "-vv", token)

	require.Equal(t, 0, result.exitCode)
	assert.Contains(t, result.stderr, "Decoded JWT successfully")
}

func fixtureToken(payload string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	claims := base64.RawURLEncoding.EncodeToString([]byte(payload))
	signature := base64.RawURLEncoding.EncodeToString([]byte("signature"))
	return header + "." + claims + "." + signature
}

func runCLI(t *testing.T, args ...string) cliResult {
	t.Helper()
	var stdout, stderr bytes.Buffer

	root := cli.NewCommandRoot(cli.Options{Version: "test"})
	root.SetArgs(args)
	root.SetOutput(&stdout, &stderr)

	err := root.Execute(t.Context())

	exitCode := 0
	if err != nil {
		exitCode = errorhandling.ExitCodeGeneralFailure
		var toolErr *errorhandling.ToolError
		if errors.As(err, &toolErr) {
			exitCode = toolErr.ResolvedExitCode()
		}
	}

	return cliResult{
		stdout:   stdout.String(),
		stderr:   stderr.String(),
		exitCode: exitCode,
	}
}
