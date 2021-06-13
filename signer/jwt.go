package signer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
)

func SignJwtHelper(ctx context.Context, claimsJson string, kid string, s Signer) (key string, signed string, err error) {
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
