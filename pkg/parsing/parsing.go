// Package parsing provides utilities for inspecting JWT tokens without
// verifying their signature. Intended for debugging/dev tooling use only.
package parsing

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jgfranco17/jwtview-cli/pkg/errorhandling"
)

// DecodedToken holds the human-readable parts of a JWT.
type DecodedToken struct {
	Header    map[string]any `json:"header"`
	Payload   jwt.MapClaims  `json:"payload"`
	Signature string         `json:"signature"` // base64url, unverified
}

// ToJSON returns the decoded token as JSON, pretty-printed or compact.
func (d *DecodedToken) ToJSON(pretty bool) (string, error) {
	var out []byte
	var err error
	if pretty {
		out, err = json.MarshalIndent(d, "", "  ")
	} else {
		out, err = json.Marshal(d)
	}
	if err != nil {
		return "", fmt.Errorf("marshal decoded token: %w", err)
	}

	return string(out), nil
}

// Read decodes a JWT's header and claims into a DecodedToken.
func Read(token string) (*DecodedToken, error) {
	if token == "" {
		return nil, &errorhandling.ToolError{
			BaseError: errors.New("raw token cannot be empty"),
		}
	}

	parser := jwt.NewParser()
	claims := jwt.MapClaims{}
	parsed, _, err := parser.ParseUnverified(token, claims)
	if err != nil {
		return nil, &errorhandling.ToolError{
			BaseError: fmt.Errorf("parse failure: %w", err),
		}
	}

	return &DecodedToken{
		Header:    parsed.Header,
		Payload:   claims,
		Signature: base64.RawURLEncoding.EncodeToString(parsed.Signature),
	}, nil
}
