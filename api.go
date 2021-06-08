package signer

import "context"

type Signer interface {
	ServiceAccount() string
	Signer(context.Context) func([]byte) ([]byte, error)
}
