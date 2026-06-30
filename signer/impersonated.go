package signer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/apstndb/adcplus/internal"
)

func signerFromImpersonatedServiceAccountJSON(ctx context.Context, credentialJSON []byte) (Signer, error) {
	var parsed struct {
		Type                           string          `json:"type"`
		ServiceAccountImpersonationURL string          `json:"service_account_impersonation_url"`
		Delegates                      []string        `json:"delegates"`
		SourceCredentials              json.RawMessage `json:"source_credentials"`
	}
	if err := json.Unmarshal(credentialJSON, &parsed); err != nil {
		return nil, err
	}
	if parsed.Type != internal.ImpersonatedServiceAccountKey {
		return nil, fmt.Errorf("signerFromImpersonatedServiceAccountJSON: unexpected credential type %q", parsed.Type)
	}
	if parsed.ServiceAccountImpersonationURL == "" || len(parsed.SourceCredentials) == 0 {
		return nil, errors.New("impersonated_service_account credentials require source_credentials and service_account_impersonation_url")
	}

	target, err := serviceAccountEmailFromImpersonationURL(parsed.ServiceAccountImpersonationURL)
	if err != nil {
		return nil, err
	}

	sourceCred, err := internal.GoogleCredentialsFromJSON(ctx, parsed.SourceCredentials, iamScope, userinfoEmailScope)
	if err != nil {
		return nil, err
	}
	return newIamCredentialsSigner(target, parsed.Delegates, sourceCred.TokenSource)
}

func serviceAccountEmailFromImpersonationURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid service_account_impersonation_url: %w", err)
	}

	const marker = "/serviceAccounts/"
	idx := strings.Index(u.Path, marker)
	if idx < 0 {
		return "", fmt.Errorf("invalid service_account_impersonation_url: missing serviceAccounts segment in %q", rawURL)
	}

	rest := u.Path[idx+len(marker):]
	email, _, ok := strings.Cut(rest, ":")
	if !ok || email == "" {
		return "", fmt.Errorf("invalid service_account_impersonation_url: missing service account email in %q", rawURL)
	}
	return email, nil
}
