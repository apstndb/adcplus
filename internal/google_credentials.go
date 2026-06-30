package internal

import (
	"context"

	"golang.org/x/oauth2/google"
)

// GoogleCredentialsFromJSON loads credentials using google.CredentialsFromJSONWithType
// so the JSON type must match the declared credential type field.
func GoogleCredentialsFromJSON(ctx context.Context, credentialJSON []byte, scopes ...string) (*google.Credentials, error) {
	credType, err := CredentialTypeFromJSON(credentialJSON)
	if err != nil {
		return nil, err
	}
	return google.CredentialsFromJSONWithType(ctx, credentialJSON, google.CredentialsType(credType), scopes...)
}
