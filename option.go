package adcplus

import (
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
