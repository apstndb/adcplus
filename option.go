package adcplus

import (
	"golang.org/x/oauth2"

	"github.com/apstndb/adcplus/internal/config"
)

type Option func(*config.AdcPlusConfig) error

// WithExperimentalAppEngineSigner returns Option which specifies to use appengine.SignBytes by Signer.
// Caution: It makes the signer to require an App Engine context. If you don't know this meaning, don't set this option.
func WithExperimentalAppEngineSigner(enable bool) Option {
	return func(config *config.AdcPlusConfig) error {
		config.EnableAppEngineSigner = enable
		return nil
	}
}

// WithJWTAccessWithScope returns Option which specifies to use scope-based
// self-signed JWT access tokens when SmartAccessTokenSource is clearly using
// service account JSON credentials.
//
// This is an AIP-4111 opt-in optimization for services known to support
// self-signed JWT access tokens with scopes. It only applies to service account
// JSON credentials supplied by WithCredentialsJSON or WithCredentialsFile, and
// does not apply to impersonation, WithTokenSource, ADC, ID tokens, or
// non-service-account credential types.
func WithJWTAccessWithScope(enable bool) Option {
	return func(config *config.AdcPlusConfig) error {
		config.UseJWTAccessWithScope = enable
		return nil
	}
}

// WithTargetPrincipal returns Option which specifies the target principal for impersonation.
func WithTargetPrincipal(targetPrincipal string) Option {
	return func(config *config.AdcPlusConfig) error {
		config.TargetPrincipal = targetPrincipal
		return nil
	}
}

// WithDelegates returns Option which specifies the delegate chain for impersonation.
func WithDelegates(delegates ...string) Option {
	return func(config *config.AdcPlusConfig) error {
		config.Delegates = delegates
		return nil
	}
}

// WithScopes returns Option which specifies the scopes of the access token.
func WithScopes(scopes ...string) Option {
	return func(config *config.AdcPlusConfig) error {
		config.Scopes = scopes
		return nil
	}
}

// WithCredentialsFile returns Option which specifies the path of credentials.
// If filename is empty string, it will be ignored.
func WithCredentialsFile(filename string) Option {
	return func(config *config.AdcPlusConfig) error {
		config.CredentialsFile = filename
		return nil
	}
}

// WithCredentialsJSON returns Option which specifies the content of credentials.
// If j is nil or empty slice, it will be ignored.
func WithCredentialsJSON(j []byte) Option {
	return func(config *config.AdcPlusConfig) error {
		config.CredentialsJSON = j
		return nil
	}
}

// WithTokenSource returns Option which overrides the underlying oauth2.TokenSource.
// When set, SmartAccessTokenSource and SmartIDTokenSource use this source instead of
// deriving one from ADC credentials. SmartSigner uses it for IAM Credentials API auth.
func WithTokenSource(ts oauth2.TokenSource) Option {
	return func(config *config.AdcPlusConfig) error {
		config.TokenSource = ts
		return nil
	}
}
