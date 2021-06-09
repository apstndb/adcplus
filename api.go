package signer

import "context"

type Signer interface {
	ServiceAccount(context.Context) string
	Signer(context.Context) func([]byte) ([]byte, error)
}
