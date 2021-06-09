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
}

func (s *serviceAccountSigner) ServiceAccount(context.Context) string {
	return s.clientEmail
}

func (s serviceAccountSigner) Signer(context.Context) func([]byte) ([]byte, error) {
	return func(b []byte) ([]byte, error) {
		sum := sha256.Sum256(b)
		return rsa.SignPKCS1v15(rand.Reader, s.rsaKey, crypto.SHA256, sum[:])
	}
}

// ServiceAccountSigner returns Signer which can sign without any network access.
func ServiceAccountSigner(jsonKey []byte) (Signer, error) {
	// google.JWTConfigFromJSON is actually alternative json.Marshal because credentialFile is not exported
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
	return &serviceAccountSigner{clientEmail: cfg.Email, rsaKey: rsaKey}, nil
}
