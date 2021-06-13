package tokensource

import (
	"context"

	"github.com/apstndb/adcplus"
	"github.com/apstndb/adcplus/internal"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/impersonate"
)

const cloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"

// SmartIDTokenSource generate oauth2.TokenSource which generates ID token and supports CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT environment variable.
func SmartIDTokenSource(ctx context.Context, audience string, options ...adcplus.Option) (oauth2.TokenSource, error) {
	config, err := internal.CalcAdcPlusConfig(options...)
	if err != nil {
		return nil, err
	}
	if config.TargetPrincipal != "" {
		idCfg := impersonate.IDTokenConfig{
			Audience:        audience,
			TargetPrincipal: config.TargetPrincipal,
			Delegates:       config.Delegates,
			// Cloud IAP requires email claim.
			IncludeEmail: true,
		}
		return impersonate.IDTokenSource(ctx, idCfg)
	}

	return idtoken.NewTokenSource(ctx, audience)
}

// SmartAccessTokenSource generate oauth2.TokenSource which generates access token and supports CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT environment variable.
func SmartAccessTokenSource(ctx context.Context, options ...adcplus.Option) (oauth2.TokenSource, error) {
	config, err := internal.CalcAdcPlusConfig(options...)
	if err != nil {
		return nil, err
	}
	if len(config.Scopes) == 0 {
		config.Scopes = []string{cloudPlatformScope}
	}
	if config.TargetPrincipal != "" {
		return impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: config.TargetPrincipal,
			Delegates:       config.Delegates,
			Scopes:          config.Scopes,
		})
	}
	return google.DefaultTokenSource(ctx, config.Scopes...)
}
