package parsing

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	token := fixtureToken(t)

	decoded, err := Read(token)

	require.NoError(t, err)
	assert.Equal(t, "HS256", decoded.Header["alg"])
	assert.Equal(t, "JWT", decoded.Header["typ"])
	assert.Equal(t, "user-123", decoded.Payload["sub"])
	assert.Equal(t, float64(3600), decoded.Payload["exp"])
	assert.Equal(t, "c2lnbmF0dXJl", decoded.Signature)
}

func TestReadErrors(t *testing.T) {
	tests := []struct {
		name  string
		token string
		err   string
	}{
		{
			name: "empty token",
			err:  "raw token cannot be empty",
		},
		{
			name:  "malformed token",
			token: "not-a-jwt",
			err:   "parse failure",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decoded, err := Read(test.token)

			assert.Nil(t, decoded)
			assert.ErrorContains(t, err, test.err)
		})
	}
}

func TestDecodedTokenToJSON(t *testing.T) {
	decoded := &DecodedToken{
		Header: map[string]any{"alg": "none"},
		Payload: jwt.MapClaims{
			"sub": "user-123",
		},
		Signature: "signature",
	}

	output, err := decoded.ToJSON()

	require.NoError(t, err)
	var value map[string]any
	require.NoError(t, json.Unmarshal([]byte(output), &value))
	assert.Equal(t, map[string]any{"alg": "none"}, value["header"])
	assert.Equal(t, map[string]any{"sub": "user-123"}, value["payload"])
	assert.Equal(t, "signature", value["signature"])
}

func fixtureToken(t *testing.T) string {
	t.Helper()

	parts := make([]string, 3)
	copy(parts, []string{
		base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`)),
		base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"user-123","exp":3600}`)),
		base64.RawURLEncoding.EncodeToString([]byte("signature")),
	})
	return parts[0] + "." + parts[1] + "." + parts[2]
}
