package internal

import (
	"context"
	"testing"
)

func TestGoogleCredentialsFromJSON_authorizedUser(t *testing.T) {
	cred, err := GoogleCredentialsFromJSON(
		context.Background(),
		[]byte(`{
			"type":"authorized_user",
			"client_id":"id",
			"client_secret":"secret",
			"refresh_token":"token"
		}`),
		"https://www.googleapis.com/auth/cloud-platform",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.TokenSource == nil {
		t.Fatal("expected token source")
	}
}

func TestGoogleCredentialsFromJSON_gdchServiceAccount(t *testing.T) {
	cred, err := GoogleCredentialsFromJSON(
		context.Background(),
		[]byte(`{"type":"gdch_service_account"}`),
		"https://www.googleapis.com/auth/cloud-platform",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.TokenSource == nil {
		t.Fatal("expected token source")
	}
}
