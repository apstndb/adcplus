package signer

import (
	"context"
	"errors"
	"fmt"

	"github.com/apstndb/adcplus"
	"github.com/apstndb/adcplus/internal"
	"golang.org/x/oauth2/google"
	goauth2 "google.golang.org/api/oauth2/v1"
	gapioption "google.golang.org/api/option"
)

// SmartSigner create signer for ADC with optional impersonation.
// Impersonation setting is supplied from below in descending order of priority.
// 	1. options e.g. signer.WithTargetPrincipal, signer.WithDelegates
// 	2. `CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT` environment variable
// If impersonation is not applied, all credentials except App Engine 1st gen(only Go 1.11) and Service Account Key need a Token Creator role to themselves.
// 	* https://cloud.google.com/iam/docs/creating-short-lived-service-account-credentials?hl=en
// 	* https://cloud.google.com/iam/docs/impersonating-service-accounts?hl=en
func SmartSigner(ctx context.Context, options ...adcplus.Option) (Signer, error) {
	config, err := internal.CalcAdcPlusConfig(options...)
	if err != nil {
		return nil, err
	}

	var cred *google.Credentials
	if len(config.CredentialsJSON) > 0 {
		cred, err = google.CredentialsFromJSON(ctx, config.CredentialsJSON)
	} else {
		// Find credentials in ADC manner.
		// See also https://google.aip.dev/auth/4110.
		cred, err = google.FindDefaultCredentials(ctx)
	}
	if err != nil {
		return nil, err
	}

	// If targetPrincipal is populated, use ADC with impersonation
	if config.TargetPrincipal != "" {
		return newIamCredentialsSigner(config.TargetPrincipal, config.Delegates, cred.TokenSource)
	}

	credType, err := internal.InferADCCredentialType(cred)
	if err != nil {
		return nil, err
	}

	switch credType {
	case internal.UserCredentialsKey:
		return nil, fmt.Errorf("authorized_user is unsupported so set CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT or use other credentials")
	case internal.ServiceAccountKey:
		return newServiceAccountSigner(cred.JSON)
	case internal.ExternalAccountKey:
		// fallthrough to IAM Credentials
	case internal.ComputeMetadataCredential:
		// Ensure initialization doesn't need an appengine context.
		if config.EnableAppEngineSigner && isSupportedAppEngineRuntime() {
			return newAppEngineSigner()
		}
		// fallthrough to IAM Credentials because metadata server doesn't have SignBlob
	default:
		// fallthrough to IAM Credentials
	}

	ts := cred.TokenSource

	// Get email from tokeninfo of ADC
	oauth2Svc, err := goauth2.NewService(ctx, gapioption.WithTokenSource(ts))
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
	return newIamCredentialsSigner(resp.Email, nil, ts)
}
