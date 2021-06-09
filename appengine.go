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

func (a appengineTokenSigner) SignBlob(ctx context.Context, b[]byte) (string, []byte, error) {
		return appengine.SignBytes(ctx, b)
}

func AppEngineSigner() (Signer, error){
	return &appengineTokenSigner{}, nil
}

