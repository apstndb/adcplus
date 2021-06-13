package option

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
