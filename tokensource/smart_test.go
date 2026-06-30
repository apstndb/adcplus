package tokensource

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
