//go:build integration

package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/apstndb/adcplus"
	"github.com/apstndb/adcplus/tokensource"
)

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
