package signer

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/apstndb/adcplus/internal"
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

// newServiceAccountSigner returns Signer which can sign without any network access.
func newServiceAccountSigner(jsonKey []byte) (Signer, error) {
	var key struct {
		Type         string `json:"type"`
		ClientEmail  string `json:"client_email"`
		PrivateKeyID string `json:"private_key_id"`
		PrivateKey   string `json:"private_key"`
	}
	if err := json.Unmarshal(jsonKey, &key); err != nil {
		return nil, err
	}
	switch key.Type {
	case internal.ServiceAccountKey, internal.GDCHServiceAccountKey:
	default:
		return nil, fmt.Errorf("unsupported local signing credential type %q", key.Type)
	}
	if key.ClientEmail == "" || key.PrivateKey == "" {
		return nil, errors.New("service account credentials require client_email and private_key")
	}

	block, _ := pem.Decode([]byte(key.PrivateKey))
	if block == nil {
		return nil, errors.New("failed to decode PEM private key")
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key failed rsa.PrivateKey type assertion")
	}
	return &serviceAccountSigner{clientEmail: key.ClientEmail, rsaKey: rsaKey, keyId: key.PrivateKeyID}, nil
}
