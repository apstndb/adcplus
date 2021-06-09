package signer

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/iam/credentials/apiv1"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
)

type iamCredentialsSigner struct {
	target    string
	delegates []string
	ts        oauth2.TokenSource
}

func (s *iamCredentialsSigner) ServiceAccount(context.Context) string {
	return s.target
}

func (s *iamCredentialsSigner) Signer(ctx context.Context) func([]byte) ([]byte, error) {
	return func(b []byte) ([]byte, error) {
		// Actually, WithTokenSource(nil) will be ignored so this condition doesn't make any changes.
		var opts []option.ClientOption
		if s.ts != nil {
			opts = []option.ClientOption{option.WithTokenSource(s.ts)}
		}

		client, err := credentials.NewIamCredentialsClient(ctx, opts...)
		if err != nil {
			return nil, err
		}
		defer client.Close()

		resp, err := client.SignBlob(ctx, &credentialspb.SignBlobRequest{
			Name:      s.target,
			Delegates: s.delegates,
			Payload:   b,
		})
		if err != nil {
			return nil, fmt.Errorf("iamCredentialsSigner can't call SignBlob as %s: %w", s.target, err)
		}
		return resp.GetSignedBlob(), nil
	}
}

// IamCredentialsSigner makes new Signer.
// targetPrincipal and delegates is passed to iamcredentials.SignBlob.
// if ts is nil, ADC will be used.
func IamCredentialsSigner(targetPrincipal string, delegates []string, ts oauth2.TokenSource) (Signer, error) {
	if targetPrincipal == "" {
		return nil, errors.New("signer.IamCredentialsSigner requires non-empty targetPrincipal")
	}
	return &iamCredentialsSigner{
		target:    targetPrincipal,
		delegates: delegates,
		ts:        ts,
	}, nil
}
