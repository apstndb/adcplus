package signer

import (
	"context"
	"strings"
	"testing"

	"github.com/apstndb/adcplus"
)

func TestServiceAccountEmailFromImpersonationURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		want    string
		wantErr bool
	}{
		{
			name:   "generateAccessToken URL",
			rawURL: "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/target@p.iam.gserviceaccount.com:generateAccessToken",
			want:   "target@p.iam.gserviceaccount.com",
		},
		{
			name:   "generateIdToken URL",
			rawURL: "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/target@p.iam.gserviceaccount.com:generateIdToken",
			want:   "target@p.iam.gserviceaccount.com",
		},
		{
			name:    "missing service account segment",
			rawURL:  "https://example.com/v1/projects/-/foo/bar",
			wantErr: true,
		},
		{
			name:    "missing email suffix",
			rawURL:  "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := serviceAccountEmailFromImpersonationURL(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("serviceAccountEmailFromImpersonationURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("serviceAccountEmailFromImpersonationURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSmartSigner_unsupportedCredentialTypes(t *testing.T) {
	ctx := context.Background()

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
			_, err := SmartSigner(ctx, adcplus.WithCredentialsJSON([]byte(tt.credential)))
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantSubstr) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantSubstr)
			}
		})
	}
}

func TestSignerFromImpersonatedServiceAccountJSON_missingFields(t *testing.T) {
	_, err := signerFromImpersonatedServiceAccountJSON(context.Background(), []byte(`{"type":"impersonated_service_account"}`))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
