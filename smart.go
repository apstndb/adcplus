package signer

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	goauth2 "google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
)

const impSaEnvName = "CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT"

// All patterns are defined in
// https://github.com/golang/oauth2/blob/f6687ab2804cbebdfdeef385bee94918b1ce83de/google/google.go#L93-L98
const (
	serviceAccountKey  = "service_account"
	userCredentialsKey = "authorized_user"
	externalAccountKey = "external_account"
)

type smartSignerConfig struct {
	targetPrincipal string
	delegates       []string
	enableAppengineSigner bool
}

type Option func(*smartSignerConfig) error

func WithExperimentalAppEngineSigner(enable bool) Option {
	return func(config *smartSignerConfig) error {
		config.enableAppengineSigner = enable
		return nil
	}
}

func WithTargetPrincipal(targetPrincipal string) Option {
	return func(config *smartSignerConfig) error {
		config.targetPrincipal = targetPrincipal
		return nil
	}
}

func WithDelegates(delegates ...string) Option {
	return func(config *smartSignerConfig) error {
		config.delegates = delegates
		return nil
	}
}

func calcSmartSignerConfig(opts ...Option) (*smartSignerConfig, error) {
	var config smartSignerConfig
	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, err
		}
	}

	if config.targetPrincipal == "" && len(config.delegates) > 0 {
		return nil, fmt.Errorf("targetPrincipal is set but delegates is not set: %s", config.delegates)
	}
	if impSaVal := os.Getenv(impSaEnvName); config.targetPrincipal == "" && impSaVal != "" {
		config.targetPrincipal, config.delegates = parseDelegateChain(impSaVal)
	}
	return &config, nil
}

// SmartSigner create signer for ADC with optional impersonation.
// Impersonation setting is supplied from below in descending order of priority.
// 	1. options e.g. signer.WithTargetPrincipal, signer.WithDelegates
// 	2. `CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT` environment variable
// If impersonation is not applied, all credentials except App Engine 1st gen(only Go 1.11) and Service Account Key need a Token Creator role to themselves.
// 	* https://cloud.google.com/iam/docs/creating-short-lived-service-account-credentials?hl=en
// 	* https://cloud.google.com/iam/docs/impersonating-service-accounts?hl=en
func SmartSigner(ctx context.Context, options ...Option) (Signer, error) {
	config, err := calcSmartSignerConfig(options...)
	if err != nil {
		return nil, err
	}

	// Find credentials in ADC manner.
	// See also https://google.aip.dev/auth/4110.
	cred, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}

	// If targetPrincipal is populated, use ADC with impersonation
	if config.targetPrincipal != "" {
		return IamCredentialsSigner(config.targetPrincipal, config.delegates, cred.TokenSource)
	}

	if len(cred.JSON) != 0 {
		t, err := credentialType(cred.JSON)
		if err != nil {
			return nil, err
		}

		switch t {
		case userCredentialsKey:
			return nil, fmt.Errorf("authorized_user is unsupported so set CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT or use other credentials")
		case serviceAccountKey:
			return ServiceAccountSigner(cred.JSON)
		case externalAccountKey:
			fallthrough
		default:
			// fallthrough to IAM Credentials
		}
	} else {
		// App Engine or metadata server credentials are possible in this branch.
		// appengine.SignBytes can sign blob without Token Creator roles in go111 runtime.
		// Ensure initialization doesn't need an appengine context.
		if config.enableAppengineSigner && isSupportedAppEngineRuntime() {
			return AppEngineSigner()
		}
		// fall through to IAM Credentials because metadata server doesn't have SignBlob
	}

	ts := cred.TokenSource

	// Get email from tokeninfo of ADC
	oauth2Svc, err := goauth2.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, err
	}
	resp, err := oauth2Svc.Tokeninfo().Do()
	if err != nil {
		return nil, err
	}

	if resp.Email == "" {
		return nil, errors.New("signer.SmartSigner can't infer email")
	}
	// Use itself as target
	return IamCredentialsSigner(resp.Email, nil, ts)
}
