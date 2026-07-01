package signer

import (
	"context"
	"strings"
	"testing"

	"golang.org/x/oauth2"

	"github.com/apstndb/adcplus"
)

func TestSmartSigner_withTokenSourceAndTargetPrincipal(t *testing.T) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "override-token"})
	stubLookupTokeninfoEmail(t, func(context.Context, oauth2.TokenSource) (string, error) {
		t.Fatal("tokeninfo should not be called when target principal is set")
		return "", nil
	})

	s, err := SmartSigner(ctx,
		adcplus.WithTokenSource(ts),
		adcplus.WithTargetPrincipal("target@project.iam.gserviceaccount.com"),
	)
	if err != nil {
		t.Fatalf("SmartSigner() error = %v", err)
	}
	if got := s.ServiceAccount(ctx); got != "target@project.iam.gserviceaccount.com" {
		t.Fatalf("ServiceAccount() = %q, want %q", got, "target@project.iam.gserviceaccount.com")
	}
}

func TestSmartSigner_withTokenSourceInfersServiceAccountEmail(t *testing.T) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "override-token"})
	stubLookupTokeninfoEmail(t, func(context.Context, oauth2.TokenSource) (string, error) {
		return "inferred@project.iam.gserviceaccount.com", nil
	})

	s, err := SmartSigner(ctx, adcplus.WithTokenSource(ts))
	if err != nil {
		t.Fatalf("SmartSigner() error = %v", err)
	}
	if got := s.ServiceAccount(ctx); got != "inferred@project.iam.gserviceaccount.com" {
		t.Fatalf("ServiceAccount() = %q, want %q", got, "inferred@project.iam.gserviceaccount.com")
	}
}

func TestSmartSigner_withTokenSourceRejectsUserEmail(t *testing.T) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "override-token"})
	stubLookupTokeninfoEmail(t, func(context.Context, oauth2.TokenSource) (string, error) {
		return "user@example.com", nil
	})

	_, err := SmartSigner(ctx, adcplus.WithTokenSource(ts))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "non-service-account email") {
		t.Fatalf("error = %q, want non-service-account email", err.Error())
	}
	if !strings.Contains(err.Error(), "adcplus.WithTargetPrincipal") {
		t.Fatalf("error = %q, want WithTargetPrincipal guidance", err.Error())
	}
}

func TestSmartSigner_withTokenSourceRejectsEmptyEmail(t *testing.T) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "override-token"})
	stubLookupTokeninfoEmail(t, func(context.Context, oauth2.TokenSource) (string, error) {
		return "", nil
	})

	_, err := SmartSigner(ctx, adcplus.WithTokenSource(ts))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "can't infer email from the provided TokenSource") {
		t.Fatalf("error = %q, want tokeninfo email inference failure", err.Error())
	}
}

func stubLookupTokeninfoEmail(t *testing.T, lookup func(context.Context, oauth2.TokenSource) (string, error)) {
	t.Helper()
	old := lookupTokeninfoEmail
	lookupTokeninfoEmail = lookup
	t.Cleanup(func() {
		lookupTokeninfoEmail = old
	})
}
