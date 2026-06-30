package tokensource

import (
	"context"
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
