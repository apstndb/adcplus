package internal

import (
	"os"
	"reflect"
	"testing"

	"github.com/apstndb/adcplus"
)

func TestParseDelegateChain(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		wantTarget      string
		wantDelegates   []string
		wantErr         bool
	}{
		{
			name:          "single element (no delegates)",
			input:         "sa@project.iam.gserviceaccount.com",
			wantTarget:    "sa@project.iam.gserviceaccount.com",
			wantDelegates: []string{},
		},
		{
			name:          "multiple elements (with delegates)",
			input:         "d1@p.iam.gserviceaccount.com,d2@p.iam.gserviceaccount.com,target@p.iam.gserviceaccount.com",
			wantTarget:    "target@p.iam.gserviceaccount.com",
			wantDelegates: []string{"d1@p.iam.gserviceaccount.com", "d2@p.iam.gserviceaccount.com"},
		},
		{
			name:    "empty string returns error",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, delegates, err := ParseDelegateChain(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseDelegateChain(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if target != tt.wantTarget {
				t.Errorf("target = %q, want %q", target, tt.wantTarget)
			}
			if !reflect.DeepEqual(delegates, tt.wantDelegates) {
				t.Errorf("delegates = %v, want %v", delegates, tt.wantDelegates)
			}
		})
	}
}

func TestCredentialTypeFromJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr bool
	}{
		{
			name:  "service_account",
			input: []byte(`{"type":"service_account"}`),
			want:  "service_account",
		},
		{
			name:  "authorized_user",
			input: []byte(`{"type":"authorized_user"}`),
			want:  "authorized_user",
		},
		{
			name:  "external_account",
			input: []byte(`{"type":"external_account"}`),
			want:  "external_account",
		},
		{
			name:    "empty type",
			input:   []byte(`{"type":""}`),
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{invalid`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CredentialTypeFromJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CredentialTypeFromJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("CredentialTypeFromJSON() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCalcAdcPlusConfig(t *testing.T) {
	t.Run("applies options", func(t *testing.T) {
		cfg, err := CalcAdcPlusConfig(
			adcplus.WithTargetPrincipal("sa@p.iam.gserviceaccount.com"),
			adcplus.WithScopes("scope1", "scope2"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.TargetPrincipal != "sa@p.iam.gserviceaccount.com" {
			t.Errorf("TargetPrincipal = %q, want %q", cfg.TargetPrincipal, "sa@p.iam.gserviceaccount.com")
		}
		if !reflect.DeepEqual(cfg.Scopes, []string{"scope1", "scope2"}) {
			t.Errorf("Scopes = %v, want %v", cfg.Scopes, []string{"scope1", "scope2"})
		}
	})

	t.Run("delegates without targetPrincipal returns error", func(t *testing.T) {
		_, err := CalcAdcPlusConfig(
			adcplus.WithDelegates("d1@p.iam.gserviceaccount.com"),
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("CredentialsJSON and CredentialsFile are mutually exclusive", func(t *testing.T) {
		_, err := CalcAdcPlusConfig(
			adcplus.WithCredentialsJSON([]byte(`{"type":"service_account"}`)),
			adcplus.WithCredentialsFile("/some/file.json"),
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("env CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT parsed", func(t *testing.T) {
		os.Setenv("CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT", "d1@p.iam.gserviceaccount.com,target@p.iam.gserviceaccount.com")
		defer os.Unsetenv("CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT")

		cfg, err := CalcAdcPlusConfig()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.TargetPrincipal != "target@p.iam.gserviceaccount.com" {
			t.Errorf("TargetPrincipal = %q, want %q", cfg.TargetPrincipal, "target@p.iam.gserviceaccount.com")
		}
		if !reflect.DeepEqual(cfg.Delegates, []string{"d1@p.iam.gserviceaccount.com"}) {
			t.Errorf("Delegates = %v, want %v", cfg.Delegates, []string{"d1@p.iam.gserviceaccount.com"})
		}
	})

	t.Run("env ignored when targetPrincipal already set", func(t *testing.T) {
		os.Setenv("CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT", "env-target@p.iam.gserviceaccount.com")
		defer os.Unsetenv("CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT")

		cfg, err := CalcAdcPlusConfig(
			adcplus.WithTargetPrincipal("opt-target@p.iam.gserviceaccount.com"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.TargetPrincipal != "opt-target@p.iam.gserviceaccount.com" {
			t.Errorf("TargetPrincipal = %q, want %q", cfg.TargetPrincipal, "opt-target@p.iam.gserviceaccount.com")
		}
	})

	t.Run("no options returns empty config", func(t *testing.T) {
		// Clear env to avoid interference
		os.Unsetenv("CLOUDSDK_AUTH_IMPERSONATE_SERVICE_ACCOUNT")

		cfg, err := CalcAdcPlusConfig()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.TargetPrincipal != "" {
			t.Errorf("TargetPrincipal = %q, want empty", cfg.TargetPrincipal)
		}
	})
}
