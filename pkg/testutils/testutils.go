package testutils

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
)

const (
	envRunIntegrationTests = "RUN_INTEGRATION"
)

func SkipUnlessIntegrationTestsEnabled(t *testing.T) {
	t.Helper()

	flagValue := os.Getenv(envRunIntegrationTests)
	truthyValues := []string{"1", "true", "yes", "on"}

	if !slices.Contains(truthyValues, strings.ToLower(flagValue)) {
		message := fmt.Sprintf("Skipping test %s since %s flag is not active", t.Name(), envRunIntegrationTests)
		t.Skip(message)
	}
}
