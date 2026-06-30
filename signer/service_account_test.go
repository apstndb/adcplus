package signer

import (
	"strings"
	"testing"
)

func TestNewServiceAccountSigner_invalidPEM(t *testing.T) {
	_, err := newServiceAccountSigner([]byte(`{
		"type":"service_account",
		"client_email":"sa@test.iam.gserviceaccount.com",
		"private_key_id":"id",
		"private_key":"not-a-pem"
	}`))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "PEM") {
		t.Fatalf("error = %q, want PEM decode failure", err.Error())
	}
}

func TestNewServiceAccountSigner_unsupportedType(t *testing.T) {
	_, err := newServiceAccountSigner([]byte(`{
		"type":"authorized_user",
		"client_email":"sa@test.iam.gserviceaccount.com",
		"private_key":"not-a-pem"
	}`))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported local signing credential type") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestNewServiceAccountSigner_gdchType(t *testing.T) {
	_, err := newServiceAccountSigner([]byte(`{
		"type":"gdch_service_account",
		"client_email":"sa@test.iam.gserviceaccount.com",
		"private_key_id":"id",
		"private_key":"not-a-pem"
	}`))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "PEM") {
		t.Fatalf("error = %q, want PEM decode failure", err.Error())
	}
}
