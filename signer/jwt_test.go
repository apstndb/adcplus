package signer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type mockSigner struct {
	serviceAccount string
	keyID          string
	signature      []byte
	signBlobErr    error
}

func (m *mockSigner) ServiceAccount(ctx context.Context) string {
	return m.serviceAccount
}

func (m *mockSigner) SignBlob(ctx context.Context, b []byte) (string, []byte, error) {
	if m.signBlobErr != nil {
		return "", nil, m.signBlobErr
	}
	return m.keyID, m.signature, nil
}

func (m *mockSigner) SignJwt(ctx context.Context, claims string) (string, error) {
	return "", errors.New("not implemented")
}

func TestSignJwtHelper(t *testing.T) {
	t.Run("produces valid JWT structure", func(t *testing.T) {
		s := &mockSigner{
			keyID:     "key-123",
			signature: []byte("test-signature"),
		}

		key, signed, err := signJwtHelper(context.Background(), `{"sub":"test"}`, "kid-456", s)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if key != "key-123" {
			t.Errorf("key = %q, want %q", key, "key-123")
		}

		parts := strings.Split(signed, ".")
		if len(parts) != 3 {
			t.Fatalf("JWT has %d parts, want 3", len(parts))
		}

		// Verify header
		headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
		if err != nil {
			t.Fatalf("failed to decode header: %v", err)
		}
		var header struct {
			Alg string `json:"alg"`
			Kid string `json:"kid"`
			Typ string `json:"typ"`
		}
		if err := json.Unmarshal(headerJSON, &header); err != nil {
			t.Fatalf("failed to unmarshal header: %v", err)
		}
		if header.Alg != "RS256" {
			t.Errorf("alg = %q, want %q", header.Alg, "RS256")
		}
		if header.Kid != "kid-456" {
			t.Errorf("kid = %q, want %q", header.Kid, "kid-456")
		}
		if header.Typ != "JWT" {
			t.Errorf("typ = %q, want %q", header.Typ, "JWT")
		}

		// Verify claims
		claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			t.Fatalf("failed to decode claims: %v", err)
		}
		if string(claimsJSON) != `{"sub":"test"}` {
			t.Errorf("claims = %q, want %q", string(claimsJSON), `{"sub":"test"}`)
		}

		// Verify signature
		sig, err := base64.RawURLEncoding.DecodeString(parts[2])
		if err != nil {
			t.Fatalf("failed to decode signature: %v", err)
		}
		if string(sig) != "test-signature" {
			t.Errorf("signature = %q, want %q", string(sig), "test-signature")
		}
	})

	t.Run("SignBlob error propagates", func(t *testing.T) {
		s := &mockSigner{
			signBlobErr: errors.New("sign failed"),
		}

		_, _, err := signJwtHelper(context.Background(), `{"sub":"test"}`, "kid", s)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "sign failed" {
			t.Errorf("error = %q, want %q", err.Error(), "sign failed")
		}
	})
}
