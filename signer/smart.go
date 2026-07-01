package signer

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	goauth2 "google.golang.org/api/oauth2/v1"
	gapioption "google.golang.org/api/option"

	"github.com/apstndb/adcplus"
	"github.com/apstndb/adcplus/internal"
)

const iamScope = "https://www.googleapis.com/auth/iam"
const userinfoEmailScope = "https://www.googleapis.com/auth/userinfo.email"

var lookupTokeninfoEmail = tokeninfoEmail

// SmartSigner create signer for ADC with optional impersonation.
// Impersonation setting is supplied from below in descending order of priority.
//  1. options e.g. signer.WithTargetPrincipal, signer.WithDelegates
//  2. `CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT` environment variable
//
// If impersonation is not applied, all credentials except App Engine 1st gen(only Go 1.11) and Service Account Key need a Token Creator role to themselves.
//   - https://cloud.google.com/iam/docs/creating-short-lived-service-account-credentials?hl=en
//   - https://cloud.google.com/iam/docs/impersonating-service-accounts?hl=en
func SmartSigner(ctx context.Context, options ...adcplus.Option) (Signer, error) {
	config, err := internal.CalcAdcPlusConfig(options...)
	if err != nil {
		return nil, err
	}

	if config.TokenSource != nil {
		if config.TargetPrincipal != "" {
			return newIamCredentialsSigner(config.TargetPrincipal, config.Delegates, config.TokenSource)
		}
		email, err := inferServiceAccountEmailFromTokenSource(ctx, config.TokenSource, "the provided TokenSource")
		if err != nil {
			return nil, err
		}
		return newIamCredentialsSigner(email, nil, config.TokenSource)
	}

	var cred *google.Credentials
	if len(config.CredentialsJSON) > 0 {
		cred, err = internal.GoogleCredentialsFromJSON(ctx, config.CredentialsJSON, iamScope, userinfoEmailScope)
	} else {
		// Find credentials in ADC manner.
		// See also https://google.aip.dev/auth/4110.
		cred, err = google.FindDefaultCredentials(ctx, iamScope, userinfoEmailScope)
	}
	if err != nil {
		return nil, err
	}

	ts := cred.TokenSource

	// If targetPrincipal is populated, use ADC with impersonation
	if config.TargetPrincipal != "" {
		return newIamCredentialsSigner(config.TargetPrincipal, config.Delegates, ts)
	}

	credType, err := internal.InferADCCredentialType(cred)
	if err != nil {
		return nil, err
	}

	var email string
	switch credType {
	case internal.UserCredentialsKey:
		return nil, fmt.Errorf("authorized_user is unsupported so set CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT or use other credentials")
	case internal.ExternalAccountAuthorizedUserKey:
		return nil, fmt.Errorf("external_account_authorized_user is unsupported so set CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT or use other credentials")
	case internal.ServiceAccountKey, internal.GDCHServiceAccountKey:
		return newServiceAccountSigner(cred.JSON)
	case internal.ImpersonatedServiceAccountKey:
		return signerFromImpersonatedServiceAccountJSON(ctx, cred.JSON)
	case internal.ExternalAccountKey:
		// fallthrough to IAM Credentials
	case internal.ComputeMetadataCredential:
		// Ensure initialization doesn't need an appengine context.
		if config.EnableAppEngineSigner && isSupportedAppEngineRuntime() {
			return newAppEngineSigner()
		}
		email, err = metadata.Email("default")
		if err != nil {
			return nil, err
		}
		// fallthrough to IAM Credentials because metadata server doesn't have SignBlob
	default:
		return nil, fmt.Errorf("unsupported credential type %q for SmartSigner without impersonation", credType)
	}

	if email == "" {
		// Get email from tokeninfo of ADC.
		email, err = inferServiceAccountEmailFromTokenSource(ctx, ts, "ADC TokenSource")
		if err != nil {
			return nil, err
		}
	}
	// Use itself as target
	return newIamCredentialsSigner(email, nil, ts)
}

func inferServiceAccountEmailFromTokenSource(ctx context.Context, ts oauth2.TokenSource, sourceDescription string) (string, error) {
	email, err := lookupTokeninfoEmail(ctx, ts)
	if err != nil {
		return "", err
	}
	if email == "" {
		return "", fmt.Errorf("signer.SmartSigner can't infer email from %s", sourceDescription)
	}
	if !isServiceAccountEmail(email) {
		return "", fmt.Errorf(
			"signer.SmartSigner inferred non-service-account email %q from %s; set adcplus.WithTargetPrincipal to the service account to sign as",
			email,
			sourceDescription,
		)
	}
	return email, nil
}

func tokeninfoEmail(ctx context.Context, ts oauth2.TokenSource) (string, error) {
	oauth2Svc, err := goauth2.NewService(ctx, gapioption.WithTokenSource(ts))
	if err != nil {
		return "", err
	}
	resp, err := oauth2Svc.Tokeninfo().Do()
	if err != nil {
		return "", err
	}
	return resp.Email, nil
}

func isServiceAccountEmail(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	return strings.Contains(email, "@") && strings.HasSuffix(email, ".gserviceaccount.com")
}
