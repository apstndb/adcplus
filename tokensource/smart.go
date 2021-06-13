package tokensource

import (
	"context"

	"github.com/apstndb/adcplus"
	"github.com/apstndb/adcplus/internal"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

const cloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"

// SmartIDTokenSource generate oauth2.TokenSource which generates ID token and supports CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT environment variable.
func SmartIDTokenSource(ctx context.Context, audience string, options ...adcplus.Option) (oauth2.TokenSource, error) {
	config, err := internal.CalcAdcPlusConfig(options...)
	if err != nil {
		return nil, err
	}

	var copts []option.ClientOption
	if len(config.CredentialsJSON) > 0 {
		copts = []option.ClientOption{option.WithCredentialsJSON(config.CredentialsJSON)}
	}

	if config.TargetPrincipal != "" {
		idCfg := impersonate.IDTokenConfig{
			Audience:        audience,
			TargetPrincipal: config.TargetPrincipal,
			Delegates:       config.Delegates,
			// Cloud IAP requires email claim.
			IncludeEmail: true,
		}
		return impersonate.IDTokenSource(ctx, idCfg, copts...)
	}
	return idtoken.NewTokenSource(ctx, audience, copts...)
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
		var copts []option.ClientOption
		if len(config.CredentialsJSON) > 0 {
			copts = []option.ClientOption{option.WithCredentialsJSON(config.CredentialsJSON)}
		}
		return impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: config.TargetPrincipal,
			Delegates:       config.Delegates,
			Scopes:          config.Scopes,
		}, copts...)
	}

	if len(config.CredentialsJSON) > 0 {
		cred, err := google.CredentialsFromJSON(ctx, config.CredentialsJSON, config.Scopes...)
		if err != nil {
			return nil, err
		}
		return cred.TokenSource, nil
	}
	return google.DefaultTokenSource(ctx, config.Scopes...)
}
