package signer

import "context"

type Signer interface {
	ServiceAccount(context.Context) string
	SignBlob(context.Context, []byte) (string, []byte, error)
}

func SignWithoutKeyAdaptor(ctx context.Context, signer Signer) func([]byte) ([]byte, error){
	return func(b []byte) ([]byte, error) {
		_, sig, err := signer.SignBlob(ctx, b)
		return sig, err
	}
}
