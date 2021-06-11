package signer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
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
		key, signed, err := signJwtHelper(ctx, c, kid, s)
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

func signJwtHelper(ctx context.Context, claimsJson string, kid string, s Signer) (key string, signed string, err error) {
	var buf bytes.Buffer
	type jwtHeader struct {
		Alg string `json:"alg"`
		Kid string `json:"kig"`
		Typ string `json:"typ"`
	}
	header, err := json.Marshal(jwtHeader{
		Alg: "RS256",
		Kid: kid,
		Typ: "JWT",
	})
	if err != nil {
		return "", "", err
	}

	buf.WriteString(base64.RawURLEncoding.EncodeToString(header))
	buf.WriteByte('.')
	buf.WriteString(base64.RawURLEncoding.EncodeToString([]byte(claimsJson)))
	key, sig, err := s.SignBlob(ctx, buf.Bytes())
	if err != nil {
		return "", "", err
	}
	buf.WriteByte('.')
	buf.WriteString(base64.RawURLEncoding.EncodeToString(sig))
	return key, buf.String(), nil
}

func (s *appengineSigner) ServiceAccount(ctx context.Context) string {
	email, err := appengine.ServiceAccount(ctx)
	if err != nil{
		return ""
	}
	return email
}

func (s *appengineSigner) SignBlob(ctx context.Context, b[]byte) (string, []byte, error) {
		return appengine.SignBytes(ctx, b)
}

func AppEngineSigner() (Signer, error){
	return &appengineSigner{}, nil
}

func isSupportedAppEngineRuntime() bool {
	return appengine.IsStandard() && os.Getenv("GAE_RUNTIME") == "go111"
}
