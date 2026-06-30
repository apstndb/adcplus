package signer

import (
	"context"
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

func TestSmartSigner_withTokenSourceAndTargetPrincipal(t *testing.T) {
	ctx := context.Background()
	ts := staticAccessTokenSource{token: &oauth2.Token{AccessToken: "override-token"}}

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
