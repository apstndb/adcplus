package signer

import (
	"context"

	"google.golang.org/appengine"
)

type appengineTokenSigner struct{}

func (a appengineTokenSigner) ServiceAccount(ctx context.Context) string {
	email, err := appengine.ServiceAccount(ctx)
	if err != nil{
		return ""
	}
	return email
}

func (a appengineTokenSigner) Signer(ctx context.Context) func([]byte) ([]byte, error) {
	return func(b []byte) ([]byte, error) {
		_, signed, err := appengine.SignBytes(ctx, b)
		return signed, err
	}
}

func AppEngineSigner() (Signer, error){
	return &appengineTokenSigner{}, nil
}

