package signer

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"golang.org/x/oauth2/google"
)

type serviceAccountSigner struct {
	clientEmail string
	rsaKey      *rsa.PrivateKey
	keyId       string
}

func (s *serviceAccountSigner) SignJwt(ctx context.Context, claims string) (string, error) {
	_, signed, err := signJwtHelper(ctx, claims, s.keyId, s)
	return signed, err
}

func (s *serviceAccountSigner) ServiceAccount(context.Context) string {
	return s.clientEmail
}
func (s serviceAccountSigner) SignBlob(_ context.Context, b []byte) (string, []byte, error) {
	sum := sha256.Sum256(b)
	sig, err := rsa.SignPKCS1v15(rand.Reader, s.rsaKey, crypto.SHA256, sum[:])
	return s.keyId, sig, err
}

// ServiceAccountSigner returns Signer which can sign without any network access.
func ServiceAccountSigner(jsonKey []byte) (Signer, error) {
	// google.JWTConfigFromJSON is used to extract client_email, private_key, private_key_id
	// because credentialFile struct is not exported.
	cfg, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(cfg.PrivateKey)
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key failed rsa.PrivateKey type assertion")
	}
	return &serviceAccountSigner{clientEmail: cfg.Email, rsaKey: rsaKey, keyId: cfg.PrivateKeyID}, nil
}
