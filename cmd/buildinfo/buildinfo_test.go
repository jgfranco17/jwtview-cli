package buildinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	tcases := []struct {
		name            string
		definition      []byte
		expectedError   string
		expectedVersion string
	}{{
		name:            "empty definition",
		expectedError:   "unexpected end of JSON input",
		expectedVersion: "0.0.0",
	}, {
		name:            "invalid json definition",
		definition:      []byte("invalid"),
		expectedError:   "invalid character",
		expectedVersion: "0.0.0",
	}, {
		name:            "version missing in definition",
		definition:      []byte("{}"),
		expectedVersion: "0.0.0",
	}, {
		name:            "all good",
		definition:      []byte(`{"version":"1.2.3"}`),
		expectedVersion: "1.2.3",
	}}

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := GetVersion(tc.definition)
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, v, tc.expectedVersion)
		})
	}
}
