package adcplus

import (
	"github.com/apstndb/adcplus/internal/config"
)

type Option func(*config.AdcPlusConfig) error

func WithExperimentalAppEngineSigner(enable bool) Option {
	return func(config *config.AdcPlusConfig) error {
		config.EnableAppengineSigner = enable
		return nil
	}
}

func WithTargetPrincipal(targetPrincipal string) Option {
	return func(config *config.AdcPlusConfig) error {
		config.TargetPrincipal = targetPrincipal
		return nil
	}
}

func WithDelegates(delegates ...string) Option {
	return func(config *config.AdcPlusConfig) error {
		config.Delegates = delegates
		return nil
	}
}

func WithScopes(scopes ...string) Option {
	return func(config *config.AdcPlusConfig) error {
		config.Scopes = scopes
		return nil
	}
}
func WithCredentialsFile(filename string) Option {
	return func(config *config.AdcPlusConfig) error {
		config.CredentialsFile = filename
		return nil
	}
}
func WithCredentialsJSON(j []byte) Option {
	return func(config *config.AdcPlusConfig) error {
		config.CredentialsJSON = j
		return nil
	}
}
