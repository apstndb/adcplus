package tokensource

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"

	"github.com/apstndb/adcplus"
	"github.com/apstndb/adcplus/internal"
)

const cloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"

const smartIDTokenSourceImpersonationHint = "set CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT, adcplus.WithTargetPrincipal (optionally with adcplus.WithDelegates), or use other credentials"

// SmartIDTokenSource generate oauth2.TokenSource which generates ID token and supports CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT environment variable.
func SmartIDTokenSource(ctx context.Context, audience string, options ...adcplus.Option) (oauth2.TokenSource, error) {
	config, err := internal.CalcAdcPlusConfig(options...)
	if err != nil {
		return nil, err
	}

	var copts []option.ClientOption
	if config.TokenSource != nil {
		copts = []option.ClientOption{option.WithTokenSource(config.TokenSource)}
	} else if len(config.CredentialsJSON) > 0 {
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

	if config.TokenSource != nil {
		return config.TokenSource, nil
	}

	if len(config.CredentialsJSON) > 0 {
		credType, err := internal.CredentialTypeFromJSON(config.CredentialsJSON)
		if err != nil {
			return nil, err
		}
		if err := validateSmartIDTokenSourceCredentialType(credType); err != nil {
			return nil, err
		}
	} else {
		cred, err := google.FindDefaultCredentials(ctx)
		if err != nil {
			return nil, err
		}
		credType, err := internal.InferADCCredentialType(cred)
		if err != nil {
			return nil, err
		}
		if err := validateSmartIDTokenSourceCredentialType(credType); err != nil {
			return nil, err
		}
		// Reuse the ADC lookup for idtoken.NewTokenSource. When cred.JSON is empty
		// (compute_metadata / App Engine ADC), idtoken still uses the metadata
		// identity endpoint rather than the access-token TokenSource.
		copts = append(copts, option.WithCredentials(cred))
	}

	return idtoken.NewTokenSource(ctx, audience, copts...)
}

func validateSmartIDTokenSourceCredentialType(credType string) error {
	switch credType {
	case internal.UserCredentialsKey:
		return fmt.Errorf("authorized_user is unsupported for SmartIDTokenSource without impersonation; %s", smartIDTokenSourceImpersonationHint)
	case internal.ExternalAccountKey:
		return fmt.Errorf("external_account is unsupported for SmartIDTokenSource without impersonation; %s (STS support is tracked in https://github.com/apstndb/adcplus/issues/3)", smartIDTokenSourceImpersonationHint)
	case internal.ExternalAccountAuthorizedUserKey:
		return fmt.Errorf("external_account_authorized_user is unsupported for SmartIDTokenSource without impersonation; %s", smartIDTokenSourceImpersonationHint)
	default:
		return nil
	}
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
		if config.TokenSource != nil {
			copts = []option.ClientOption{option.WithTokenSource(config.TokenSource)}
		} else if len(config.CredentialsJSON) > 0 {
			copts = []option.ClientOption{option.WithCredentialsJSON(config.CredentialsJSON)}
		}
		return impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: config.TargetPrincipal,
			Delegates:       config.Delegates,
			Scopes:          config.Scopes,
		}, copts...)
	}

	if config.TokenSource != nil {
		return config.TokenSource, nil
	}

	if len(config.CredentialsJSON) > 0 {
		cred, err := internal.GoogleCredentialsFromJSON(ctx, config.CredentialsJSON, config.Scopes...)
		if err != nil {
			return nil, err
		}
		return cred.TokenSource, nil
	}
	return google.DefaultTokenSource(ctx, config.Scopes...)
}
