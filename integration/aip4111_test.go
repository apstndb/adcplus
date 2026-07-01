//go:build integration

package integration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/apstndb/adcplus"
	"github.com/apstndb/adcplus/tokensource"
)

const (
	aip4111ServiceAccountKeyJSONEnv = "ADCPLUS_AIP4111_SERVICE_ACCOUNT_KEY_JSON"
	aip4111ServiceAccountKeyFileEnv = "ADCPLUS_AIP4111_SERVICE_ACCOUNT_KEY_FILE"
	aip4111ProjectIDEnv             = "ADCPLUS_AIP4111_PROJECT_ID"
	aip4111DefaultProjectID         = "apstndb-sandbox"
	aip4111CloudPlatformScope       = "https://www.googleapis.com/auth/cloud-platform"
)

func TestSmartAccessTokenSource_AIP4111JWTAccessWithScope_serviceUsage(t *testing.T) {
	keyJSON := aip4111ServiceAccountKeyJSON(t)
	if len(keyJSON) == 0 {
		t.Skipf("%s or %s not set", aip4111ServiceAccountKeyJSONEnv, aip4111ServiceAccountKeyFileEnv)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	ts, err := tokensource.SmartAccessTokenSource(
		ctx,
		adcplus.WithCredentialsJSON(keyJSON),
		adcplus.WithScopes(aip4111CloudPlatformScope),
		adcplus.WithJWTAccessWithScope(true),
	)
	if err != nil {
		t.Fatalf("SmartAccessTokenSource() error = %v", err)
	}

	tok, err := ts.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if tok.AccessToken == "" {
		t.Fatal("expected non-empty access token")
	}
	if tok.TokenType != "Bearer" {
		t.Fatalf("TokenType = %q, want Bearer", tok.TokenType)
	}

	claims := decodeJWTClaimsForIntegration(t, tok.AccessToken)
	if got := claims["scope"]; got != aip4111CloudPlatformScope {
		t.Fatalf("scope claim = %v, want %q", got, aip4111CloudPlatformScope)
	}

	projectID := os.Getenv(aip4111ProjectIDEnv)
	if projectID == "" {
		projectID = aip4111DefaultProjectID
	}
	url := fmt.Sprintf("https://serviceusage.googleapis.com/v1/projects/%s/services?pageSize=1", projectID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext() error = %v", err)
	}
	tok.SetAuthHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Service Usage request error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		t.Fatalf("Service Usage status = %d; reading error body: %v", resp.StatusCode, err)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		t.Fatalf("Service Usage returned 401; AIP-4111 self-signed JWT access token was not accepted: %s", string(body))
	case http.StatusForbidden:
		t.Fatalf("Service Usage returned 403; check roles/serviceusage.serviceUsageViewer on project %q: %s", projectID, string(body))
	default:
		t.Fatalf("Service Usage status = %d, want 200: %s", resp.StatusCode, string(body))
	}
}

func aip4111ServiceAccountKeyJSON(t *testing.T) []byte {
	t.Helper()

	if keyJSON := os.Getenv(aip4111ServiceAccountKeyJSONEnv); keyJSON != "" {
		return []byte(keyJSON)
	}

	keyFile := os.Getenv(aip4111ServiceAccountKeyFileEnv)
	if keyFile == "" {
		return nil
	}
	keyJSON, err := os.ReadFile(keyFile)
	if err != nil {
		t.Fatalf("ReadFile(%s) error = %v", keyFile, err)
	}
	if len(keyJSON) == 0 {
		t.Fatalf("%s points to an empty file", aip4111ServiceAccountKeyFileEnv)
	}
	return keyJSON
}

func decodeJWTClaimsForIntegration(t *testing.T, token string) map[string]any {
	t.Helper()

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("access token has %d JWT segments, want 3", len(parts))
	}
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("DecodeString() error = %v", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	return claims
}
