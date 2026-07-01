package tokensource

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/oauth2"

	"github.com/apstndb/adcplus"
)

type staticAccessTokenSource struct {
	token *oauth2.Token
}

func (s staticAccessTokenSource) Token() (*oauth2.Token, error) {
	return s.token, nil
}

func testServiceAccountJSON(t *testing.T, tokenURL string) []byte {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}
	privateKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("MarshalPKCS8PrivateKey() error = %v", err)
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKey,
	})

	credential := map[string]string{
		"type":           "service_account",
		"project_id":     "test-project",
		"private_key_id": "test-key-id",
		"private_key":    string(privateKeyPEM),
		"client_email":   "test-sa@test-project.iam.gserviceaccount.com",
		"client_id":      "1234567890",
		"token_uri":      tokenURL,
	}
	j, err := json.Marshal(credential)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	return j
}

func decodeJWTClaims(t *testing.T, token string) map[string]any {
	t.Helper()

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("JWT has %d parts, want 3", len(parts))
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

func oauthTokenServer(t *testing.T, token string) (*httptest.Server, *bool) {
	t.Helper()

	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"access_token":"` + token + `","token_type":"Bearer","expires_in":3600}`)); err != nil {
			t.Errorf("Write() error = %v", err)
		}
	}))
	return server, &called
}

func TestSmartIDTokenSource_withTokenSource(t *testing.T) {
	ctx := context.Background()
	want := &oauth2.Token{AccessToken: "override-id-token"}

	ts, err := SmartIDTokenSource(ctx, "https://example.com", adcplus.WithTokenSource(staticAccessTokenSource{token: want}))
	if err != nil {
		t.Fatalf("SmartIDTokenSource() error = %v", err)
	}

	got, err := ts.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if got.AccessToken != want.AccessToken {
		t.Fatalf("AccessToken = %q, want %q", got.AccessToken, want.AccessToken)
	}
}

func TestSmartAccessTokenSource_withTokenSource(t *testing.T) {
	ctx := context.Background()
	want := &oauth2.Token{AccessToken: "override-token"}

	ts, err := SmartAccessTokenSource(ctx, adcplus.WithTokenSource(staticAccessTokenSource{token: want}))
	if err != nil {
		t.Fatalf("SmartAccessTokenSource() error = %v", err)
	}

	got, err := ts.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if got.AccessToken != want.AccessToken {
		t.Fatalf("AccessToken = %q, want %q", got.AccessToken, want.AccessToken)
	}
}

func TestSmartAccessTokenSource_jwtAccessWithScope(t *testing.T) {
	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected token endpoint request: %s %s", r.Method, r.URL.String())
	}))
	defer server.Close()

	ts, err := SmartAccessTokenSource(
		ctx,
		adcplus.WithCredentialsJSON(testServiceAccountJSON(t, server.URL)),
		adcplus.WithScopes("scope1", "scope2"),
		adcplus.WithJWTAccessWithScope(true),
	)
	if err != nil {
		t.Fatalf("SmartAccessTokenSource() error = %v", err)
	}

	tok, err := ts.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if tok.TokenType != "Bearer" {
		t.Errorf("TokenType = %q, want Bearer", tok.TokenType)
	}

	claims := decodeJWTClaims(t, tok.AccessToken)
	wantEmail := "test-sa@test-project.iam.gserviceaccount.com"
	if claims["iss"] != wantEmail {
		t.Errorf("iss = %v, want %q", claims["iss"], wantEmail)
	}
	if claims["sub"] != wantEmail {
		t.Errorf("sub = %v, want %q", claims["sub"], wantEmail)
	}
	if claims["scope"] != "scope1 scope2" {
		t.Errorf("scope = %v, want %q", claims["scope"], "scope1 scope2")
	}
	if aud, ok := claims["aud"]; ok && aud != "" {
		t.Errorf("aud = %v, want empty or absent", aud)
	}
}

func TestSmartAccessTokenSource_jwtAccessWithScope_nonServiceAccount(t *testing.T) {
	ctx := context.Background()
	credential := []byte(`{
		"type":"authorized_user",
		"client_id":"cid",
		"client_secret":"secret",
		"refresh_token":"token"
	}`)

	_, err := SmartAccessTokenSource(
		ctx,
		adcplus.WithCredentialsJSON(credential),
		adcplus.WithScopes("scope1"),
		adcplus.WithJWTAccessWithScope(true),
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "requires service_account credentials JSON") {
		t.Fatalf("error = %q, want service_account requirement", err.Error())
	}
}

func TestSmartAccessTokenSource_noJWTAccessOptInUsesOAuthTokenSource(t *testing.T) {
	ctx := context.Background()
	server, called := oauthTokenServer(t, "oauth-token")
	defer server.Close()

	ts, err := SmartAccessTokenSource(
		ctx,
		adcplus.WithCredentialsJSON(testServiceAccountJSON(t, server.URL)),
		adcplus.WithScopes("scope1"),
	)
	if err != nil {
		t.Fatalf("SmartAccessTokenSource() error = %v", err)
	}

	tok, err := ts.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if tok.AccessToken != "oauth-token" {
		t.Fatalf("AccessToken = %q, want oauth-token", tok.AccessToken)
	}
	if !*called {
		t.Fatal("expected OAuth token endpoint request")
	}
}

func TestSmartAccessTokenSource_jwtAccessWithScope_skippedRoutes(t *testing.T) {
	ctx := context.Background()

	t.Run("WithTokenSource overrides JWT access", func(t *testing.T) {
		want := &oauth2.Token{AccessToken: "override-token"}
		ts, err := SmartAccessTokenSource(
			ctx,
			adcplus.WithCredentialsJSON([]byte(`{"type":"service_account","private_key":"not-a-pem"}`)),
			adcplus.WithScopes("scope1"),
			adcplus.WithJWTAccessWithScope(true),
			adcplus.WithTokenSource(staticAccessTokenSource{token: want}),
		)
		if err != nil {
			t.Fatalf("SmartAccessTokenSource() error = %v", err)
		}
		got, err := ts.Token()
		if err != nil {
			t.Fatalf("Token() error = %v", err)
		}
		if got.AccessToken != want.AccessToken {
			t.Fatalf("AccessToken = %q, want %q", got.AccessToken, want.AccessToken)
		}
	})

	t.Run("no explicit scopes uses OAuth token source", func(t *testing.T) {
		server, called := oauthTokenServer(t, "default-scope-oauth-token")
		defer server.Close()

		ts, err := SmartAccessTokenSource(
			ctx,
			adcplus.WithCredentialsJSON(testServiceAccountJSON(t, server.URL)),
			adcplus.WithJWTAccessWithScope(true),
		)
		if err != nil {
			t.Fatalf("SmartAccessTokenSource() error = %v", err)
		}
		tok, err := ts.Token()
		if err != nil {
			t.Fatalf("Token() error = %v", err)
		}
		if tok.AccessToken != "default-scope-oauth-token" {
			t.Fatalf("AccessToken = %q, want default-scope-oauth-token", tok.AccessToken)
		}
		if !*called {
			t.Fatal("expected OAuth token endpoint request")
		}
	})

	t.Run("impersonation is checked before JWT access", func(t *testing.T) {
		_, err := SmartAccessTokenSource(
			ctx,
			adcplus.WithCredentialsJSON([]byte(`{"type":"service_account","private_key":"not-a-pem"}`)),
			adcplus.WithScopes("scope1"),
			adcplus.WithJWTAccessWithScope(true),
			adcplus.WithTargetPrincipal("target@test-project.iam.gserviceaccount.com"),
			adcplus.WithTokenSource(staticAccessTokenSource{token: &oauth2.Token{AccessToken: "source-token"}}),
		)
		if err != nil {
			t.Fatalf("SmartAccessTokenSource() error = %v", err)
		}
	})
}

func TestSmartIDTokenSource_unsupportedCredentialTypes(t *testing.T) {
	ctx := context.Background()
	audience := "https://example.com"

	tests := []struct {
		name       string
		credential string
		wantSubstr string
	}{
		{
			name: "authorized_user",
			credential: `{
				"type":"authorized_user",
				"client_id":"cid",
				"client_secret":"secret",
				"refresh_token":"token"
			}`,
			wantSubstr: "authorized_user is unsupported",
		},
		{
			name: "external_account_authorized_user",
			credential: `{
				"type":"external_account_authorized_user",
				"client_id":"cid",
				"client_secret":"secret",
				"refresh_token":"token",
				"token_url":"https://example.com/token"
			}`,
			wantSubstr: "external_account_authorized_user is unsupported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SmartIDTokenSource(ctx, audience, adcplus.WithCredentialsJSON([]byte(tt.credential)))
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantSubstr) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantSubstr)
			}
		})
	}
}

func TestSmartIDTokenSource_externalAccount_delegatesToIDToken(t *testing.T) {
	ctx := context.Background()
	audience := "https://example.com"
	credential := `{
		"type":"external_account",
		"audience":"//iam.googleapis.com/projects/123/locations/global/workloadIdentityPools/pool/providers/provider",
		"subject_token_type":"urn:ietf:params:oauth:token-type:jwt",
		"service_account_impersonation_url":"https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/target@p.iam.gserviceaccount.com:generateAccessToken",
		"token_url":"https://sts.googleapis.com/v1/token",
		"credential_source":{"file":"/tmp/nonexistent-adcplus-wif-token"}
	}`

	ts, err := SmartIDTokenSource(ctx, audience, adcplus.WithCredentialsJSON([]byte(credential)))
	if err != nil {
		t.Fatalf("SmartIDTokenSource() error = %v", err)
	}
	_, err = ts.Token()
	if err == nil {
		t.Fatal("expected token fetch error, got nil")
	}
	if strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("Token() error = %q, should delegate to idtoken not reject credential type", err.Error())
	}
}
