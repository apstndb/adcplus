package signer

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestSignWithoutKeyAdaptor(t *testing.T) {
	t.Run("returns signature discarding key", func(t *testing.T) {
		s := &mockSigner{
			keyID:     "key-123",
			signature: []byte("the-signature"),
		}

		signFn := SignWithoutKeyAdaptor(context.Background(), s)
		got, err := signFn([]byte("data-to-sign"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, []byte("the-signature")) {
			t.Errorf("signature = %q, want %q", got, "the-signature")
		}
	})

	t.Run("propagates SignBlob error", func(t *testing.T) {
		s := &mockSigner{
			signBlobErr: errors.New("sign failed"),
		}

		signFn := SignWithoutKeyAdaptor(context.Background(), s)
		_, err := signFn([]byte("data"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "sign failed" {
			t.Errorf("error = %q, want %q", err.Error(), "sign failed")
		}
	})
}
