//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/apstndb/adcplus"
	"github.com/apstndb/adcplus/signer"
	"github.com/apstndb/adcplus/tokensource"
)

// wifIDTokenAudience is the audience for SmartIDTokenSource impersonation tests.
const wifIDTokenAudience = "https://adcplus-ci.test"

func TestSmartAccessTokenSource_WIF(t *testing.T) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("GOOGLE_APPLICATION_CREDENTIALS not set; run under GitHub Actions WIF or local ADC")
	}

	ts, err := tokensource.SmartAccessTokenSource(context.Background())
	if err != nil {
		t.Fatalf("SmartAccessTokenSource: %v", err)
	}

	tok, err := ts.Token()
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	if tok.AccessToken == "" {
		t.Fatal("expected non-empty access token")
	}
}

func TestSmartAccessTokenSource_WIF_impersonation(t *testing.T) {
	target := os.Getenv("ADCPLUS_IMPERSONATION_TARGET")
	if target == "" {
		t.Skip("ADCPLUS_IMPERSONATION_TARGET not set")
	}

	ts, err := tokensource.SmartAccessTokenSource(
		context.Background(),
		adcplus.WithTargetPrincipal(target),
	)
	if err != nil {
		t.Fatalf("SmartAccessTokenSource: %v", err)
	}

	tok, err := ts.Token()
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	if tok.AccessToken == "" {
		t.Fatal("expected non-empty access token")
	}
}

func TestSmartSigner_WIF(t *testing.T) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("GOOGLE_APPLICATION_CREDENTIALS not set; run under GitHub Actions WIF or local ADC")
	}

	ctx := context.Background()
	s, err := signer.SmartSigner(ctx)
	if err != nil {
		t.Fatalf("SmartSigner: %v", err)
	}

	email := s.ServiceAccount(ctx)
	if email == "" {
		t.Fatal("expected non-empty service account email")
	}

	keyID, sig, err := s.SignBlob(ctx, []byte("adcplus-wif-integration"))
	if err != nil {
		t.Fatalf("SignBlob: %v", err)
	}
	if keyID == "" {
		t.Fatal("expected non-empty key ID from SignBlob")
	}
	if len(sig) == 0 {
		t.Fatal("expected non-empty signature from SignBlob")
	}

	now := time.Now()
	payload := fmt.Sprintf(
		`{"iss":%q,"sub":%q,"aud":%q,"iat":%d,"exp":%d}`,
		email, email, wifIDTokenAudience, now.Unix(), now.Add(time.Hour).Unix(),
	)
	jwt, err := s.SignJwt(ctx, payload)
	if err != nil {
		t.Fatalf("SignJwt: %v", err)
	}
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		t.Fatalf("SignJwt: expected 3 JWT segments, got %d", len(parts))
	}
	for i, part := range parts {
		if part == "" {
			t.Fatalf("SignJwt: segment %d is empty", i)
		}
	}
}

func TestSmartIDTokenSource_WIF_impersonation(t *testing.T) {
	target := os.Getenv("ADCPLUS_IMPERSONATION_TARGET")
	if target == "" {
		t.Skip("ADCPLUS_IMPERSONATION_TARGET not set")
	}

	ts, err := tokensource.SmartIDTokenSource(
		context.Background(),
		wifIDTokenAudience,
		adcplus.WithTargetPrincipal(target),
	)
	if err != nil {
		t.Fatalf("SmartIDTokenSource: %v", err)
	}

	tok, err := ts.Token()
	if err != nil {
		t.Fatalf("Token: %v", err)
	}
	if tok.AccessToken == "" {
		t.Fatal("expected non-empty ID token")
	}
}
