package signer

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/appengine"
)

type appengineSigner struct{}

func (s *appengineSigner) SignJwt(ctx context.Context, c string) (string, error) {
	certificates, err := appengine.PublicCertificates(ctx)
	if err != nil {
		return "", err
	}
	for _, certificate := range certificates {
		kid := certificate.KeyName
		key, signed, err := SignJwtHelper(ctx, c, kid, s)
		if err != nil {
			return "", err
		}
		if key != certificate.KeyName {
			continue
		}
		return signed, nil
	}
	return "", fmt.Errorf("key not matched")
}

func (s *appengineSigner) ServiceAccount(ctx context.Context) string {
	email, err := appengine.ServiceAccount(ctx)
	if err != nil {
		return ""
	}
	return email
}

func (s *appengineSigner) SignBlob(ctx context.Context, b []byte) (string, []byte, error) {
	return appengine.SignBytes(ctx, b)
}

func newAppEngineSigner() (Signer, error) {
	return &appengineSigner{}, nil
}

func isSupportedAppEngineRuntime() bool {
	return appengine.IsStandard() && os.Getenv("GAE_RUNTIME") == "go111"
}
