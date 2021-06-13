package signer

import (
	"context"
	"errors"

	"cloud.google.com/go/iam/credentials/apiv1"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
	"google.golang.org/api/option"
	credentials2 "google.golang.org/genproto/googleapis/iam/credentials/v1"
)

type iamCredentialsSigner struct {
	target    string
	delegates []string
	ts        oauth2.TokenSource
}

func (s *iamCredentialsSigner) SignJwt(ctx context.Context, c string) (string, error) {
	// Actually, WithTokenSource(nil) will be ignored so this condition doesn't make any changes.
	var opts []option.ClientOption
	if s.ts != nil {
		opts = []option.ClientOption{option.WithTokenSource(s.ts)}
	}

	client, err := credentials.NewIamCredentialsClient(ctx, opts...)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.SignJwt(ctx, &credentials2.SignJwtRequest{
		Name:      s.target,
		Delegates: s.delegates,
		Payload:   c,
	})
	if err != nil {
		return "", xerrors.Errorf("iamCredentialsSigner can't call SignBlob as %s: %w", s.target, err)
	}
	return resp.GetSignedJwt(), nil
}

func (s *iamCredentialsSigner) ServiceAccount(context.Context) string {
	return s.target
}

func (s *iamCredentialsSigner) SignBlob(ctx context.Context, b []byte) (string, []byte, error) {
	// Actually, WithTokenSource(nil) will be ignored so this condition doesn't make any changes.
	var opts []option.ClientOption
	if s.ts != nil {
		opts = []option.ClientOption{option.WithTokenSource(s.ts)}
	}

	client, err := credentials.NewIamCredentialsClient(ctx, opts...)
	if err != nil {
		return "", nil, err
	}
	defer client.Close()

	resp, err := client.SignBlob(ctx, &credentials2.SignBlobRequest{
		Name:      s.target,
		Delegates: s.delegates,
		Payload:   b,
	})
	if err != nil {
		return "", nil, xerrors.Errorf("iamCredentialsSigner can't call SignBlob as %s: %w", s.target, err)
	}
	return resp.GetKeyId(), resp.GetSignedBlob(), nil
}

// newIamCredentialsSigner makes new Signer.
// targetPrincipal and delegates is passed to iamcredentials.SignBlob.
// if ts is nil, ADC will be used.
func newIamCredentialsSigner(targetPrincipal string, delegates []string, ts oauth2.TokenSource) (Signer, error) {
	if targetPrincipal == "" {
		return nil, errors.New("signer.newIamCredentialsSigner requires non-empty targetPrincipal")
	}
	return &iamCredentialsSigner{
		target:    targetPrincipal,
		delegates: delegates,
		ts:        ts,
	}, nil
}
